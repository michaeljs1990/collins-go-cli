package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	color "github.com/logrusorgru/aurora"
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

// Given a string name give it the proper color scheme
func colorType(sev string) string {
	switch sev {
	case "EMERGENCY":
		return color.Magenta(sev).String()
	case "ALERT", "ERROR":
		return color.Red(sev).String()
	case "INFORMATIONAL":
		return color.Green(sev).String()
	case "DEBUG":
		return color.Blue(sev).String()
	case "NOTE":
		return color.Cyan(sev).String()
	case "CRITICAL":
		return color.Blue(sev).String()
	case "WARNING":
		return color.Brown(sev).String()
	default:
		return color.Gray(sev).String()
	}
}

func fieldToLogStruct(field string, log collins.Log) string {
	switch field {
	case "time":
		return log.Created
	case "creator":
		return log.CreatedBy
	case "severity":
		return colorType(log.Type)
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

	formatterSlice := []string{
		"%-" + strconv.Itoa(maxColumnWidth["time"]) + "s:",
		"%-" + strconv.Itoa(maxColumnWidth["creator"]) + "s",
		"%-" + strconv.Itoa(maxColumnWidth["severity"]) + "s",
		"%-" + strconv.Itoa(maxColumnWidth["tag"]) + "s",
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
		"severity",
		"tag",
	}

	seenLogIDs := map[int]bool{}
	for {
		logsThisRun := []collins.Log{}
		for _, tag := range tags {
			var logs []collins.Log
			var err error
			if c.IsSet("all") {
				logs, _, err = col.Logs.GetAll(&opts)
			} else {
				logs, _, err = col.Logs.Get(tag, &opts)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to fetch logs for "+tag+": "+err.Error())
			}

			logsThisRun = append(logsThisRun, logs...)
		}

		// Filter logs based on severity passed in
		filteredLogs := []collins.Log{}
		if c.IsSet("severity") {
			for _, log := range logsThisRun {
				if log.Type == strings.ToUpper(c.String("severity")) {
					filteredLogs = append(filteredLogs, log)
				}
			}
		}

		// Get the format for logs this run, sort them, and print
		format := logGetFormat(uniqueSet, filteredLogs)
		sort.Slice(filteredLogs, func(i, j int) bool {
			return filteredLogs[i].ID < filteredLogs[j].ID
		})

		// Only print out messages that we have not seen before and
		// pop the message onto the end of the string formatter
		for _, log := range filteredLogs {
			var fields []interface{}

			for _, column := range uniqueSet {
				fields = append(fields, fieldToLogStruct(column, log))
			}

			fields = append(fields, log.Message)

			if _, ok := seenLogIDs[log.ID]; ok {
				continue
			} else {
				fmt.Printf(format+"\n", fields...)
				seenLogIDs[log.ID] = true
			}
		}

		if !c.IsSet("follow") {
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func logRunCommand(c *cli.Context) error {
	client := getCollinsClient(c)

	if c.IsSet("severity") {
		check := []string{
			"EMERGENCY",
			"ALERT",
			"ERROR",
			"INFORMATIONAL",
			"DEBUG",
			"NOTE",
			"CRITICAL",
			"WARNING",
			"NOTICE"}

		valid := false
		toCheck := strings.ToUpper(c.String("severity"))
		for _, s := range check {
			if s == toCheck {
				valid = true
			}
		}

		if !valid {
			logAndDie("Log severities " + toCheck + " are invalid! Use one of EMERGENCY, ALERT, CRITICAL, ERROR, WARNING, NOTICE, INFORMATIONAL, DEBUG, NOTE")
		}
	}

	if c.IsSet("tags") || c.IsSet("all") {
		tags := strings.Split(c.String("tags"), ",")
		handleLogs(c, client, tags)
	}

	return nil
}
