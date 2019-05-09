package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	color "github.com/logrusorgru/aurora"
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
			cli.IntFlag{
				Name:     "n, number",
				Usage:    "number of operations to run concurrently",
				Category: "Power options",
				Value:    10,
			},
		},
		Action: powerRunCommand,
	}
}

// For the same behavior as the collins-cli command if both power and
// status are passed in we just ignore the status request and run the
// specified power command. Maybe remove this and allow you to specify
// both in a later release.
func powerActionByTag(wg *sync.WaitGroup, ctx *cli.Context, col *collins.Client, tag string) {
	defer wg.Done()

	if ctx.IsSet("power") {

		debugLog("Trying to perform power action '" + ctx.String("power") + "' on " + tag)

		var err error
		action := strings.ToLower(ctx.String("power"))
		msg := tag + " performing " + action + " ... "
		switch action {
		case "reboot", "rebootsoft":
			_, err = col.Management.SoftReboot(tag)
		case "reboothard":
			_, err = col.Management.HardReboot(tag)
		case "poweron", "on":
			_, err = col.Management.PowerOn(tag)
		case "poweroff", "off":
			_, err = col.Management.PowerOff(tag)
		case "identify":
			_, err = col.Management.Identify(tag)
		case "verify":
			_, err = col.Management.Verify(tag)
		default:
			logAndDie("Unknown power action '" + action + "' , expecting one of reboot,rebootsoft,reboothard,on,off,poweron,poweroff,identify,verify")
		}

		// We dont' use printSuccess or printError here since these are happening in a go routine
		// and although it's unlikely it's possibly that these don't print in the right order if it's
		// not all printed to screen in one fmt.Print function.
		if err != nil {
			gotError = true
			fmt.Print(msg, color.Red("ERROR"), "\n")
		} else {
			fmt.Print(msg, color.Green("SUCCESS"), "\n")
		}

		return
	}

	if ctx.IsSet("status") {
		debugLog("Trying to get current power status from " + tag)

		msg := tag + " checking power status ... "

		stat, _, err := col.Management.PowerStatus(tag)

		if stat == "" {
			stat = "Unknown"
		}

		if err != nil {
			gotError = true
		}

		fmt.Print(msg+"(", color.Magenta(stat), ")\n")

		return
	}

}

func powerRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.Int("number") <= 0 {
		logAndDie("--number must be greater than 0")
	}

	var wg sync.WaitGroup
	powerActionsRun := 0
	if c.IsSet("tags") {
		tags := strings.Split(c.String("tags"), ",")
		for _, tag := range tags {
			powerActionsRun++
			wg.Add(1)
			go powerActionByTag(&wg, c, client, tag)
			if 0 == powerActionsRun%c.Int("number") {
				debugLog("hit max number of outstanding power actions waiting for some to finish before continuing")
				wg.Wait()
			}
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
				powerActionsRun++
				wg.Add(1)
				go powerActionByTag(&wg, c, client, tag[0])
			}

			// In the case that the actions run are a mod of number
			// we take a break and wait for them to finish up before
			// issuing more commands.
			if 0 == powerActionsRun%c.Int("number") {
				debugLog("hit max number of outstanding power actions waiting for some to finish before continuing")
				wg.Wait()
			}
		}
	}

	// Ensure nothing is outstanding
	debugLog("Will now wait to ensure that everything finishes before returning")
	wg.Wait()

	if gotError {
		return errors.New("Some commands failed to run to success")
	}

	return nil
}
