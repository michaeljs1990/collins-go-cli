package cmds

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

func PowerSubcommand() cli.Command {
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
			fmt.Print(tag + " performing " + ctx.String("power") + " ...")
			_, err = col.Management.SoftReboot(tag)
		case "reboothard":
			fmt.Print(tag + " performing " + ctx.String("power") + " ...")
			_, err = col.Management.HardReboot(tag)
		case "poweron", "on":
			fmt.Print(tag + " performing " + ctx.String("power") + " ...")
			_, err = col.Management.PowerOn(tag)
		case "poweroff", "off":
			fmt.Print(tag + " performing " + ctx.String("power") + " ...")
			_, err = col.Management.SoftPowerOff(tag)
		case "identify":
			fmt.Print(tag + " performing " + ctx.String("power") + " ...")
			_, err = col.Management.Identify(tag)
		case "verify":
			fmt.Print(tag + " performing " + ctx.String("power") + " ...")
			_, err = col.Management.Verify(tag)
		default:
			logAndDie("Unknown power action rebootx, expecting one of reboot,rebootsoft,reboothard,on,off,poweron,poweroff,identify,verify")
		}

		if err != nil {
			gotError = true
			printError(err.Error())
		} else {
			printSuccess()
		}

		return
	}

	if ctx.IsSet("status") {
		msg := tag + " checking power status ... "
		fmt.Print(msg)

		stat, _, err := col.Management.PowerStatus(tag)

		if stat == "" {
			stat = "Unknown"
		}

		if err != nil {
			gotError = true
		}

		fmt.Print("(" + stat + ")\n")

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
				logAndDie(err.Error())
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
