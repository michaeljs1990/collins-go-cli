package main

import (
	cli "github.com/urfave/cli"
)

func provisionSubcommand() cli.Command {
	return cli.Command{
		Name:  "provision",
		Usage: "Provision assets",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
