package cmds

import (
	"strings"
	"testing"

	monkey "github.com/bouk/monkey"
	cli "github.com/urfave/cli"
	collins "gopkg.in/tumblr/go-collins.v0/collins"
)

// This is largely copied from the main file to create a wrapper in which we
// can run tests in a sane way that allows for testing of short and long form flags
func queryContext(fn func(*cli.Context), cmd []string) {
	app := cli.App{
		Commands: []cli.Command{
			{
				Name: "query",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "t, tag",
						Usage:    "Assets with tag[s] value[,...]",
						Category: "Query options",
					},
					cli.BoolFlag{
						Name:     "Z, remote-lookup",
						Usage:    "Query remote datacenters for asset",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "T, type",
						Usage:    "Only show asset with type value",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "n, nodeclass",
						Usage:    "Assets in nodeclass value[,...]",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "p, pool",
						Usage:    "Assets in pool value[,...]",
						Category: "Query options",
					},
					cli.IntFlag{
						Name:     "s, size",
						Usage:    "Number of assets to return per page",
						Value:    100,
						Category: "Query options",
					},
					cli.IntFlag{
						Name:     "limit",
						Usage:    "Limit total results of assets",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "r, role",
						Usage:    "Assets in primary role",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "R, secondary-role",
						Usage:    "Assets in secondary role",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "i, ip-address",
						Usage:    "Assets with IP address[es]",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "S, status",
						Usage:    "Asset status (and optional state after :)",
						Category: "Query options",
					},
					cli.StringSliceFlag{
						Name:     "a, attribute",
						Usage:    "Arbitrary attributes and values to match in query. : between key and value",
						Category: "Query options",
					},
					cli.StringFlag{
						Name:     "o, operation",
						Usage:    "Sets if your query will be joined with AND or OR",
						Value:    "AND",
						Category: "Query options",
					},
					cli.BoolFlag{
						Name:     "H, show-header",
						Usage:    "Show header fields in output",
						Category: "Table formatting",
					},
					cli.StringFlag{
						Name:     "c, columns",
						Usage:    "Attributes to output as columns, comma separated",
						Value:    "tag,hostname,nodeclass,status,pool,primary_role,secondary_role",
						Category: "Table formatting",
					},
					cli.StringFlag{
						Name:     "x, extra-columns",
						Usage:    "Show these columns in addition to the default columns, comma separated",
						Category: "Table formatting",
					},
					cli.StringFlag{
						Name:     "f, field-separator",
						Usage:    "Separator between columns in output",
						Value:    "\t",
						Category: "Table formatting",
					},
					cli.BoolFlag{
						Name:     "l, link",
						Usage:    "Output link to assets found in web UI",
						Category: "Robot formatting",
					},
					cli.BoolFlag{
						Name:     "j, json",
						Usage:    "Output results in JSON",
						Category: "Robot formatting",
					},
					cli.BoolFlag{
						Name:     "y, yaml",
						Usage:    "Output results in YAML",
						Category: "Robot formatting",
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

func TestQueryGetColumns(t *testing.T) {
	// Simple check to see if columns flag override works
	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[0] != "test" || out[1] != "in" || out[2] != "order" || len(out) != 3 {
			t.Error("Parsing columns flag failed test,in,order")
		}
	}, []string{"cmd", "query", "--columns", "test,in,order"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[0] != "first" || len(out) != 1 {
			t.Error("Parsing columns flag failed first")
		}
	}, []string{"cmd", "query", "--columns", "first"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		defaults := []string{
			"tag",
			"hostname",
			"nodeclass",
			"status",
			"pool",
			"primary_role",
			"secondary_role",
		}

		for i, v := range defaults {
			if out[i] != v {
				t.Error("Failed getting default values when no flag is set. want: ", v, " got: ", out[i])
			}
		}

		if len(out) != 7 {
			t.Error("Column not being set returned more flags than expected")
		}

	}, []string{"cmd", "query"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[7] != "thing" {
			t.Error("Adding one extra column seem to be broken")
		}
	}, []string{"cmd", "query", "--extra-columns", "thing"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[7] != "thing" || out[8] != "two" {
			t.Error("Adding two extra column seem to be broken")
		}
	}, []string{"cmd", "query", "--extra-columns", "thing,two"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[0] != "doing" || out[1] != "this" || out[2] != "thing" || len(out) != 3 {
			t.Error("Setting both columns and extra didn't return what we expect")
		}
	}, []string{"cmd", "query", "--extra-columns", "thing", "--columns", "doing,this"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[0] != "thing" || len(out) != 1 {
			t.Error("Columns short flag doesn't seem to be working")
		}

	}, []string{"cmd", "query", "-c", "thing"})

	queryContext(func(ctx *cli.Context) {
		out := queryGetColumns(ctx)
		if out[0] != "thing" || out[1] != "lol" || len(out) != 2 {
			t.Error("Columns short flag doesn't seem to be working")
		}

	}, []string{"cmd", "query", "-c", "thing", "-x", "lol"})

}

func TestGetOutputFormat(t *testing.T) {
	queryContext(func(ctx *cli.Context) {
		out := getOutputFormat(ctx)
		if out != "link" {
			t.Error("Expected to get link formatter got ", out)
		}
	}, []string{"cmd", "query", "--link"})

	queryContext(func(ctx *cli.Context) {
		out := getOutputFormat(ctx)
		if out != "link" {
			t.Error("Expected to get link formatter (short form) got ", out)
		}
	}, []string{"cmd", "query", "-l"})

	queryContext(func(ctx *cli.Context) {
		out := getOutputFormat(ctx)
		if out != "json" {
			t.Error("Expected to get json formatter (short form) got ", out)
		}
	}, []string{"cmd", "query", "-j"})

	queryContext(func(ctx *cli.Context) {
		out := getOutputFormat(ctx)
		if out != "table" {
			t.Error("Expected to get table formatter got ", out)
		}
	}, []string{"cmd", "query"})

	queryContext(func(ctx *cli.Context) {
		out := getOutputFormat(ctx)
		if out != "yaml" {
			t.Error("Expected to get yaml formatter got ", out)
		}
	}, []string{"cmd", "query", "--yaml"})
}

func TestBuildOptionsQuery(t *testing.T) {
	queryContext(func(ctx *cli.Context) {
		out := buildOptionsQuery(ctx, "")
		expected := "(NODECLASS = somenode)"
		if out != expected {
			t.Error("Expected ", expected, " got ", out)
		}
	}, []string{"cmd", "query", "-n", "somenode"})

	queryContext(func(ctx *cli.Context) {
		out := buildOptionsQuery(ctx, "")
		expected := "(NODECLASS = snode) AND (PRIMARY_ROLE = arole) AND (SECONDARY_ROLE = srole)"
		if out != expected {
			t.Error("Expected ", expected, " got ", out)
		}
	}, []string{"cmd", "query", "-n", "snode", "-r", "arole", "-R", "srole"})

	queryContext(func(ctx *cli.Context) {
		out := buildOptionsQuery(ctx, "")
		expected := "(NODECLASS = snode) OR (PRIMARY_ROLE = arole2) OR (SECONDARY_ROLE = srole2)"
		if out != expected {
			t.Error("Expected ", expected, " got ", out)
		}
	}, []string{"cmd", "query", "-n", "snode", "--role", "arole2", "-R", "srole2", "-o", "OR"})

	queryContext(func(ctx *cli.Context) {
		out := buildOptionsQuery(ctx, "")
		expected := "(TAG = U001) OR (NODECLASS = snode) OR (POOL = DEV) OR (PRIMARY_ROLE = arole2) OR (SECONDARY_ROLE = srole2) OR (IP_ADDRESS = 10.0.0.5)"
		if out != expected {
			t.Error("Expected ", expected, " got ", out)
		}
	}, []string{"cmd", "query", "-n", "snode", "--role", "arole2", "-R", "srole2", "-o", "OR", "-t", "U001", "--pool", "DEV", "-i", "10.0.0.5"})

	queryContext(func(ctx *cli.Context) {
		out := buildOptionsQuery(ctx, "")
		expected := "(TAG = U001) AND (NODECLASS = snode) AND (POOL = DEV) AND (PRIMARY_ROLE = arole2) AND (SECONDARY_ROLE = srole2) AND (IP_ADDRESS = 10.0.0.5)"
		if out != expected {
			t.Error("Expected ", expected, " got ", out)
		}
	}, []string{"cmd", "query", "-n", "snode", "--role", "arole2", "-R", "srole2", "-t", "U001", "--pool", "DEV", "-i", "10.0.0.5"})

	queryContext(func(ctx *cli.Context) {
		out := buildOptionsQuery(ctx, "dev")
		expected := "(HOSTNAME = dev) AND (TEST = THING)"
		if out != expected {
			t.Error("Expected ", expected, " got ", out)
		}
	}, []string{"cmd", "query", "-a", "test:thing", "dev"})

	queryContext(func(ctx *cli.Context) {
		hitFatalError := false
		monkey.Patch(logAndDie, func(msg string) {
			hitFatalError = true
		})
		buildOptionsQuery(ctx, "")
		if hitFatalError == false {
			t.Error("LOL should throw a fatal error")
		}
	}, []string{"cmd", "query", "-o", "LOL"})

}

func TestQueryBuildOptions(t *testing.T) {
	queryContext(func(ctx *cli.Context) {
		out := queryBuildOptions(ctx, "")
		expected := "(NODECLASS = somenode)"
		if out.Query != expected || out.Status != "allocated" || out.State != "running" {
			t.Error("Building simple query options failed")
		}
	}, []string{"cmd", "query", "-n", "somenode", "-S", "allocated:running"})

	queryContext(func(ctx *cli.Context) {
		out := queryBuildOptions(ctx, "")
		expected := "(NODECLASS = somenode)"

		if out.Type != "SOME_TYPE" {
			t.Error("1 - Want: ", "SOME_TYPE", " Got: ", out.Type)
		}

		if out.Query != expected {
			t.Error("2 - Want: ", expected, " Got: ", out.Query)
		}

		if out.RemoteLookup != false {
			t.Error("3 - Want: ", false, " Got: ", out.RemoteLookup)
		}

	}, []string{"cmd", "query", "-n", "somenode", "-T", "SOME_TYPE"})

	queryContext(func(ctx *cli.Context) {
		out := queryBuildOptions(ctx, "")
		expected := "((NODECLASS = somenode) OR (NODECLASS = someother))"

		if out.Query != expected {
			t.Error("2 - Want: ", expected, " Got: ", out.Query)
		}
	}, []string{"cmd", "query", "-n", "somenode,someother"})

	queryContext(func(ctx *cli.Context) {
		out := queryBuildOptions(ctx, "")
		expected := "((TAG = M1) OR (TAG = M2) OR (TAG = M3)) AND ((NODECLASS = somenode) OR (NODECLASS = someother)) AND ((PRIMARY_ROLE = test) OR (PRIMARY_ROLE = test2)) AND (IP_ADDRESS = 10.0.0.1)"

		if out.Query != expected {
			t.Error("2 - Want: ", expected, " Got: ", out.Query)
		}
	}, []string{"cmd", "query", "-n", "somenode,someother", "-t", "M1,M2,M3", "--role", "test,test2", "-i", "10.0.0.1"})

	queryContext(func(ctx *cli.Context) {
		out := queryBuildOptions(ctx, "")
		expected := "(TAG = M1) AND (NODECLASS = somenode) AND ((SECONDARY_ROLE = test) OR (SECONDARY_ROLE = test2))"

		if out.Query != expected {
			t.Error("2 - Want: ", expected, " Got: ", out.Query)
		}
	}, []string{"cmd", "query", "-n", "somenode", "-t", "M1", "--secondary-role", "test,test2"})

	queryContext(func(ctx *cli.Context) {
		out := queryBuildOptions(ctx, "hi")
		want := ""
		if out.Attribute != "" {
			t.Error("Want:", want, " Got:", out.Attribute)
		}
	}, []string{"cmd", "query", "-a", "hostname:test", "hi"})
}

// Note that the commands we use don't matter it just shows what command was used when
// generating the json for this.
func TestQueryCLIBasicWorkflow(t *testing.T) {
	client := setup()
	monkey.Patch(getCollinsClient, func(c *cli.Context) *collins.Client {
		return client
	})
	defer teardown()

	SetupGET(201, "/api/assets", "../assets/TestQueryCLIBasicWorkflow.json", t)

	queryContext(func(ctx *cli.Context) {
		c, o, w := captureStdout()
		err := queryRunCommand(ctx)
		result := returnStdout(c, o, w)
		if err != nil {
			t.Error(err.Error())
		}

		tags := []string{"tag30", "tag31"}
		status := []string{"New", "New"}
		for i, line := range strings.Split(result, "\n") {
			parts := strings.Split(line, "\t")
			if parts[0] != tags[i] {
				t.Error("Expected ", tags[i], " got ", parts[0])
			}

			if parts[3] != status[i] {
				t.Error("Expected ", status[i], " got ", parts[3])
			}
		}
	}, []string{"cmd", "query", "-t", "tag30,tag31"})
}

func TestQueryCLIBasicWorkflowSeparator(t *testing.T) {
	client := setup()
	monkey.Patch(getCollinsClient, func(c *cli.Context) *collins.Client {
		return client
	})
	defer teardown()

	SetupGET(201, "/api/assets", "../assets/TestQueryCLIBasicWorkflow.json", t)

	queryContext(func(ctx *cli.Context) {
		c, o, w := captureStdout()
		err := queryRunCommand(ctx)
		result := returnStdout(c, o, w)
		if err != nil {
			t.Error(err.Error())
		}

		tags := []string{"tag30", "tag31"}
		status := []string{"New", "New"}
		for i, line := range strings.Split(result, "\n") {
			parts := strings.Split(line, "lol")
			if parts[0] != tags[i] {
				t.Error("Expected ", tags[i], " got ", parts[0])
			}

			if parts[3] != status[i] {
				t.Error("Expected ", status[i], " got ", parts[3])
			}
		}
	}, []string{"cmd", "query", "-t", "tag30,tag31", "-f", "lol"})
}

func TestQueryCLIGetByAttribute(t *testing.T) {
	client := setup()
	monkey.Patch(getCollinsClient, func(c *cli.Context) *collins.Client {
		return client
	})
	defer teardown()

	SetupGET(201, "/api/assets", "../assets/TestQueryCLIGetByAttribute.json", t)

	queryContext(func(ctx *cli.Context) {
		c, o, w := captureStdout()
		err := queryRunCommand(ctx)
		result := returnStdout(c, o, w)
		if err != nil {
			t.Error(err.Error())
		}

		rows := map[int][]string{
			0: []string{"M0000002", "dev-a7e8c3277b.pit1.terame.com", "devnode", "Provisioned", "DEVELOPMENT", "DEVELOPMENT", "", "bye"},
			1: []string{"M0000001", "plex-8316f3de71.pit1.terame.com", "plexnode", "Allocated", "PRODUCTION", "PLEX", "", "bye"},
		}

		for i, line := range strings.Split(result, "\n") {
			parts := strings.Split(line, "\t")
			for idx, value := range parts {
				if strings.TrimSpace(value) != rows[i][idx] {
					t.Error("Expected ", value, " got ", rows[i][idx])
				}
			}
		}
	}, []string{"cmd", "query", "-a", "hi:bye"})
}

func TestQueryCLIGetGPU(t *testing.T) {
	client := setup()
	monkey.Patch(getCollinsClient, func(c *cli.Context) *collins.Client {
		return client
	})
	defer teardown()

	SetupGET(201, "/api/assets", "../assets/TestQueryCLIGetGPU.json", t)

	queryContext(func(ctx *cli.Context) {
		c, o, w := captureStdout()
		err := queryRunCommand(ctx)
		result := returnStdout(c, o, w)
		if err != nil {
			t.Error(err.Error())
		}

		rows := map[int][]string{
			0: []string{"tumblrtag304", "", "", "New", "", "", "", "NVIDIA Corporation", "GM200GL [Quadro M6000]"},
		}

		for i, line := range strings.Split(result, "\n") {
			parts := strings.Split(line, "\t")
			for idx, value := range parts {
				if strings.TrimSpace(value) != rows[i][idx] {
					t.Error("Expected ", value, " got ", rows[i][idx])
				}
			}
		}
	}, []string{"cmd", "query", "-a", "gpu_vendor:nvidia", "-x", "gpu_product"})
}
