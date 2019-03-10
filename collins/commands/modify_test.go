package commands

import (
	"testing"

	monkey "github.com/bouk/monkey"
	cli "github.com/urfave/cli"
)

// This is largely copied from the main file to create a wrapper in which we
// can run tests in a sane way that allows for testing of short and long form flags
func modifyContext(fn func(*cli.Context), cmd []string) {
	app := cli.App{
		Commands: []cli.Command{
			{
				Name: "modify",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:     "a, set-attribute",
						Usage:    "Set attribute=value. : between key and value. attribute will be uppercased",
						Category: "Modify options",
					},
					cli.StringSliceFlag{
						Name:     "d, delete-attribute",
						Usage:    "Delete attribute",
						Category: "Modify options",
					},
					cli.StringFlag{
						Name:     "S, set-state",
						Usage:    "Set status (and optionally state) to status:state. Requires --reason",
						Category: "Modify options",
					},
					cli.StringFlag{
						Name:     "r, reason",
						Usage:    "Reason for changing status/state",
						Category: "Modify options",
					},
					cli.StringFlag{
						Name:     "l, log",
						Usage:    "Create a log entry",
						Category: "Modify options",
					},
					cli.StringFlag{
						Name:     "L, level",
						Usage:    "Set log level. Default level is NOTE.",
						Value:    "NOTE",
						Category: "Modify options",
					},
					cli.StringFlag{
						Name:     "t, tags",
						Usage:    "Tags to work on, comma separated",
						Category: "Modify options",
					},
				},
				Action: func(ctx *cli.Context) error {
					fn(ctx)
					return nil
				},
			},
		},
	}

	app.Run(cmd)
}

func TestLogCreateOpts(t *testing.T) {

	modifyContext(func(ctx *cli.Context) {
		out := logCreateOpts(ctx)
		if out.Message != "some message" || out.Type != "NOTE" {
			t.Error("Simple log set failed")
		}
	}, []string{"cmd", "modify", "--log", "some message"})

	modifyContext(func(ctx *cli.Context) {
		out := logCreateOpts(ctx)
		if out.Message != "some message" || out.Type != "ERROR" {
			t.Error("Log set level failed")
		}
	}, []string{"cmd", "modify", "--log", "some message", "--level", "ERROR"})

	modifyContext(func(ctx *cli.Context) {
		hitFatalError := false
		monkey.Patch(logAndDie, func(msg string) {
			hitFatalError = true
		})

		logCreateOpts(ctx)
		if !hitFatalError {
			t.Error("level LOL is not valid but no error was thrown")
		}
	}, []string{"cmd", "modify", "--log", "some message", "--level", "LOL"})

}

func TestStatusUpdateOpts(t *testing.T) {

	modifyContext(func(ctx *cli.Context) {
		out := statusUpdateOpts(ctx)
		if out.Status != "allocated" || out.State != "running" {
			t.Error("Simple status set failed")
		}
	}, []string{"cmd", "modify", "--set-state", "allocated:running", "--reason", "lol"})

	modifyContext(func(ctx *cli.Context) {
		out := statusUpdateOpts(ctx)
		if out.Status != "allocated" || out.State != "" {
			t.Error("State was not set but for some reason it has a value")
		}
	}, []string{"cmd", "modify", "--set-state", "allocated", "--reason", "lol"})

	modifyContext(func(ctx *cli.Context) {
		hitFatalError := false
		monkey.Patch(logAndDie, func(msg string) {
			hitFatalError = true
		})

		statusUpdateOpts(ctx)
		if !hitFatalError {
			t.Error("An error should be thrown when reason isn't provided but was not")
		}
	}, []string{"cmd", "modify", "--set-state", "allocated"})

}

func TestAttributeDeleteStrings(t *testing.T) {

	modifyContext(func(ctx *cli.Context) {
		out := attributeDeleteStrings(ctx)
		if out[0] != "hack" {
			t.Error("Really how did you mess this one up?")
		}
	}, []string{"cmd", "modify", "-d", "hack"})

	modifyContext(func(ctx *cli.Context) {
		out := attributeDeleteStrings(ctx)
		if out[0] != "hack" || out[1] != "hacker" {
			t.Error("Really how did you mess this one up?")
		}
	}, []string{"cmd", "modify", "-d", "hack", "-d", "hacker"})

	modifyContext(func(ctx *cli.Context) {
		out := attributeDeleteStrings(ctx)
		if len(out) != 0 {
			t.Error("Really how did you mess this one up?")
		}
	}, []string{"cmd", "modify"})

}

func TestAttributeUpdateOpts(t *testing.T) {

	modifyContext(func(ctx *cli.Context) {
		out := attributeUpdateOpts(ctx)
		if out[0].Attribute != "test;val" {
			t.Error("Simple attribute set failed")
		}
	}, []string{"cmd", "modify", "-a", "test:val"})

	modifyContext(func(ctx *cli.Context) {
		out := attributeUpdateOpts(ctx)
		if out[0].Attribute != "test;val:;!@:" {
			t.Error("Setting attribute with garbage failed")
		}
	}, []string{"cmd", "modify", "-a", "test:val:;!@:"})

	modifyContext(func(ctx *cli.Context) {
		out := attributeUpdateOpts(ctx)
		if out[0].Attribute != "t;v" || out[1].Attribute != "t3;b3" {
			t.Error("Setting multiple attributes failed")
		}
	}, []string{"cmd", "modify", "-a", "t:v", "-a", "t3:b3"})

	modifyContext(func(ctx *cli.Context) {
		hitFatalError := false
		monkey.Patch(logAndDie, func(msg string) {
			hitFatalError = true
		})

		attributeUpdateOpts(ctx)
		if !hitFatalError {
			t.Error("Attribute argument was not set properly but did not throw an error")
		}
	}, []string{"cmd", "modify", "-a", "t"})

}
