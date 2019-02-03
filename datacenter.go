package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func datacenterSubcommand() cli.Command {
	return cli.Command{
		Name:    "datacenter",
		Aliases: []string{"dc"},
		Usage:   "Manage multiple Collins configurations",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "n, new",
				Usage:    "Create a new configuration file for Value at ~/.collins.yml.Vaulue",
				Category: "Datacenter options",
			},
			cli.StringFlag{
				Name:     "H, host",
				Usage:    "Use value for host when setting up new datacenter",
				Category: "Datacenter options",
			},
			cli.StringFlag{
				Name:     "u, username",
				Usage:    "Use value for username when setting up new datacenter",
				Category: "Datacenter options",
			},
			cli.StringFlag{
				Name:     "p, password",
				Usage:    "Use value for password when setting up new datacenter",
				Category: "Datacenter options",
			},
			cli.BoolFlag{
				Name:     "l, list",
				Usage:    "List configured collins instances",
				Category: "Datacenter options",
			},
		},
		Action: datacenterRunCommand,
	}
}

func makeNewDatacenterConfig(c *cli.Context, conf string) {
	confDefaults := struct {
		Timeout  int    `yaml:"timeout"`
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password,omitempty"`
	}{
		Timeout:  120,
		Host:     c.String("host"),
		Username: c.String("username"),
		Password: c.String("password"),
	}

	reader := bufio.NewReader(os.Stdin)
	if !c.IsSet("host") {
		fmt.Print("Enter Collins URI for " + c.String("new") + " (i.e. https://collins." + c.String("new") + ".company.net): ")
		host, _ := reader.ReadString('\n')
		confDefaults.Host = strings.TrimSpace(host)
	}

	if !c.IsSet("username") {
		userDefault := os.Getenv("USER")
		fmt.Print("Enter username (default: " + userDefault + "): ")
		username, _ := reader.ReadString('\n')
		if username == "" {
			confDefaults.Username = userDefault
		} else {
			confDefaults.Username = strings.TrimSpace(username)
		}
	}

	if !c.IsSet("password") {
		fmt.Print("Enter password: ")
		pass, _ := reader.ReadString('\n')
		confDefaults.Password = strings.TrimSpace(pass)
	}

	ybytes, _ := yaml.Marshal(&confDefaults)
	outYamlFile := path.Join(os.Getenv("HOME"), ".collins.yml."+c.String("new"))
	err := ioutil.WriteFile(outYamlFile, ybytes, 0600)
	if err != nil {
		logAndDie(err.Error())
	}
}

func datacenterRunCommand(c *cli.Context) error {
	// Check if the main config file is a symlink which means that we possibly control it
	confFile := path.Join(os.Getenv("HOME"), ".collins.yml")
	fd, err := os.Lstat(confFile)
	if err != nil {
		logAndDie(err.Error())
	}

	if fd.Mode()&os.ModeSymlink != os.ModeSymlink {
		logAndDie("Unable to determine default Collins datacenter: " + confFile + " is not a symlink, which means \"collins dc\" is not managing this configuration")
	}

	if c.IsSet("new") && !c.IsSet("list") {
		makeNewDatacenterConfig(c, confFile)
	}

	if c.IsSet("list") {
		files, err := filepath.Glob(os.Getenv("HOME") + "/.collins.yml.*")
		if err != nil {
			logAndDie(err.Error())
		}

		currentConf, err := os.Readlink(confFile)
		if err != nil {
			logAndDie(err.Error())
		}

		for _, file := range files {
			filename := path.Base(file)
			prettyOutput := strings.SplitAfterN(filename, ".", 4)[3]
			if file == currentConf {
				fmt.Println(prettyOutput + " *")
			} else {
				fmt.Println(prettyOutput)
			}
		}
	}

	if !c.IsSet("list") && !c.IsSet("new") {
		switchToConf := ""
		if c.NArg() > 0 {
			switchToConf = c.Args().Get(0)

			files, err := filepath.Glob(os.Getenv("HOME") + "/.collins.yml.*")
			if err != nil {
				logAndDie(err.Error())
			}

			validDc := false
			for _, file := range files {
				filename := path.Base(file)
				prettyOutput := strings.SplitAfterN(filename, ".", 4)[3]
				if prettyOutput == switchToConf {
					validDc = true
				}
			}

			if !validDc {
				logAndDie("No Collins configuration for datacenter \"" + switchToConf + "\" found. Perhaps you want to create it with 'collins dc --new " + switchToConf + "'?")
			}

			if err = os.Remove(confFile); err != nil {
				logAndDie(err.Error())
			}
			if err = os.Symlink(confFile+"."+switchToConf, confFile); err != nil {
				logAndDie(err.Error())
			}

		}
	}

	return nil
}
