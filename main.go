package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func main() {
	_, err := collins.NewClientFromYaml()
	if err != nil {
		log.Info("You can use COLLINS_CLIENT_CONFIG env to set the location of your config")
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "collins"
	app.Usage = "Interface with http://tumblr.github.io/collins/"

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

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
