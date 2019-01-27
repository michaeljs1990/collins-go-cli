package main

import (
	cli "github.com/urfave/cli"
)

func datacenterSubcommand() cli.Command {
	return cli.Command{
		Name:    "datacenter",
		Aliases: []string{"dc"},
		Usage:   "Manage multiple Collins configurations",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
