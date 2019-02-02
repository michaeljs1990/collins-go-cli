package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func logSubcommand() cli.Command {
	return cli.Command{
		Name:  "log",
		Usage: "Display log messages on assets",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:     "a, all",
				Usage:    "Show logs from ALL assets",
				Category: "Log options",
			},
			cli.IntFlag{
				Name:     "n, number",
				Usage:    "Show the last LINES log entries",
				Value:    20,
				Category: "Log options",
			},
			cli.BoolFlag{
				Name:     "f, follow",
				Usage:    "Poll for logs every 2 seconds",
				Category: "Log options",
			},
			cli.StringFlag{
				Name:     "s, severity",
				Usage:    "Separator between columns in output",
				Value:    "all",
				Category: "Log options",
			},
			cli.StringFlag{
				Name:     "t, tags",
				Usage:    "Assets with tag[s] value[,...]",
				Category: "Log options",
			},
		},
		Action: logRunCommand,
	}
}

func fieldToLogStruct(field string, log collins.Log) string {
	switch field {
	case "time":
		return log.Created
	case "creator":
		return log.CreatedBy
	case "severity":
		return log.Type
	case "tag":
		return log.AssetTag
	}

	return ""
}

func logGetFormat(columns []string, logs []collins.Log) string {
	maxColumnWidth := make(map[string]int)
	for _, column := range columns {
		var max int
		for _, log := range logs {
			length := len(fieldToLogStruct(column, log))
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

	formatter := strings.Join(formatterSlice, " ")

	// Add on last string for the message
	return formatter + " %s"
}

func handleLogs(c *cli.Context, col *collins.Client, tags []string) {
	opts := collins.LogGetOpts{}
	opts.PageOpts = collins.PageOpts{
		Size: c.Int("number"),
	}

	uniqueSet := UniqueOrderedSet{
		"time",
		"creator",
		"serverity",
		"tag",
	}

	seenLogIDs := map[int]bool{}
	for {
		logsThisRun := []collins.Log{}
		var resp *collins.Response
		for _, tag := range tags {
			logs, r, err := col.Logs.Get(tag, &opts)
			resp = r
			if err != nil {
				fmt.Println("Unable to fetch logs for " + tag + ": " + err.Error())
			}

			logsThisRun = append(logsThisRun, logs...)
		}

		// Get the format for logs this run, sort them, and print
		format := logGetFormat(uniqueSet, logsThisRun)
		sort.Slice(logsThisRun, func(i, j int) bool {
			return logsThisRun[i].ID < logsThisRun[j].ID
		})

		for _, log := range logsThisRun {
			var fields []interface{}

			for _, column := range uniqueSet {
				fields = append(fields, fieldToLogStruct(column, log))
			}

			fields = append(fields, log.Message)

			fmt.Printf(format+"\n", fields...)
			seenLogIDs[log.ID] = true
		}

		if !c.IsSet("follow") {
			break
		}

		if resp.NextPage == resp.CurrentPage { // No more pages
			break
		} else { // Fetch next page
			opts.PageOpts.Page++
		}
	}
}

func logRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("tags") {
		tags := strings.Split(c.String("tags"), ",")
		handleLogs(c, client, tags)
	}

	return nil
}
