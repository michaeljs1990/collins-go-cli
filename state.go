package main

import (
	cli "github.com/urfave/cli"
)

func stateSubcommand() cli.Command {
	return cli.Command{
		Name:    "state",
		Aliases: []string{"status"},
		Usage:   "Show and manage states and statuses via State API",
		Action: func(c *cli.Context) error {
			return nil
		},
	}

}
