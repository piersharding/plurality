package command

import (
    "github.com/codegangsta/cli"
    "os"
    "fmt"
    "github.com/op/go-logging"
    "os/exec"
    "encoding/json"
    "io/ioutil"
)

func CmdRun(c *cli.Context) error {
    var log = logging.MustGetLogger("plurality")
    if c.GlobalBool("debug") {
        logging.SetLevel(logging.DEBUG, "")
    }

    // make sure all directories exist
    _ = os.MkdirAll(container_root, os.ModePerm)

    // fmt.Printf("No. Args: %d\n", c.NArg())
    // if c.NArg() > 0 {
    //     for i := 0; i < c.NArg(); i++ {
    //         fmt.Printf("Arg: %d - %q\n", i, c.Args().Get(i))
    //     }
    // }
    // fmt.Printf("Arg: %q\n", c.Args()[1:])

    if c.NArg() < 2 {
        return cli.NewExitError("Must supply <container name> <list of execution args>", -1)
    }
    // grab first arg - first is Docker image
    cname := c.Args().Get(0)
    target := container_root + "/" + cname

    log.Debug("Running container ", cname)

    exists, err := fileExits(target)
    if ! exists || err != nil {
        return cli.NewExitError(fmt.Sprintf("Container %s does not exist at: %s ", cname, target), -1)
    }

    // find runC
    runc, err := exec.LookPath("runc")
    if err != nil {
        log.Fatal("We cannot find runc")
        return cli.NewExitError("Could not find runC", -1)
    }
    log.Debug("runC is available at ", runc)

    // check that container exists
    var dat map[string]interface{}
    file, e := ioutil.ReadFile(target+"/config.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

    if err := json.Unmarshal(file, &dat); err != nil {
        panic(err)
    }
    bat := dat["process"].(map[string]interface{})
    bat["args"] = c.Args()[1:]
    file, err = json.Marshal(dat)
    err = ioutil.WriteFile(target+"/config.json", file, 0644)

    // try launching with command
    log.Debug("Going to run: ", runc, "--root", container_run, "run", "--bundle", target, cname)
    cmd := exec.Command(runc, "--root", container_run, "run", "--bundle", target, cname)
    cmd.Stdin = os.Stdin
    log.Debug("runC cmd: ", cmd)
    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        return cli.NewExitError(fmt.Sprintf("Could not execute: %s: %s [%s]", cname, stdoutStderr, err), -1)
    }
    fmt.Printf("%s", stdoutStderr)

    // do we background ?
    return nil
}
