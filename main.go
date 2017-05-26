package main

import (
	"os"
	"github.com/codegangsta/cli"
    "github.com/op/go-logging"
)

var log = logging.MustGetLogger("plurality")
var format = logging.MustStringFormatter(
    "%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func main() {

    logging.SetFormatter(format)
    logging.SetLevel(logging.INFO, "")
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "piersharding"
	app.Email = ""
	app.Usage = "demonstration shell around runC - launch a container super locked down!"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

    // if cli.Context.String("errlevel") {
    //     fmt.Printf("errlevel: %s\n", app.Context.String("errlevel"))
    // }

	app.Run(os.Args)
}
