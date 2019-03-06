package cmds

import (
	"bufio"
	"io"
	"os"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func QuerySubcommand() cli.Command {
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
				Usage:    "Arbitrary attributes. ':'  or ':~' between key and value. : is a positive match and :~ is a negative match",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "o, operation",
				Usage:    "Sets if your query will be joined with AND or OR",
				Value:    "AND",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "q, query",
				Usage:    "Specify a CQL query. This will overwrite most other flags see docs for more info.",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "u, pipe",
				Usage:    "This sets the attribute to match against when piping to stdin. When not set it defaults to tags.",
				Category: "Query options",
			},
			cli.IntFlag{
				Name:     "U, pipe-size",
				Usage:    "This sets the number of keys sent to collins at once when piping.",
				Value:    30,
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

func queryBuildOptions(c *cli.Context, hostname string, fromStdin []string) collins.AssetFindOpts {
	opts := collins.AssetFindOpts{}

	if c.IsSet("remote-lookup") {
		opts.RemoteLookup = true
	}

	if c.IsSet("query") {
		opts.Query = c.String("query")
	} else {
		opts.Query = buildOptionsQuery(c, hostname, fromStdin)
	}

	debugLog("CQL executed - " + opts.Query)

	return opts
}

// Build a CQL query from the passed in flag taking into account comma seperated values.
func cqlQuery(c *cli.Context, flagVal string, field string) string {
	query := []string{}
	nmatchNum := 0
	for _, val := range strings.Split(flagVal, ",") {
		// If the first value is a ~ in the field we are doing
		// a negative match.
		nmatch := (val[0] == "~"[0])
		if nmatch && len(val) == 1 {
			logAndDie("You tried to do a negative match but didn't pass in any value to do it on")
		} else if nmatch {
			nmatchNum++
			query = append(query, "("+field+" != "+val[1:]+")")
		} else {
			query = append(query, "("+field+" = "+val+")")
		}
	}

	// This isn't needed but keeps the resulting CQL from not looking so crazy.
	// Additionally if all that we are doing is negative matches it makes no sense
	// to join with OR. In the case that a flag like -n "~devnode,~aaanode" is passed
	// in we return everything that is not a devnode and not an aaanode.
	if len(query) == 1 {
		return query[0]
	} else if nmatchNum == len(query) {
		return "(" + strings.Join(query, " AND ") + ")"
	} else {
		return "(" + strings.Join(query, " OR ") + ")"
	}
}

// This is broke out of build options just for the sake of making testing easier
func buildOptionsQuery(c *cli.Context, hostname string, fromStdin []string) string {
	cql := []string{}

	// The go client isn't as friendly as the ruby one which is fine we will just
	// take everything else and convert it into CQL to talk to collins.
	if len(fromStdin) > 0 {
		// By default we match stdin against a tag unless pipe is set
		// in which the user can set any field to match on.
		key := "TAG"
		if c.IsSet("pipe") {
			key = c.String("pipe")
		}

		query := []string{}
		for _, val := range fromStdin {
			query = append(query, "("+key+" = "+val+")")
		}

		if len(query) == 1 {
			cql = append(cql, query[0])
		} else {
			cql = append(cql, "("+strings.Join(query, " OR ")+")")
		}
	}

	if c.IsSet("status") {
		status := strings.Split(c.String("status"), ":")
		if len(status) == 2 {
			cql = append(cql, cqlQuery(c, status[1], "STATE"))
		}
		cql = append(cql, cqlQuery(c, status[0], "STATUS"))
	}

	if c.IsSet("type") {
		cql = append(cql, cqlQuery(c, c.String("type"), "TYPE"))
	}

	if c.IsSet("tag") {
		cql = append(cql, cqlQuery(c, c.String("tag"), "TAG"))
	}

	if c.IsSet("nodeclass") {
		cql = append(cql, cqlQuery(c, c.String("nodeclass"), "NODECLASS"))
	}

	if c.IsSet("pool") {
		cql = append(cql, cqlQuery(c, c.String("pool"), "POOL"))
	}

	if c.IsSet("role") {
		cql = append(cql, cqlQuery(c, c.String("role"), "PRIMARY_ROLE"))
	}

	if c.IsSet("secondary-role") {
		cql = append(cql, cqlQuery(c, c.String("secondary-role"), "SECONDARY_ROLE"))
	}

	if c.IsSet("ip-address") {
		cql = append(cql, cqlQuery(c, c.String("ip-address"), "IP_ADDRESS"))
	}

	if hostname != "" {
		cql = append(cql, "(HOSTNAME = "+hostname+")")
	}

	if c.IsSet("attribute") || c.IsSet("a") {
		for _, attr := range c.StringSlice("attribute") {
			var attrSplit []string
			equal := false
			attrSplit = strings.SplitN(attr, ":~", 2)
			if len(attrSplit) != 2 {
				attrSplit = strings.SplitN(attr, ":", 2)
				equal = true
			}
			if len(attrSplit) != 2 {
				logAndDie("--attribute and -a requires attribute:value, missing :value")
			}
			attrKey := strings.ToUpper(attrSplit[0])
			attrValue := strings.ToUpper(attrSplit[1])

			if equal {
				cql = append(cql, "("+attrKey+" = "+attrValue+")")
			} else {
				cql = append(cql, "("+attrKey+" != "+attrValue+")")
			}
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
			var attrSplit []string
			attrSplit = strings.SplitN(attr, ":~", 2)
			if len(attrSplit) != 2 {
				attrSplit = strings.SplitN(attr, ":", 2)
			}
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
	// If the user passes in an argument we treat it as a
	// hostname and pass it along it overwrites hostname
	// in the case you set it as an attribute
	hostname := ""
	if c.NArg() > 0 {
		hostname = c.Args().Get(0)
	}

	// The use case for this seems odd but it's actual very helpful to be
	// able to take a list of tags and query to get more info about them.
	// We check here if the command is being piped to or not so that you can
	// specify --tags and pipe to the file.
	fromStdin := []string{}
	fi, err := os.Stdin.Stat()
	if err != nil {
		logAndDie(err.Error())
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
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
				fromStdin = append(fromStdin, tag[0])
			}
		}
	}

	if len(os.Args) == 2 && len(fromStdin) == 0 {
		logAndDie("See --help for collins query usage")
	}

	client := getCollinsClient(c)

	// In the case that we were piped a ton of data from stdin we break it up
	// so we do not trigger an error when querying collins. For more info on
	// this see https://github.com/michaeljs1990/collins-go-cli/issues/28
	opts := []collins.AssetFindOpts{}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		if c.Int("pipe-size") == 0 {
			logAndDie("Setting pipe-size to zero is not valid.")
		}
		batchSize := c.Int("pipe-size")
		var batches [][]string
		for batchSize < len(fromStdin) {
			fromStdin, batches = fromStdin[batchSize:], append(batches, fromStdin[0:batchSize:batchSize])
		}
		batches = append(batches, fromStdin)

		for _, batch := range batches {
			opts = append(opts, queryBuildOptions(c, hostname, batch))
		}
	} else {
		opts = append(opts, queryBuildOptions(c, hostname, fromStdin))
	}

	// Limit and pipe-size don't make sense to use together and will
	// not do what you want so we don't allow it.
	if c.IsSet("limit") && c.IsSet("pipe-size") {
		logAndDie("--limit and --pipe-size can't be set at the same time")
	}

	// Kinda hacky but if limit is set we just set
	// that as the page size and break after the first
	// call to get assets.
	size := c.Int("size")
	if c.IsSet("limit") {
		size = c.Int("limit")
	}

	for i, _ := range opts {
		opts[i].PageOpts = collins.PageOpts{
			Size: size,
		}
	}

	runs := 0
	var allAssets []collins.Asset
	for {
		assets, resp, err := client.Assets.Find(&opts[runs])

		if err != nil {
			logAndDie(err.Error())
		}

		allAssets = append(allAssets, assets...)

		// Limit was set break out of the loop now
		if c.IsSet("limit") {
			break
		}

		if resp.NextPage == resp.CurrentPage && len(opts) == runs+1 { // No more pages
			break
		} else if resp.NextPage == resp.CurrentPage {
			runs++
		} else { // Fetch next page
			opts[runs].PageOpts.Page++
		}
	}

	columns := queryGetColumns(c)
	format := getOutputFormat(c)
	showHeaders := c.Bool("show-header")
	formatAssets(format, c.String("field-separator"), showHeaders, client.BaseURL.String(), c.Bool("remote-lookup"), columns, allAssets)

	return nil
}
