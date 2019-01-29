package main

import (
	"fmt"
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
	for _, asset := range assets {
		var row string

		for _, column := range columns {
			row = row + " " + fieldToAssetStruct(column, asset)
		}

		fmt.Println(row)
	}
}
