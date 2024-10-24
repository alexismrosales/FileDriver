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
func errorArgumentsException(ctx *cli.Context, args ...string) cli.ExitCoder {
	argsString := strings.Join(args, ", ")
	if ctx.Args().Len() != len(args) || ctx.Args().Len() == 0 {
		return cli.Exit("ERROR:\tArguments "+argsString+" missing.", 0)
	}
	return nil
}

// errorException show all errors not produced by the CLI
func errorException(errs ...error) cli.ExitCoder {
	for _, err := range errs {
		fmt.Println("ERROR: ", err)
	}
	return cli.Exit("", 0)
}

func RunApp() *cli.App {
	fm := &FileManger{
		CurrentDir: "/",
		Flags:      []string{},
	}
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
				Name:      "Set direcction",
				Aliases:   []string{"setdir"},
				Category:  "CONNECTION",
				Usage:     "Connects to the server with an IP address",
				UsageText: "filedriver setdir [address] [port]",
				Action: func(ctx *cli.Context) error {
					if err := errorArgumentsException(ctx, "[address]", "[port]"); err != nil {
						return err
					}
					address := ctx.Args().Get(0)
					port := ctx.Args().Get(1)
					Connect(address, port)
					return nil
				},
			},
			{
				Name:      "disconnect",
				Aliases:   []string{"disconn"},
				Category:  "CONNECTION",
				Usage:     "Disconnect server connection",
				UsageText: "filedriver disconnect",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() != 0 {
						return cli.Exit("ERROR: too many arguments", 0)
					}
					Disconnect()
					return nil
				},
			},
			{
				Name:      "pwd",
				Category:  "FILE NAVIGATION",
				Usage:     "Get directory to the current working directory",
				UsageText: "filedriver pwd",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() != 0 {
						return cli.Exit("ERROR: too many arguments", 0)
					}
					fm.Pwd()
					return nil
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
					if err := errorArgumentsException(ctx, "[Directory]..."); err != nil {
						return err
					}
					fm.Mkdir(paths...)
					return nil
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
						fm.Flags = []string{"a"}
						return nil
					}
					fm.Ls(paths...)
					return nil
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
					fm.Cd(ctx.Args().Get(0))
					return nil
				},
			},
			{
				Name:      "rm",
				Category:  "FILE NAVIGATION",
				Usage:     "Remove files",
				UsageText: "filedriver rm [File|Directory]",
				Action: func(ctx *cli.Context) error {
					paths := ctx.Args().Slice()
					if err := errorArgumentsException(ctx, "[File|Directory]"); err != nil {
						return err
					}
					if ctx.Bool("recursive") {
						fm.Flags = []string{"r"}
					}
					if ctx.Bool("force") {
						fm.Flags = []string{"f"}
					}
					if ctx.Bool("recursiveforced") || (ctx.Bool("force") && ctx.Bool("recursive")) {
						fm.Flags = []string{"r", "f"}
					}
					fm.Rm(paths...)
					return nil
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
					if err := errorArgumentsException(ctx, "[File1]..."); err != nil {
						return err
					}
					fm.Upload(files...)
					return nil
				},
			},
			{
				Name:      "download",
				Category:  "FILE MANAGEMENT",
				Usage:     "Download a file to the server",
				UsageText: "filedriver download [file2] [file2] ...",
				Action: func(ctx *cli.Context) error {
					files := ctx.Args().Slice()
					if err := errorArgumentsException(ctx, "[File1]..."); err != nil {
						return err
					}
					fm.Download(files...)
					return nil
				},
			},
		},
	}

	return app
}
