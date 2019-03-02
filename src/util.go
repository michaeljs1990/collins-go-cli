package cmds

import (
	"fmt"
	"math"
	"os"
	"strconv"

	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

// This is kinda dumb but go has pretty limited data structure types
// this means no Sets. We only support adding values via the Add method
// creating a UniqueOrderedSet with duplicate values will not cause it to
// be filtered. This is not efficient.
type UniqueOrderedSet []string

func (u UniqueOrderedSet) Contains(s string) bool {

	for _, val := range u {
		if val == s {
			return true
		}
	}

	return false
}

func (u UniqueOrderedSet) Add(s string) UniqueOrderedSet {
	if u.Contains(s) {
		return u
	}

	return append(u, s)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// BytesToHumanSize takes an int and treats it as if it was bytes
// converting it to the largest human readable size.
func BytesToHumanSize(b int) string {
	size := float64(b)
	suffix := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	base := math.Log(size) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffix[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)

}

func getCollinsClient(c *cli.Context) *collins.Client {
	// We just set the COLLINS_CLIENT_CONFIG so we can use the NewClientFromYaml
	// helper function still which will still run through all the default config paths.
	if c.GlobalIsSet("config") {
		os.Setenv("COLLINS_CLIENT_CONFIG", c.GlobalString("config"))
	}

	collins, err := collins.NewClientFromYaml()
	if err != nil {
		fmt.Println("You can use COLLINS_CLIENT_CONFIG env or --config to set the location of your config")
		logAndDie(err.Error())
	}

	return collins
}
