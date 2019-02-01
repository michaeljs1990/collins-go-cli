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

func powerSubcommand() cli.Command {
	return cli.Command{
		Name:  "power",
		Usage: "Control and show power status",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:     "s, status",
				Usage:    "Show IPMI power status",
				Category: "Power options",
			},
			cli.StringFlag{
				Name:     "p, power",
				Usage:    "Perform IPMI power ACTION",
				Category: "Power options",
			},
			cli.StringFlag{
				Name:     "t, tags",
				Usage:    "Tags to work on, comma separated",
				Category: "Power options",
			},
		},

		Action: powerRunCommand,
	}
}

// For the same behavior as the collins-cli command if both power and
// status are passed in we just ignore the status request and run the
// specified power command. Maybe remove this and allow you to specify
// both in a later release.
func powerActionByTag(ctx *cli.Context, col *collins.Client, tag string) {

	if ctx.IsSet("power") {
		var err error
		switch ctx.String("power") {
		case "reboot", "rebootSoft":
			_, err = col.Management.SoftReboot(tag)
		case "reboothard":
			_, err = col.Management.HardReboot(tag)
		case "poweron", "on":
			_, err = col.Management.PowerOn(tag)
		case "poweroff", "off":
			_, err = col.Management.SoftPowerOff(tag)
		case "identify":
			_, err = col.Management.Identify(tag)
		case "verify":
			_, err = col.Management.Verify(tag)
		default:
			log.Fatal("Unknown power action rebootx, expecting one of reboot,rebootsoft,reboothard,on,off,poweron,poweroff,identify,verify")
		}

		msg := tag + " performing " + ctx.String("power") + " ..."
		if err != nil {
			gotError = true
			log.Error(msg)
		} else {
			log.Print(msg)
		}

		return
	}

	if ctx.IsSet("status") {
		stat, _, err := col.Management.PowerStatus(tag)

		if stat == "" {
			stat = "Unknown"
		}

		msg := tag + " checking power status ... (" + stat + ")"
		if err != nil {
			gotError = true
			log.Error(msg)
		} else {
			log.Print(msg)
		}

		return
	}

}

func powerRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("tags") {
		tags := strings.Split(c.String("tags"), ",")
		for _, tag := range tags {
			powerActionByTag(c, client, tag)
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

			// If a newline was all that was recieved from stdin
			// ignore it and keep going.
			tag := strings.Fields(line)
			if len(tag) >= 1 {
				powerActionByTag(c, client, tag[0])
			}
		}
	}

	if gotError {
		return errors.New("Some commands failed to run to success")
	}

	return nil
}
