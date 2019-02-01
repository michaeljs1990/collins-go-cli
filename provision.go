package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func provisionSubcommand() cli.Command {
	return cli.Command{
		Name:  "provision",
		Usage: "Provision assets",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "n, nodeclass",
				Usage:    "Nodeclass to provision as. (Required)",
				Category: "Provision options",
			},
			cli.StringFlag{
				Name:     "p, pool",
				Usage:    "Provision with pool POOL",
				Category: "Provision options",
			},
			cli.StringFlag{
				Name:     "r, role",
				Usage:    "Provision with primary role ROLE",
				Category: "Provision options",
			},
			cli.StringFlag{
				Name:     "R, secondary-role",
				Usage:    "Provision with secondary role ROLE",
				Category: "Provision options",
			},
			cli.StringFlag{
				Name:     "s, suffix",
				Usage:    "Provision with suffix SUFFIX",
				Category: "Provision options",
			},
			cli.StringFlag{
				Name:     "b, build-contact",
				Usage:    "Build contact",
				Value:    os.Getenv("USER"),
				Category: "Provision options",
			},
			cli.StringFlag{
				Name:     "t, tags",
				Usage:    "Tags to work on, comma separated",
				Category: "Provision options",
			},
		},
		Action: provisionRunCommand,
	}
}

// This is a little funny with the go collins api it taks an optional opts param
// but it overrides half of this struct so we just fill the half it doesn't
func provisionMakeOpts(ctx *cli.Context, tag string) collins.ProvisionOpts {

	return collins.ProvisionOpts{
		Suffix:        ctx.String("suffix"),
		PrimaryRole:   ctx.String("role"),
		SecondaryRole: ctx.String("secondary-role"),
		Pool:          ctx.String("pool"),
	}

}

func provisionByTag(ctx *cli.Context, col *collins.Client, tag string) {
	if !ctx.IsSet("nodeclass") {
		logAndDie("You need to specify at least a nodeclass when provisioning")
	}

	opts := provisionMakeOpts(ctx, tag)
	profile := ctx.String("nodeclass")
	contact := ctx.String("build-contact")
	msg := tag + " provisioning with nodeclass:" + profile + " by " + contact + "... "
	fmt.Print(msg)

	_, err := col.Management.Provision(tag, profile, contact, opts)
	if err != nil {
		gotError = true
		printError(err.Error())
	} else {
		printSuccess()
	}
}

func provisionRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("tags") {
		tags := strings.Split(c.String("tags"), ",")
		for _, tag := range tags {
			provisionByTag(c, client, tag)
		}
	} else {
		// No tag was passed in try to read from stdin
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				logAndDie(err.Error())
			}

			// If a newline was all that was recieved from stdin
			// ignore it and keep going.
			tag := strings.Fields(line)
			if len(tag) >= 1 {
				provisionByTag(c, client, tag[0])
			}
		}
	}

	if gotError {
		return errors.New("Some commands failed to run to success")
	}

	return nil
}
