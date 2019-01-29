package main

import (
	cli "github.com/urfave/cli"
)

func logSubcommand() cli.Command {
	return cli.Command{
		Name:  "log",
		Usage: "Display log messages on assets",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
