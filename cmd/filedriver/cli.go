package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	FileDriverIntro = `FileDriver is a tool that works for upload files on a server`
)

// errorArgumentsException manage all args into an error
func errorArgumentsException(ctx *cli.Context, nArguments bool, args ...string) cli.ExitCoder {
	if nArguments && ctx.Args().Len() != 0 {
		return nil
	}

	if !nArguments && ctx.Args().Len() == len(args) {
		return nil
	}
	argsString := strings.Join(args, ", ")
	if len(args) == 0 {
		return cli.Exit("ERROR:\tNo arguments required", 0)
	}
	return cli.Exit("ERROR:\tArguments "+argsString+" missing.", 0)
}

// errorException show all errors not produced by the CLI
func errorException(errs ...error) cli.ExitCoder {
	for _, err := range errs {
		fmt.Println("ERROR: ", err)
	}
	return cli.Exit("", 0)
}

// RunApp show the cli output on terminal, after a command is written,
// all the exceptions for arguments on the CLI are managed and call
// the command_parser functions
func RunApp() *cli.App {
	cp := NewCommandParser("/", "")
	app := &cli.App{
		Name:      "filedriver",
		Compiled:  time.Now(),
		Usage:     FileDriverIntro,
		UsageText: "filedriver <command> [arguments]",
		Action: func(*cli.Context) error {
			fmt.Println(FileDriverIntro)
			fmt.Println("\nGLOBAL COMMANDS")
			fmt.Println(cli.HelpFlag)
			fmt.Println(cli.VersionFlag)
			fmt.Println("\nAUTHOR")
			fmt.Println("Name: \t\tAlexis M.Rosales")
			fmt.Println("Github: \thttps://github.com/alexismrosales")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "setAddress",
				Aliases:   []string{"setaddr"},
				Category:  "CONNECTION",
				Usage:     "Connects to the server with an IP address",
				UsageText: "filedriver setaddr [address] [port]",
				Action: func(ctx *cli.Context) error {
					if err := errorArgumentsException(ctx, false, "[address]", "[port]"); err != nil {
						return err
					}
					address := ctx.Args().Get(0)
					port := ctx.Args().Get(1)
					return cp.FirstConnection(address, port)
				},
			},
			{
				Name:      "disconnect",
				Aliases:   []string{"disconn"},
				Category:  "CONNECTION",
				Usage:     "Disconnect server connection",
				UsageText: "filedriver disconnect",
				Action: func(ctx *cli.Context) error {
					if err := errorArgumentsException(ctx, false); err != nil {
						return err
					}
					return cp.Disconnect()
				},
			},
			{
				Name:      "pwd",
				Category:  "FILE NAVIGATION",
				Usage:     "Get directory to the current working directory",
				UsageText: "filedriver pwd",
				Action: func(ctx *cli.Context) error {
					if err := errorArgumentsException(ctx, false); err != nil {
						return err
					}
					return cp.Pwd()
				},
			},
			{
				Name:      "mkdir",
				Aliases:   []string{"md"},
				Category:  "FILE NAVIGATION",
				Usage:     "Create a new directory",
				UsageText: "filedriver mkdir [Directory1] [Directory2] ...",
				Action: func(ctx *cli.Context) error {
					paths := ctx.Args().Slice()
					if err := errorArgumentsException(ctx, true, "[Directory]..."); err != nil {
						return err
					}
					return cp.Mkdir(paths...)
				},
			},
			{
				Name:      "ls",
				Category:  "FILE NAVIGATION",
				Usage:     "Show list of the content of a directory",
				UsageText: "filedriver ls [Optional]",
				Action: func(ctx *cli.Context) error {
					paths := ctx.Args().Slice()
					if ctx.Bool("all") {
						cp.flags = []string{"a"}
						return nil
					}
					return cp.Ls(paths...)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:               "all",
						Aliases:            []string{"a"},
						Usage:              "Show all files",
						DisableDefaultText: true,
					},
				},
			},
			{
				Name:      "cd",
				Category:  "FILE NAVIGATION",
				Usage:     "Navigate between directories",
				UsageText: "filedriver cd [Path]",
				Action: func(ctx *cli.Context) error {
					if err := errorArgumentsException(ctx, false, "[Path]"); err != nil {
						return err
					}
					return cp.Cd(ctx.Args().Get(0))
				},
			},
			{
				Name:      "mv",
				Category:  "FILE EDITION",
				Usage:     "Move/Rename files",
				UsageText: "filedriver mv [Directory1] [Directory2] ... [DestinationDirectory]",
				Action: func(ctx *cli.Context) error {
					paths := ctx.Args().Slice()

					return cp.Mv(paths...)
				},
			},
			{
				Name:      "rm",
				Category:  "FILE NAVIGATION",
				Usage:     "Remove files",
				UsageText: "filedriver rm [File|Directory]",
				Action: func(ctx *cli.Context) error {
					paths := ctx.Args().Slice()
					if err := errorArgumentsException(ctx, true, "[File|Directory]"); err != nil {
						return err
					}
					if ctx.Bool("recursive") {
						cp.flags = []string{"r"}
					}
					if ctx.Bool("force") {
						cp.flags = []string{"f"}
					}
					if ctx.Bool("recursiveforced") || (ctx.Bool("force") && ctx.Bool("recursive")) {
						cp.flags = []string{"r", "f"}
					}
					return cp.Rm(paths...)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:               "recursive",
						Aliases:            []string{"r"},
						Usage:              "Recursive delete",
						DisableDefaultText: true,
					},
					&cli.BoolFlag{
						Name:               "force",
						Aliases:            []string{"f"},
						Usage:              "Forced deletion",
						DisableDefaultText: true,
					},
					&cli.BoolFlag{
						Name:               "recursiveforced",
						Aliases:            []string{"rf"},
						DisableDefaultText: true,
					},
				},
			},
			{
				Name:      "upload",
				Category:  "FILE MANAGEMENT",
				Usage:     "Upload a file to the server",
				UsageText: "filedriver upload [file1] [file2] ...",
				Action: func(ctx *cli.Context) error {
					files := ctx.Args().Slice()
					if err := errorArgumentsException(ctx, true, "[File1]..."); err != nil {
						return err
					}
					return cp.Upload(files...)
				},
			},
			{
				Name:      "download",
				Category:  "FILE MANAGEMENT",
				Usage:     "Download a file to the server",
				UsageText: "filedriver download [file2] [file2] ...",
				Action: func(ctx *cli.Context) error {
					files := ctx.Args().Slice()
					if err := errorArgumentsException(ctx, true, "[File1]..."); err != nil {
						return err
					}
					return cp.Download(files...)
				},
			},
		},
	}
	return app
}
