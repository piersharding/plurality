package command

import (
    "fmt"
    "github.com/codegangsta/cli"
    "github.com/op/go-logging"
    "os/exec"
)

func CmdDelete(c *cli.Context) error {
   var log = logging.MustGetLogger("plurality")
    if c.GlobalBool("debug") {
        logging.SetLevel(logging.DEBUG, "")
    }

    if c.NArg() != 1 {
        return cli.NewExitError("Must supply <container name>", -1)
    }
    cname := c.Args().Get(0)
    target := container_root + "/" + cname
    log.Debug("Deleting container ", cname)

    exists, err := fileExits(target)
    if ! exists || err != nil {
        return cli.NewExitError(fmt.Sprintf("Container %s does not exist at: %s ", cname, target), -1)
    }

    // find sudo
    sudo, err := exec.LookPath("sudo")
    if err != nil {
        log.Fatal("We cannot find sudo")
        return cli.NewExitError("Could not find sudo", -1)
    }
    log.Debug("sudo is available at ", sudo)

    out, err := exec.Command("sudo", "rm", "-rf", target).CombinedOutput()
    log.Debug("tar ouput: ", out)
    if err != nil {
        log.Fatal(err)
        return cli.NewExitError("Could not delete container image", -1)
    }

    return nil
}
