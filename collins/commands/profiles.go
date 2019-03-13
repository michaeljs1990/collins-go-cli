package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func ProfilesSubcommand() cli.Command {
	return cli.Command{
		Name:  "profiles",
		Usage: "Show all currently configured collins profiles",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:     "H, show-header",
				Usage:    "Show header fields",
				Category: "Profile options",
			},
			cli.BoolFlag{
				Name:     "l, list",
				Usage:    "List all profiles",
				Category: "Profile options",
			},
		},
		Action: profilesRunCommand,
	}
}

func fieldToProfilesStruct(field string, profile collins.Profile) string {

	switch field {
	case "profile":
		return profile.Profile
	case "label":
		return profile.Label
	case "prefix":
		return profile.Prefix
	case "suffix_allowed":
		return strconv.FormatBool(profile.SuffixAllowed)
	case "primary_role":
		return profile.PrimaryRole
	case "requires_primary_role":
		return strconv.FormatBool(profile.RequiresPrimaryRole)
	case "pool":
		return profile.Pool
	case "requires_pool":
		return strconv.FormatBool(profile.RequiresPool)
	case "secondary_role":
		return profile.SecondaryRole
	case "requires_secondary_role":
		return strconv.FormatBool(profile.RequiresSecondaryRole)
	}

	return ""

}

func renderProfiles(c *cli.Context, profiles []collins.Profile) {

	columns := []string{
		"profile",
		"label",
		"prefix",
		"suffix_allowed",
		"primary_role",
		"requires_primary_role",
		"pool",
		"requires_pool",
		"secondary_role",
		"requires_secondary_role"}

	showHeaders := c.IsSet("show-header")

	// Find the longest column in each field so the final output is pretty.
	maxColumnWidth := make(map[string]int)
	for _, column := range columns {
		var max int
		for _, profile := range profiles {
			length := len(fieldToProfilesStruct(column, profile))
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

	for _, profile := range profiles {
		// We use an interface instead of a slice becasue Printf requires this.
		var fields []interface{}

		for _, column := range columns {
			fields = append(fields, fieldToProfilesStruct(column, profile))
		}

		fmt.Printf(formatter+"\n", fields...)
	}
}

func profilesRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("list") {
		profs, _, err := client.Management.GetProvisioningProfiles()
		if err != nil {
			logAndDie(err.Error())
		}
		renderProfiles(c, profs)
	}

	return nil
}
