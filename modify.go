package main

import (
	cli "github.com/urfave/cli"
)

func modifySubcommand() cli.Command {
	return cli.Command{
		Name:    "modify",
		Aliases: []string{"set"},
		Usage:   "Add and remove attributes, change statuses, and log to assets",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
