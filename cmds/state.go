package cmds

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func StateSubcommand() cli.Command {
	return cli.Command{
		Name:    "state",
		Aliases: []string{"status"},
		Usage:   "Show and manage states and statuses via State API",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:     "l, list",
				Usage:    "Show IPMI power status",
				Category: "State options",
			},
			cli.BoolFlag{
				Name:     "H, show-header",
				Usage:    "Show header fields in output",
				Category: "Table formatting",
			},
			cli.StringFlag{
				Name:     "f, field-separator",
				Usage:    "Separator between columns in output",
				Category: "Table formatting",
			},
		},
		Action: stateRunCommand,
	}
}

func fieldToStateStruct(field string, state collins.State) string {
	switch field {
	case "status_name":
		if state.Status.Name == "" {
			return "Any"
		}

		return state.Status.Name
	case "state_name":
		return state.Name
	case "status_description":
		if state.Status.Name == "" {
			return "Any status"
		}

		return state.Status.Description
	case "description":
		return state.Description
	}

	return ""
}

func stateRenderTable(columns []string, showHeaders bool, states []collins.State) {
	// Find the longest column in each field so the final output is pretty.
	maxColumnWidth := make(map[string]int)
	for _, column := range columns {
		var max int
		for _, state := range states {
			length := len(fieldToStateStruct(column, state))
			if length > max {
				max = length
			}
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

	for _, state := range states {
		// We use an interface instead of a slice becasue Printf requires this.
		var fields []interface{}

		for _, column := range columns {
			fields = append(fields, fieldToStateStruct(column, state))
		}

		fmt.Printf(formatter+"\n", fields...)
	}
}

func stateRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("list") {
		columns := []string{
			"status_name",
			"state_name",
			"status_description",
			"description"}

		headers := c.IsSet("show-header")
		states, _, err := client.States.List()
		if err != nil {
			logAndDie(err.Error())
		}

		stateRenderTable(columns, headers, states)
	}

	return nil
}
