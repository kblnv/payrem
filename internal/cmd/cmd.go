package cmd

import (
	"flag"
	"os"
)

type Command string

const (
	CommandNone Command = ""
	CommandInit Command = "init"
	CommandRun  Command = "run"
)

type InParams struct {
	ConfigPath string
}

type OutParams struct {
	Command    Command
	ConfigPath string
}

type Cmd struct {
	defaultConfigPath string
}

func New(defaultParams InParams) *Cmd {
	return &Cmd{
		defaultConfigPath: defaultParams.ConfigPath,
	}
}

func (c *Cmd) Parse() OutParams {
	args := os.Args[1:]
	if len(args) == 0 {
		printGlobalUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		return c.parseInit()
	case "run":
		return c.parseRun()
	case "-h", "--help":
		printGlobalUsage()
		os.Exit(0)
	default:
		printGlobalUsage()
		os.Exit(1)
	}

	return OutParams{}
}

func printGlobalUsage() {
	println("Usage: payr <command>")
	println()
	println("Commands:")
	println("  init   initialize new config")
	println("  run    run the server")
	println()
	println("Run 'payr <command> --help' for more information on a command.")
}

func (c *Cmd) parseInit() OutParams {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	configPath := fs.String("path", c.defaultConfigPath, "config path")
	fs.Usage = func() {
		println("Usage: payr init")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])

	return OutParams{
		Command:    CommandInit,
		ConfigPath: *configPath,
	}
}

func (c *Cmd) parseRun() OutParams {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	configPath := fs.String("config", c.defaultConfigPath, "config path")
	fs.Usage = func() {
		println("Usage: payr run")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])

	return OutParams{
		Command:    CommandRun,
		ConfigPath: *configPath,
	}
}
