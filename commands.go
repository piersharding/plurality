package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/piersharding/plurality/command"
    "flag"
)

var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "ENV_DEBUG",
		Name:   "debug",

		Usage: "",
	},
    cli.BoolFlag{
        EnvVar: "ENV_NOSUDO",
        Name:   "nosudo",

        Usage: "",
    },
	cli.StringFlag{
		EnvVar: "ENV_ERRLEVEL",
		Name:   "errlevel",
        Value:  "ERROR",
		Usage: "",
	},
}

var Commands = []cli.Command{
    {
        Name:   "create",
        Aliases: []string{"c"},
        Usage:  "<image:tag> <container name>",
        Action: command.CmdCreate,
        Flags:  []cli.Flag{
                    cli.BoolFlag{
                        EnvVar: "ENV_NOPULL",
                        Name:   "nopull",
                        Usage: "do not pull image",
                    },
        },
    },
	{
		Name:   "run",
        Aliases: []string{"r"},
		Usage:  "[--nodaemon] <container name> [command line arguments]",
		Action: command.CmdRun,
		Flags:  []cli.Flag{
                    cli.BoolFlag{
                        EnvVar: "ENV_NODAEMON",
                        Name:   "nodaemon",
                        Usage: "do not daemonise",
                    },
        },
	},
    {
        Name:   "delete",
        Usage:  "<container name>",
        Action: command.CmdDelete,
        Flags:  []cli.Flag{},
    },
	{
		Name:   "status",
		Usage:  "<container name>",
		Action: command.CmdStatus,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "stop",
		Usage:  "<container name>",
		Action: command.CmdStop,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "list",
		Usage:  "",
		Action: command.CmdList,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
    flag.PrintDefaults()
	os.Exit(2)
}
