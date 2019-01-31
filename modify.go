package main

import (
	cli "github.com/urfave/cli"
)

// Allowed values (uppercase or lowercase is accepted):
//   Status:State (-S,--set-status):
//     See `collins state --list`
//   Log levels (-L,--level):
//     ERROR, DEBUG, EMERGENCY, ALERT, CRITICAL, WARNING, NOTICE, INFORMATIONAL, NOTE

func modifySubcommand() cli.Command {
	return cli.Command{
		Name:    "modify",
		Aliases: []string{"set"},
		Usage:   "Add and remove attributes, change statuses, and log to assets",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "a, set-attribute",
				Usage:    "Set attribute=value. : between key and value. attribute will be uppercased",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "d, delete-attribute",
				Usage:    "Delete attribute",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "S, set-state",
				Usage:    "Set attribute=value. : between key and value. attribute will be uppercased",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "S, set-state",
				Usage:    "Set status (and optionally state) to status:state. Requires --reason",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "r, reason",
				Usage:    "Reason for changing status/state",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "l, log",
				Usage:    "Create a log entry",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "L, level",
				Usage:    "Set log level. Default level is NOTE.",
				Value:    "NOTE",
				Category: "Modify options",
			},
			cli.StringFlag{
				Name:     "t, tags",
				Usage:    "Tags to work on, comma separated",
				Category: "Modify options",
			},
		},
		Action: modifyRunCommand,
	}
}

func modifyRunCommand(c *cli.Context) error {
	// client := getCollinsClient(c)
	// assetService := client.AssetServices
	return nil
}
