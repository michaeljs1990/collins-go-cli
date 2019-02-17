package cmds

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

// -s, --show-pools                 Show IP pools
// -H, --show-header                Show header fields in --show-pools output
// -a, --allocate POOL              Allocate addresses in POOL
// -n, --number [NUM]               Allocate NUM addresses (Defaults to 1 if omitted)
// -d, --delete [POOL]              Delete addresses in POOL. Deletes ALL addresses if POOL is omitted
// -t, --tags TAG[,...]             Tags to work on, comma separated

func IpamSubcommand() cli.Command {
	return cli.Command{
		Name:    "ipam",
		Aliases: []string{"address", "ipaddress"},
		Usage:   "Allocate and delete IPs, show IP pools",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:     "s, show-pools",
				Usage:    "Show IP pools",
				Category: "IPAM options",
			},
			cli.BoolFlag{
				Name:     "i, show-ipmi-pools",
				Usage:    "Show IPMI pools",
				Category: "IPAM options",
			},
			cli.BoolFlag{
				Name:     "H, show-header",
				Usage:    "Show header fields",
				Category: "IPAM options",
			},
		},
		Action: ipamRunCommand,
	}
}

func fieldToPoolStruct(field string, pool collins.Pool) string {

	switch field {
	case "name":
		return pool.Name
	case "network":
		return pool.Network
	case "start_address":
		return pool.StartAddress
	case "specified_gateway":
		return pool.SpecifiedGateway
	case "gateway":
		return pool.Gateway
	case "broadcast":
		return pool.Broadcast
	case "possible_addresses":
		return strconv.Itoa(pool.PossibleAddresses)
	}

	return ""
}

func renderPools(c *cli.Context, col *collins.Client, pools []collins.Pool) {

	columns := []string{
		"name",
		"network",
		"start_address",
		"specified_gateway",
		"gateway",
		"broadcast",
		"possible_addresses"}

	showHeaders := c.IsSet("show-header")

	// Find the longest column in each field so the final output is pretty.
	maxColumnWidth := make(map[string]int)
	for _, column := range columns {
		var max int
		for _, pool := range pools {
			length := len(fieldToPoolStruct(column, pool))
			if length > max {
				max = length
			}
		}

		if showHeaders && len(column) > max {
			max = len(column)
		}

		maxColumnWidth[column] = max
	}

	// Make sure we build the formatter back in the correct order.
	// Golang you need more datastructures for real.
	var formatterSlice []string
	for _, col := range columns {
		fmtr := "%-" + strconv.Itoa(maxColumnWidth[col]) + "s"
		formatterSlice = append(formatterSlice, fmtr)
	}

	formatter := strings.Join(formatterSlice, "\t")

	if showHeaders {
		headers := make([]interface{}, len(columns))
		for i, v := range columns {
			headers[i] = v
		}

		fmt.Fprintf(os.Stderr, formatter+"\n", headers...)
	}

	for _, pool := range pools {
		// We use an interface instead of a slice becasue Printf requires this.
		var fields []interface{}

		for _, column := range columns {
			fields = append(fields, fieldToPoolStruct(column, pool))
		}

		fmt.Printf(formatter+"\n", fields...)
	}
}

func ipamRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	switch {
	case c.IsSet("show-pools"):
		pools, _, err := client.IPAM.Pools()
		if err != nil {
			logAndDie(err.Error())
		}
		renderPools(c, client, pools)
	case c.IsSet("show-ipmi-pools"):
		pools, _, err := client.IPAM.IpmiPools()
		if err != nil {
			logAndDie(err.Error())
		}
		renderPools(c, client, pools)
	}

	return nil
}
