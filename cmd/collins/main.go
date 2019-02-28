package main

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli"

	src "cgit.xrt0x.com/xrt0x/collins-go-cli/src"
)

var (
	version = "master"
	commit  = "4b825dc642cb6eb9a060e54bf8d69288fbee4904" // Empty tree hash
)

func main() {
	app := cli.NewApp()
	app.Name = "collins"
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Usage = "Interface with http://tumblr.github.io/collins/"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:     "timeout",
			Usage:    "Timeout in seconds (0 == forever)",
			Category: "Global setting",
		},
		cli.BoolFlag{
			Name:     "debug",
			Usage:    "Print some extra info to stderr when running commands",
			Category: "Global setting",
		},
		cli.StringFlag{
			Name:     "C, config",
			Usage:    "Use specific Collins config yaml for client",
			Category: "Global setting",
		},
	}

	// Setup any needed middleware
	app.Before = src.BeforeMiddleware

	subCmds := []cli.Command{}
	subCmds = append(subCmds, src.QuerySubcommand())
	subCmds = append(subCmds, src.ModifySubcommand())
	subCmds = append(subCmds, src.LogSubcommand())
	subCmds = append(subCmds, src.ProvisionSubcommand())
	subCmds = append(subCmds, src.PowerSubcommand())
	subCmds = append(subCmds, src.IpamSubcommand())
	subCmds = append(subCmds, src.StateSubcommand())
	subCmds = append(subCmds, src.DatacenterSubcommand())
	app.Commands = subCmds

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
