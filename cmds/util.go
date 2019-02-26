package cmds

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
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

var debugMode bool

func setDebugMode(c *cli.Context) {
	debugMode = c.GlobalIsSet("debug")
}

func debugLog(msg string) {
	if debugMode {
		fmt.Fprintln(os.Stderr, "DEBUG: "+msg)
	}
}

// Helper functions for printing
func logAndDie(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func printSuccess() {
	fmt.Println(color.Green("SUCCESS"))
}

func printSuccessWithMsg(msg string) {
	fmt.Println(color.Green("SUCCESS "), "("+msg+")")
}

func printError(e string) {
	fmt.Println(color.Red("ERROR "), "("+e+")")
}
