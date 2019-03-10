package commands

import (
	cli "github.com/urfave/cli"
)

// This handles any setup that needs to happen before every subcommand
// is called.
func BeforeMiddleware(c *cli.Context) error {
	debugMode = c.GlobalIsSet("debug")

	return nil
}
