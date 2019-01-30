package main

import (
	"flag"
	"testing"

	cli "github.com/urfave/cli"
)

func queryGetColumnsContext() (*cli.Context, *flag.FlagSet) {
	set := flag.NewFlagSet("test", 0)
	set.String("columns", "", "doc")
	set.String("extra-columns", "", "doc")
	return cli.NewContext(nil, set, nil), set
}

func TestQueryGetColumns(t *testing.T) {
	// Simple check to see if columns flag override works
	ctx, set := queryGetColumnsContext()
	set.Parse([]string{"--columns", "test,in,order"})

	out := queryGetColumns(ctx)
	if out[0] != "test" || out[1] != "in" || out[2] != "order" || len(out) != 3 {
		t.Error("Parsing columns flag failed test,in,order")
	}

	// Check to make sure setting only one field works
	ctx, set = queryGetColumnsContext()
	set.Parse([]string{"--columns", "first"})
	out = queryGetColumns(ctx)
	if out[0] != "first" || len(out) != 1 {
		t.Error("Parsing columns flag failed first")
	}

	// Ensure that when nothing is set the default values are returned
	defaults := []string{
		"tag",
		"hostname",
		"nodeclass",
		"status",
		"pool",
		"primary_role",
		"secondary_role",
	}

	ctx, set = queryGetColumnsContext()
	set.Parse([]string{""})
	out = queryGetColumns(ctx)
	for i, v := range defaults {
		if out[i] != v {
			t.Error("Failed getting default values when no flag is set. want: ", v, " got: ", out[i])
		}
	}

	if len(out) != 7 {
		t.Error("Column not being set returned more flags than expected")
	}

	// Make sure adding a single extra column works as we expect
	ctx, set = queryGetColumnsContext()
	set.Parse([]string{"--extra-columns", "thing"})
	out = queryGetColumns(ctx)
	if out[7] != "thing" {
		t.Error("Adding one extra column seem to be broken")
	}

	// Make sure adding two extra columns works as we expect
	ctx, set = queryGetColumnsContext()
	set.Parse([]string{"--extra-columns", "thing,two"})
	out = queryGetColumns(ctx)
	if out[7] != "thing" || out[8] != "two" {
		t.Error("Adding two extra column seem to be broken")
	}

	// Make sure adding columns and extra columns returns what we expect
	ctx, set = queryGetColumnsContext()
	set.Parse([]string{"--columns", "doing,this", "--extra-columns", "thing"})
	out = queryGetColumns(ctx)
	if out[0] != "doing" || out[1] != "this" || out[2] != "thing" || len(out) != 3 {
		t.Error("Setting both columns and extra didn't return what we expect")
	}
}
