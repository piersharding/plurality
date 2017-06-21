package command

import (
    "os"
    "io"
    "path/filepath"
    "fmt"
    "github.com/op/go-logging"
    "os/exec"
    "github.com/codegangsta/cli"
    "context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/filters"
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
    "archive/tar"
)

var root = os.Getenv("HOME") + "/.runc"
var container_run = root + "/run"
var container_root = root + "/container"
var prefix = "tmp-"
var log = logging.MustGetLogger("plurality")

func CmdCreate(c *cli.Context) error {
    if c.GlobalBool("debug") {
        logging.SetLevel(logging.DEBUG, "")
    }

    // make sure all directories exist
    _ = os.MkdirAll(container_root, os.ModePerm)

    if c.NArg() != 2 {
        return cli.NewExitError("Must supply <image:tag> and <container name>", -1)
    }
    // grab first arg - first is Docker image
    image := c.Args().Get(0)
    cname := c.Args().Get(1)
    target := container_root + "/" + cname

    log.Debug("Creating container ", cname, " from ", image)

    exists, err := fileExits(target)
    if exists || err != nil {
        return cli.NewExitError(fmt.Sprintf("Container directory %s already exists at: %s", cname, target), -1)
    }

    // find runC
    runc, err := exec.LookPath("runc")
    if err != nil {
        log.Fatal("We cannot find runc")
        return cli.NewExitError("Could not find runC", -1)
    }
    log.Debug("runC is available at ", runc)

    // find sudo
    sudo, err := exec.LookPath("sudo")
    if err != nil {
        log.Fatal("We cannot find sudo")
        return cli.NewExitError("Could not find sudo", -1)
    }
    log.Debug("sudo is available at ", sudo)

    // find tar
    tarpath, err := exec.LookPath("tar")
    if err != nil {
        log.Fatal("We cannot find tar")
        return cli.NewExitError("Could not find tar", -1)
    }
    log.Debug("tar is available at ", tarpath)

    dcli, err := client.NewEnvClient()
    if err != nil {
        log.Error("Could not create the docker client")
        panic(err)
    }

    // check for container
    filter := filters.NewArgs()
    filter.Add("name", prefix+cname)
    containers, err := dcli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})
    if err != nil {
        panic(err)
    }
    if len(containers) != 0 {
        log.Error("No. of containers found: ", len(containers))
        for _, container := range containers {
            log.Error("container: ", container.ID[:10], " ", container.Image)
            err = dcli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
            if err != nil {
                panic(err)
            }
        }
        return cli.NewExitError("Temporary container already exists: "+prefix+cname, -1)
    }

    // pull the image
    if ! c.Bool("nopull") {
        _, err = dcli.ImagePull(context.Background(), image, types.ImagePullOptions{All: true})
        if err != nil {
            panic(err)

        }
    }

    // now check that the image exists and is one
    filter = filters.NewArgs()
    filter.Add("reference", image)

    images, err := dcli.ImageList(context.Background(), types.ImageListOptions{Filters: filter})
    for _, image := range images {
        // img = image.ID
        log.Info("Found image: ", image.ID[:10], " ", image.RepoTags)
    }
    if len(images) != 1 {
        log.Error("No. of images found: ", len(images))
        return cli.NewExitError("Ambiguous reference to image: "+image[:10], -1)
    }

    // launch the image into an empty container ready for export
    containerConfig := &container.Config{
        Image: image,
        Hostname: prefix+cname,
        // Cmd:   []string{"true"}, // hopefully not required!
    }
    log.Info("Creating container: ", prefix+cname)
    resp, err := dcli.ContainerCreate(context.Background(), containerConfig, nil, nil, prefix+cname)
    if err != nil {
        log.Error("container: ", prefix+cname, " failed")
        panic(err)
    }

    log.Info("Running container: ", prefix+cname)
    if err := dcli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
        log.Error("container: ", resp.ID[:10], " failed")
        panic(err)
    }

    log.Info("Logs for container: ", prefix+cname)
    out, err := dcli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true})
    if err != nil {
        panic(err)
    }
    io.Copy(os.Stdout, out)

    // export the container and remove
    log.Info("Exporting container: ", prefix+cname)
    cout, err := dcli.ContainerExport(context.Background(), resp.ID)
    if err != nil {
        log.Error("Exporting container: ", prefix+cname, " failed")
        panic(err)
    }
    archivename := "/tmp/"+prefix+cname+".tar"
    log.Debug("Writing container: ", prefix+cname, " to: ", archivename)
    writeContainerTar(archivename, cout)

    err = dcli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
    if err != nil {
        panic(err)
    }

    // make sure all directories exist
    _ = os.MkdirAll(target+"/rootfs", os.ModePerm)

    // untar the file
    if c.GlobalBool("nosudo") {
        log.Debug("We are NOT using sudo to untar container: ", archivename)
        untar(archivename, target+"/rootfs")
    } else {
        extractTar(archivename, target+"/rootfs")
    }


    // generate spec file
    rout, err := exec.Command(runc, "--root", container_run, "spec", "--rootless", "--bundle", target).Output()
    if err != nil {
        log.Fatal(err)
        return cli.NewExitError("Could not generate spec file", -1)
    }
    log.Debug("runC spec: ", rout)

    return nil
}

func fileExits(path string) (bool, error) {
    log.Debug("checking path ", path)
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return true, err
}

func writeContainerTar(f string, out io.ReadCloser) {
    // err := os.Remove(f)
    outFile, err := os.Create(f)
    // handle err
    defer outFile.Close()
    _, err = io.Copy(outFile, out)
    if err != nil {
        panic(err)
    }
}

func extractTar(tarball, target string) error {
    out, err := exec.Command("sudo", "tar", "-xf", tarball, "-C", target).Output()
    log.Debug("tar ouput: ", out)
    if err != nil {
        log.Fatal(err)
        return cli.NewExitError("Could not unpack container image", -1)
    }
    return nil
}

func untar(tarball, target string) error {
    reader, err := os.Open(tarball)
    if err != nil {
        return err
    }
    defer reader.Close()
    tarReader := tar.NewReader(reader)

    for {
        header, err := tarReader.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        path := filepath.Join(target, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err = os.MkdirAll(path, info.Mode()); err != nil {
                return err
            }
            continue
        }

        file, err:= os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
        if err != nil {
            return err
        }
        defer file.Close()
        _, err =io.Copy(file, tarReader)
        if err != nil {
            return err
        }
    }
    return nil
}
