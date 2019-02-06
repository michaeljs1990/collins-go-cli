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

func emptyOrValue(sliceSize int, fn func() string) string {
	if sliceSize > 0 {
		return fn()
	} else {
		return ""
	}
}

func fieldToAssetStruct(field string, asset collins.Asset) string {
	switch field {
	case "tag":
		return asset.Metadata.Tag
	case "status":
		return asset.Metadata.Status
	case "state":
		return asset.Metadata.State.Name
	case "classification":
		return asset.Classification.Tag
	case "cpu_cores":
		return emptyOrValue(len(asset.CPUs), func() string {
			return strconv.Itoa(asset.CPUs[0].Cores)
		})
	case "cpu_threads":
		return emptyOrValue(len(asset.CPUs), func() string {
			return strconv.Itoa(asset.CPUs[0].Threads)
		})
	case "cpu_speed_ghz":
		return emptyOrValue(len(asset.CPUs), func() string {
			return strconv.FormatFloat(float64(asset.CPUs[0].SpeedGhz), 'f', 4, 32)
		})
	case "cpu_description":
		return emptyOrValue(len(asset.CPUs), func() string {
			return asset.CPUs[0].Description
		})
	case "cpu_product":
		return emptyOrValue(len(asset.CPUs), func() string {
			return asset.CPUs[0].Product
		})
	case "cpu_vendor":
		return emptyOrValue(len(asset.CPUs), func() string {
			return asset.CPUs[0].Vendor
		})
	case "memory_size_bytes":
		return emptyOrValue(len(asset.Memory), func() string {
			bytes := 0
			for _, v := range asset.Memory {
				bytes = bytes + v.Size
			}
			return strconv.Itoa(bytes)
		})
	case "memory_size_total":
		return emptyOrValue(len(asset.Memory), func() string {
			var size float64
			format := ""
			for _, v := range asset.Memory {
				split := strings.Split(v.SizeHuman, " ")
				format = split[1]
				pop, _ := strconv.ParseFloat(split[0], 64)
				size = size + pop
			}
			return strconv.FormatFloat(size, 'f', 2, 64) + " " + format
		})
	case "memory_description":
		return emptyOrValue(len(asset.Memory), func() string {
			return asset.Memory[0].Description
		})
	case "memory_banks_total":
		return emptyOrValue(len(asset.Memory), func() string {
			return strconv.Itoa((asset.Memory[len(asset.Memory)-1].Bank + 1))
		})
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
