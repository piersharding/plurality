package command

import (
    "github.com/codegangsta/cli"
    "github.com/op/go-logging"
)

func CmdStop(c *cli.Context) error {
   var log = logging.MustGetLogger("plurality")
   log.Info("TBD")

   return nil

}
