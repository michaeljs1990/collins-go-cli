package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func getCollinsClient(c *cli.Context) *collins.Client {
	collins, err := collins.NewClientFromYaml()
	if err != nil {
		log.Info("You can use COLLINS_CLIENT_CONFIG env to set the location of your config")
		log.Fatal(err)
	}
	return collins
}

func main() {
	app := cli.NewApp()
	app.Name = "collins"
	app.Version = "0.0.1"
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

	app.Run(os.Args)
}
