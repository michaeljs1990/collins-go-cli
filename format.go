package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func formatAssets(format string, showHeaders bool, columns []string, assets []collins.Asset) {
	switch format {
	case "table":
		renderTable(columns, showHeaders, assets)
	default:
		logAndDie(format + " is not a supported format")
	}
}

// Can't thing of a very nice and idiomatic way to handle converting the fields
// a user passes into a printable struct without mapping it like this.
func fieldToAssetStruct(field string, asset collins.Asset) string {
	switch field {
	case "tag":
		return asset.Metadata.Tag
	case "status":
		return asset.Metadata.Status
	default:
		// If it's not special fish it out of atts
		return asset.Attributes["0"][strings.ToUpper(field)]
	}
}

func renderTable(columns []string, showHeaders bool, assets []collins.Asset) {
	// Find the longest column in each field so the final output is pretty.
	maxColumnWidth := make(map[string]int)
	for _, column := range columns {
		var max int
		for _, asset := range assets {
			length := len(fieldToAssetStruct(column, asset))
			if length > max {
				max = length
			}
		}

		// If headers are going to be output make sure we take them into
		// account when formatting the table.
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

	for _, asset := range assets {
		// We use an interface instead of a slice becasue Printf requires this.
		var fields []interface{}

		for _, column := range columns {
			fields = append(fields, fieldToAssetStruct(column, asset))
		}

		fmt.Printf(formatter+"\n", fields...)
	}
}
