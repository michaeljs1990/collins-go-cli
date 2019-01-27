package main

import (
	cli "github.com/urfave/cli"
)

func querySubcommand() cli.Command {
	return cli.Command{
		Name:    "query",
		Aliases: []string{"find"},
		Usage:   "Search for assets in Collins",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
