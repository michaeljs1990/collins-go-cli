package main

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

var (
	version = "master"
	commit  = "4b825dc642cb6eb9a060e54bf8d69288fbee4904" // Empty tree hash
)

func getCollinsClient(c *cli.Context) *collins.Client {
	collins, err := collins.NewClientFromYaml()
	if err != nil {
		fmt.Println("You can use COLLINS_CLIENT_CONFIG env to set the location of your config")
		logAndDie(err.Error())
	}
	return collins
}

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
		cli.StringFlag{
			Name:     "C, config",
			Usage:    "Use specific Collins config yaml for client",
			Category: "Global setting",
		},
	}

	cmds := []cli.Command{}
	cmds = append(cmds, querySubcommand())
	cmds = append(cmds, modifySubcommand())
	cmds = append(cmds, logSubcommand())
	cmds = append(cmds, provisionSubcommand())
	cmds = append(cmds, powerSubcommand())
	cmds = append(cmds, ipamSubcommand())
	cmds = append(cmds, stateSubcommand())
	cmds = append(cmds, datacenterSubcommand())
	app.Commands = cmds

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
