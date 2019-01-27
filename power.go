package main

import (
	cli "github.com/urfave/cli"
)

func powerSubcommand() cli.Command {
	return cli.Command{
		Name:    "power",
		Usage:   "Control and show power status",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
