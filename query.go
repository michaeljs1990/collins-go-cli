package main

import (
	"os"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func querySubcommand() cli.Command {
	return cli.Command{
		Name:    "query",
		Aliases: []string{"find"},
		Usage:   "Search for assets in Collins",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "t, tag",
				Usage:    "Assets with tag[s] value[,...]",
				Category: "Query options",
			},
			cli.BoolFlag{
				Name:     "Z, remote-lookup",
				Usage:    "Query remote datacenters for asset",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "T, type",
				Usage:    "Only show asset with type value",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "n, nodeclass",
				Usage:    "Assets in nodeclass value[,...]",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "p, pool",
				Usage:    "Assets in pool value[,...]",
				Category: "Query options",
			},
			cli.IntFlag{
				Name:     "s, size",
				Usage:    "Number of assets to return per page",
				Value:    100,
				Category: "Query options",
			},
			cli.IntFlag{
				Name:     "limit",
				Usage:    "Limit total results of assets",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "r, role",
				Usage:    "Assets in primary role[s] value[,...]",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "R, secondary-role",
				Usage:    "Assets in secondary role[s] value[,...]",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "i, ip-address",
				Usage:    "Assets with IP address[es] value[,...]",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "S, status",
				Usage:    "Asset status (and optional state after :)",
				Category: "Query options",
			},
			cli.StringSliceFlag{
				Name:     "a, attribute",
				Usage:    "Arbitrary attributes and values to match in query. : between key and value",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "o, operation",
				Usage:    "Sets if your query will be joined with AND or OR",
				Value:    "AND",
				Category: "Query options",
			},
			cli.BoolFlag{
				Name:     "H, show-header",
				Usage:    "Show header fields in output",
				Category: "Table formatting",
			},
			cli.StringFlag{
				Name:     "c, columns",
				Usage:    "Attributes to output as columns, comma separated",
				Value:    "tag,hostname,nodeclass,status,pool,primary_role,secondary_role",
				Category: "Table formatting",
			},
			cli.StringFlag{
				Name:     "x, extra-columns",
				Usage:    "Show these columns in addition to the default columns, comma separated",
				Category: "Table formatting",
			},
			cli.StringFlag{
				Name:     "f, field-separator",
				Usage:    "Separator between columns in output",
				Value:    "\t",
				Category: "Table formatting",
			},
			cli.BoolFlag{
				Name:     "l, link",
				Usage:    "Output link to assets found in web UI",
				Category: "Robot formatting",
			},
			cli.BoolFlag{
				Name:     "j, json",
				Usage:    "Output results in JSON",
				Category: "Robot formatting",
			},
			cli.BoolFlag{
				Name:     "y, yaml",
				Usage:    "Output results in YAML",
				Category: "Robot formatting",
			},
		},
		Action: queryRunCommand,
	}
}

func queryBuildOptions(c *cli.Context, hostname string) collins.AssetFindOpts {
	opts := collins.AssetFindOpts{}

	if c.IsSet("status") {
		status := strings.Split(c.String("status"), ":")
		if len(status) == 2 {
			opts.State = status[1]
		}
		opts.Status = status[0]
	}

	if c.IsSet("type") {
		opts.Type = c.String("type")
	}

	if c.IsSet("remote-lookup") {
		opts.RemoteLookup = true
	}

	opts.Query = buildOptionsQuery(c, hostname)

	return opts
}

// Build a CQL query from the passed in flag taking into account comma seperated values.
func cqlQuery(c *cli.Context, flag string, field string) string {
	query := []string{}
	for _, val := range strings.Split(c.String(flag), ",") {
		query = append(query, "("+field+" = "+val+")")
	}

	// This isn't needed but keeps the resulting CQL from not looking so crazy
	if len(query) == 1 {
		return query[0]
	} else {
		return "(" + strings.Join(query, " OR ") + ")"
	}
}

// This is broke out of build options just for the sake of making testing easier
func buildOptionsQuery(c *cli.Context, hostname string) string {
	cql := []string{}
	// The go client isn't as friendly as the ruby one which is fine we will just
	// take everything else and convert it into CQL to talk to collins.
	if c.IsSet("tag") {
		cql = append(cql, cqlQuery(c, "tag", "TAG"))
	}

	if c.IsSet("nodeclass") {
		cql = append(cql, cqlQuery(c, "nodeclass", "NODECLASS"))
	}

	if c.IsSet("pool") {
		cql = append(cql, cqlQuery(c, "pool", "POOL"))
	}

	if c.IsSet("role") {
		cql = append(cql, cqlQuery(c, "role", "PRIMARY_ROLE"))
	}

	if c.IsSet("secondary-role") {
		cql = append(cql, cqlQuery(c, "secondary-role", "SECONDARY_ROLE"))
	}

	if c.IsSet("ip-address") {
		cql = append(cql, cqlQuery(c, "ip-address", "IP_ADDRESS"))
	}

	if hostname != "" {
		cql = append(cql, "(HOSTNAME = "+hostname+")")
	}

	if c.IsSet("attribute") || c.IsSet("a") {
		for _, attr := range c.StringSlice("attribute") {
			attrSplit := strings.SplitN(attr, ":", 2)
			if len(attrSplit) != 2 {
				logAndDie("--attribute and -a requires attribute:value, missing :value")
			}
			attrKey := strings.ToUpper(attrSplit[0])
			attrValue := strings.ToUpper(attrSplit[1])

			cql = append(cql, "("+attrKey+" = "+attrValue+")")
		}
	}

	operation := c.String("operation")
	if operation != "AND" && operation != "OR" {
		logAndDie("Operation (or o) flag may only be set to AND or OR")
	}

	return strings.Join(cql, " "+operation+" ")
}

// This uses a "trick" of using a map to create a unique list that
// we then turn into a slice before returning.
func queryGetColumns(c *cli.Context) []string {
	uniqueSet := UniqueOrderedSet{
		"tag",
		"hostname",
		"nodeclass",
		"status",
		"pool",
		"primary_role",
		"secondary_role",
	}

	if c.IsSet("attribute") || c.IsSet("a") {
		for _, attr := range c.StringSlice("attribute") {
			attrSplit := strings.SplitN(attr, ":", 2)
			if len(attrSplit) != 2 {
				logAndDie("--attribute and -a requires attribute:value, missing :value")
			}
			uniqueSet = uniqueSet.Add(attrSplit[0])
		}
	}

	if c.IsSet("columns") {
		uniqueSet = strings.Split(c.String("columns"), ",")
	}

	if c.IsSet("extra-columns") {
		extras := strings.Split(c.String("extra-columns"), ",")
		for _, column := range extras {
			uniqueSet = uniqueSet.Add(column)
		}
	}

	return uniqueSet
}

func getOutputFormat(c *cli.Context) string {
	switch {
	case c.IsSet("link"):
		return "link"
	case c.IsSet("json"):
		return "json"
	case c.IsSet("yaml"):
		return "yaml"
	default:
		return "table"
	}
}

func queryRunCommand(c *cli.Context) error {
	// Don't run if nothing is passed into the command
	if len(os.Args) == 2 {
		logAndDie("See --help for collins query usage")
	}

	// If the user passes in an argument we treat it as a
	// hostname and pass it along it overwrites hostname
	// in the case you set it as an attribute
	hostname := ""
	if c.NArg() > 0 {
		hostname = c.Args().Get(0)
	}

	client := getCollinsClient(c)
	opts := queryBuildOptions(c, hostname)
	size := c.Int("size")

	// Kinda hacky but if limit is set we just set
	// that as the page size and break after the first
	// call to get assets.
	if c.IsSet("limit") {
		size = c.Int("limit")
	}

	opts.PageOpts = collins.PageOpts{
		Size: size,
	}

	var allAssets []collins.Asset
	for {
		assets, resp, err := client.Assets.Find(&opts)

		if err != nil {
			logAndDie(err.Error())
		}

		allAssets = append(allAssets, assets...)

		// Limit was set break out of the loop now
		if c.IsSet("limit") {
			break
		}

		if resp.NextPage == resp.CurrentPage { // No more pages
			break
		} else { // Fetch next page
			opts.PageOpts.Page++
		}
	}

	columns := queryGetColumns(c)
	format := getOutputFormat(c)
	showHeaders := c.Bool("show-header")
	formatAssets(format, c.String("field-separator"), showHeaders, client.BaseURL.String(), c.Bool("remote-lookup"), columns, allAssets)

	return nil
}
