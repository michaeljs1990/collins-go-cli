package cmds

import (
	cli "github.com/urfave/cli"
)

func IpamSubcommand() cli.Command {
	return cli.Command{
		Name:    "ipam",
		Aliases: []string{"address", "ipaddress"},
		Usage:   "Allocate and delete IPs, show IP pools",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
