package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

var gotError = false

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

func attributeDeleteStrings(ctx *cli.Context) []string {
	if ctx.IsSet("d") || ctx.IsSet("delete-attribute") {
		return ctx.StringSlice("delete-attribute")
	}

	return []string{}
}

func statusUpdateOpts(ctx *cli.Context) collins.AssetUpdateStatusOpts {
	opt := collins.AssetUpdateStatusOpts{}

	if ctx.IsSet("set-state") {
		if !ctx.IsSet("reason") || ctx.String("reason") == "" {
			log.Fatal("You need to provide a --reason when changing asset states!")
		}

		status := strings.Split(ctx.String("set-state"), ":")
		if len(status) == 2 {
			opt.State = status[1]
		}
		opt.Status = status[0]

		opt.Reason = ctx.String("reason")
	}

	return opt
}

func logCreateOpts(ctx *cli.Context) collins.LogCreateOpts {
	opts := collins.LogCreateOpts{}

	if ctx.IsSet("log") && ctx.String("log") == "" {
		log.Fatal("You need to provide a message with the --log flag")
	} else if ctx.IsSet("log") {
		opts.Message = ctx.String("log")
		opts.Type = "NOTE"
	}

	validLevels := []string{
		"ERROR",
		"DEBUG",
		"EMERGENCY",
		"ALERT",
		"CRITICAL",
		"WARNING",
		"NOTICE",
		"INFORMATIONAL",
		"NOTE"}

	// Only set the level if log is also set otherwise we will
	// just ignore this option from the user.
	if ctx.IsSet("level") && ctx.IsSet("log") {
		valid := false
		level := strings.ToUpper(ctx.String("level"))
		for _, vlevel := range validLevels {
			if level == vlevel {
				valid = true
			}
		}

		if !valid {
			log.Fatal("Your log level is not valid")
		}

		opts.Type = level
	}

	return opts
}

// Run everything that will mutate the assets state in this function
func modifyAssetByTag(ctx *cli.Context, col *collins.Client, tag string) {
	// Generate all the options before doing anything so we don't half start
	// applying settings and then run into an issue with the proper flags not
	// being passed in such as status without reason.
	attrs := attributeUpdateOpts(ctx)
	status := statusUpdateOpts(ctx)
	delattrs := attributeDeleteStrings(ctx)
	logMsg := logCreateOpts(ctx)

	// Apply the options that we have set and try to output it in some kind
	// of sane format for users to see what applied and what did not.
	for _, attr := range attrs {
		_, err := col.Assets.Update(tag, &attr)
		attrSplit := strings.SplitN(attr.Attribute, ";", 2)
		msg := tag + " setting " + strings.Join(attrSplit, "=")
		if err != nil {
			gotError = true
			log.Error(msg)
		} else {
			log.Print(msg)
		}
	}

	for _, attr := range delattrs {
		_, err := col.Assets.DeleteAttribute(tag, attr)
		msg := tag + " deleting " + attr
		if err != nil {
			gotError = true
			log.Error(msg)
		} else {
			log.Print(msg)
		}
	}

	if status != (collins.AssetUpdateStatusOpts{}) {
		_, err := col.Assets.UpdateStatus(tag, &status)
		msg := tag + " changing status to " + strings.ToUpper(status.Status)
		if status.State != "" {
			msg = msg + ":" + strings.ToUpper(status.State)
		}
		if err != nil {
			gotError = true
			log.Error(msg)
		} else {
			log.Print(msg)
		}
	}

	if logMsg != (collins.LogCreateOpts{}) {
		_, _, err := col.Logs.Create(tag, &logMsg)
		msg := tag + " logging " + strings.ToLower(logMsg.Type) + "\"" + logMsg.Message + "\""
		if err != nil {
			log.Error(msg)
		} else {
			log.Print(msg)
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

	if gotError {
		return errors.New("Some commands failed to run to success")
	}

	return nil
}
