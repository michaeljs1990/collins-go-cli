package cmds

import (
	"fmt"
	"os"

	color "github.com/logrusorgru/aurora"
)

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
