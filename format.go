package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

func formatAssets(format string, columns []string, assets []collins.Asset) {
	switch format {
	case "table":
		renderTable(columns, assets)
	default:
		log.Fatal(format, " is not a supported format")
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

func renderTable(columns []string, assets []collins.Asset) {
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
		maxColumnWidth[column] = max
	}

	var formatterSlice []string
	for _, maxWidth := range maxColumnWidth {
		formatterSlice = append(formatterSlice, "%-"+strconv.Itoa(maxWidth)+"s")
	}

	formatter := strings.Join(formatterSlice, "\t")

	for _, asset := range assets {
		// We use an interface instead of a slice becasue Printf requires this.
		var fields []interface{}

		for _, column := range columns {
			fields = append(fields, fieldToAssetStruct(column, asset))
		}

		fmt.Printf(formatter+"\n", fields...)
	}
}
