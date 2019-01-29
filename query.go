package main

import (
	"fmt"
  "strings"

	log "github.com/sirupsen/logrus"
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
				Usage:    "Assets in primary role",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "R, secondary-role",
				Usage:    "Assets in secondary role",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "i, ip-address",
				Usage:    "Assets with IP address[es]",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "S, status",
				Usage:    "Asset status (and optional state after :)",
				Category: "Query options",
			},
			cli.StringFlag{
				Name:     "a, attribute",
				Usage:    "Arbitrary attributes and values to match in query. : between key and value",
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

func queryBuildOptions(c *cli.Context) collins.AssetFindOpts {
	opts := collins.AssetFindOpts{}

  if c.IsSet("status") {
    status := strings.Split(c.String("status"), ":")
    if len(status) == 2 {
      opts.State = status[1]
    }
    opts.Status = status[0]
  }

  if c.IsSet("attribute") {
    attribute := strings.Split(c.String("attribute"), ":")
    opts.Attribute = strings.Join(attribute, ";")
  }

  return opts
}

// Render the output to the screen in the given format
func queryOutputFormat(format string, assets []collins.Asset) {
	for _, asset := range assets {
		fmt.Printf("", asset.Metadata.Tag)
	}
}

func queryRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)
	opts := queryBuildOptions(c)

  var allAssets []collins.Asset
	for {
		assets, resp, err := client.Assets.Find(&opts)

		if err != nil {
			log.Fatal(err.Error())
		}

    allAssets = append(allAssets, assets...)

		if resp.NextPage == resp.CurrentPage { // No more pages
			break
		} else { // Fetch next page
			opts.PageOpts.Page++
		}
	}

  queryOutputFormat("", allAssets)

  return nil
}
