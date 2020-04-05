package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

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
			cli.StringFlag{
				Name:     "a, allocate",
				Usage:    "Allocate addresses in pool",
				Category: "IPAM options",
			},
			cli.IntFlag{
				Name:     "n, number",
				Usage:    "Allocate NUM addresses",
				Value:    1,
				Category: "IPAM options",
			},
			cli.StringFlag{
				Name:     "d, delete",
				Usage:    "Delete addresses in POOL. An empty string deletes all pools",
				Category: "IPAM options",
			},
			cli.StringFlag{
				Name:     "dip, delete-ip",
				Usage:    "Delete a single IP on an asset",
				Category: "IPAM options",
			},
			cli.StringFlag{
				Name:     "t, tags",
				Usage:    "Tags to work on, comma separated",
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

func deleteIP(c *cli.Context, col *collins.Client) {
	ip := c.String("delete-ip")

	opts := collins.AddressDeleteOpts{
		Address: ip,
	}

	// This only ever makes sense to do for one tag but we will
	// leave this here since all other commands support multiple tags
	// it will just throw some errors.
	for _, tag := range strings.Split(c.String("tags"), ",") {
		num, _, err := col.IPAM.Delete(tag, opts)
		if ip == "" {
			fmt.Printf("%s deleting all IPs... ", tag)
		} else {
			fmt.Printf("%s deleting IP %s... ", tag, ip)
		}

		if err != nil {
			printError(err.Error())
		} else {
			printSuccessWithMsg(fmt.Sprintf("Deleted %d IPs", num))
		}
	}
}

func deleteAddress(c *cli.Context, col *collins.Client) {
	pool := c.String("delete")

	opts := collins.AddressDeleteOpts{
		Pool: pool,
	}

	for _, tag := range strings.Split(c.String("tags"), ",") {
		num, _, err := col.IPAM.Delete(tag, opts)
		if pool == "" {
			fmt.Printf("%s deleting all IPs in %s... ", tag, pool)
		} else {
			fmt.Printf("%s deleting all IPs... ", tag)
		}

		if err != nil {
			printError(err.Error())
		} else {
			printSuccessWithMsg(fmt.Sprintf("Deleted %d IPs", num))
		}
	}
}

func allocateAddress(c *cli.Context, col *collins.Client) {
	num := c.Int("number")
	pool := c.String("allocate")

	opts := collins.AddressAllocateOpts{
		Count: num,
		Pool:  pool,
	}

	for _, tag := range strings.Split(c.String("tags"), ",") {
		addrs, _, err := col.IPAM.Allocate(tag, opts)
		fmt.Printf("%s allocating %d IP in %s... ", tag, num, pool)
		if err != nil {
			printError(err.Error())
		} else {
			msg := []string{"Allocated"}
			for _, addr := range addrs {
				msg = append(msg, addr.Address)
			}
			printSuccessWithMsg(strings.Join(msg, " "))
		}
	}
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
		pools, _, err := client.IPAM.IPMIPools()
		if err != nil {
			logAndDie(err.Error())
		}
		renderPools(c, client, pools)
	case c.IsSet("allocate"):
		allocateAddress(c, client)
	case c.IsSet("delete"):
		deleteAddress(c, client)
	case c.IsSet("delete-ip"):
		deleteIP(c, client)
	}

	return nil
}
