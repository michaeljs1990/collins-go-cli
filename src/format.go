package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	collins "gopkg.in/tumblr/go-collins.v0/collins"
	yaml "gopkg.in/yaml.v2"
)

func formatAssets(format string, separator string, showHeaders bool, url string, remoteLookup bool, columns []string, assets []collins.Asset) {
	switch format {
	case "table":
		renderTable(separator, showHeaders, columns, assets)
	case "yaml":
		renderYAML(assets)
	case "json":
		renderJSON(assets)
	case "link":
		renderLinks(url, remoteLookup, assets)
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
	case "asset_type":
		return asset.Metadata.Type
	case "classification":
		return asset.Classification.Tag
	case "ipmi_address":
		return asset.IPMI.Address
	case "ip_address":
		return emptyOrValue(len(asset.Addresses), func() string {
			ips := []string{}
			for _, address := range asset.Addresses {
				ips = append(ips, address.Address)
			}
			return strings.Join(ips, ",")
		})
	case "cpu_cores":
		return emptyOrValue(len(asset.CPUs), func() string {
			return strconv.Itoa(asset.CPUs[0].Cores * len(asset.CPUs))
		})
	case "cpu_threads":
		return emptyOrValue(len(asset.CPUs), func() string {
			return strconv.Itoa(asset.CPUs[0].Threads * len(asset.CPUs))
		})
	case "cpu_speed_ghz":
		return emptyOrValue(len(asset.CPUs), func() string {
			return strconv.FormatFloat(float64(asset.CPUs[0].SpeedGhz), 'f', 4, 32)
		})
	case "cpu_description":
		return emptyOrValue(len(asset.CPUs), func() string {
			return asset.CPUs[0].Description
		})
	case "gpu_description":
		return emptyOrValue(len(asset.GPUs), func() string {
			return asset.GPUs[0].Description
		})
	case "cpu_product":
		return emptyOrValue(len(asset.CPUs), func() string {
			return asset.CPUs[0].Product
		})
	case "gpu_product":
		return emptyOrValue(len(asset.GPUs), func() string {
			return asset.GPUs[0].Product
		})
	case "cpu_vendor":
		return emptyOrValue(len(asset.CPUs), func() string {
			return asset.CPUs[0].Vendor
		})
	case "gpu_vendor":
		return emptyOrValue(len(asset.GPUs), func() string {
			return asset.GPUs[0].Vendor
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
			return strconv.Itoa(len(asset.Memory))
		})
	case "disk_storage_human":
		return emptyOrValue(len(asset.Disks), func() string {
			var size float64
			for _, v := range asset.Disks {
				size = size + float64(v.Size)
			}

			return BytesToHumanSize(size)
		})
	case "disk_types":
		return emptyOrValue(len(asset.Disks), func() string {
			disks := UniqueOrderedSet{}
			for _, v := range asset.Disks {
				disks = disks.Add(v.Description)
			}
			return strings.Join(disks, ",")
		})
	default:
		// If it's not special fish it out of atts
		if val, ok := asset.Attributes["0"]; ok {
			return val[strings.ToUpper(field)]
		}

		return ""
	}
}

func renderTable(separator string, showHeaders bool, columns []string, assets []collins.Asset) {
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

	formatter := strings.Join(formatterSlice, separator)

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

func renderYAML(assets []collins.Asset) {
	out, err := yaml.Marshal(&assets)
	if err != nil {
		logAndDie(err.Error())
	}

	fmt.Println(string(out))
}

func renderJSON(assets []collins.Asset) {
	out, err := json.Marshal(&assets)
	if err != nil {
		logAndDie(err.Error())
	}

	fmt.Println(string(out))
}

func renderLinks(url string, remoteLookup bool, assets []collins.Asset) {
	for _, asset := range assets {
		if remoteLookup {
			fmt.Printf("%s/resources?tag=%s&remoteLookup=true\n", url, asset.Metadata.Tag)
		} else {
			fmt.Println(url + "/asset/" + asset.Metadata.Tag)
		}
	}
}
