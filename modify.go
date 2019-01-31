package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
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
			cli.StringSliceFlag{
				Name:     "a, set-attribute",
				Usage:    "Set attribute=value. : between key and value. attribute will be uppercased",
				Category: "Modify options",
			},
			cli.StringSliceFlag{
				Name:     "d, delete-attribute",
				Usage:    "Delete attribute",
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

func attributeUpdateOpts(ctx *cli.Context) []collins.AssetUpdateOpts {
	var opts []collins.AssetUpdateOpts
	if ctx.IsSet("a") || ctx.IsSet("set-attribute") {
		// Merge these because of a bug in the CLI library we are using
		attrs := ctx.StringSlice("set-attribute")
		for _, attr := range attrs {
			attrSplit := strings.SplitN(attr, ":", 2)
			if len(attrSplit) != 2 {
				log.Fatal("--set-attribute and -a requires attribute:value, missing :value")
			}

			attrJoin := strings.Join(attrSplit, ";")
			opts = append(opts, collins.AssetUpdateOpts{
				Attribute: attrJoin,
			})
		}
	}

	return opts
}

func modifyAssetByTag(ctx *cli.Context, col *collins.Client, tag string) {
	attrs := attributeUpdateOpts(ctx)
	for _, attr := range attrs {
		_, err := col.Assets.Update(tag, &attr)
		attrSplit := strings.SplitN(attr.Attribute, ";", 2)
		attrMsg := strings.Join(attrSplit, "=")
		if err != nil {
			log.Error(tag + " setting " + attrMsg)
		} else {
			log.Print(tag + " setting " + attrMsg)
		}
	}
}

func modifyRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("tags") {
		tags := strings.Split(c.String("tags"), ",")
		for _, tag := range tags {
			modifyAssetByTag(c, client, tag)
		}
	} else {
		// No tag was passed in try to read from stdin
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			modifyAssetByTag(c, client, line)
		}
	}

	return nil
}
