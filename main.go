package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Name:  "replace",
		Usage: "find and replace tool for the command line",
		Authors: []*cli.Author{
			{Name: "Brandon Jaus", Email: "brandon.jaus@gmail.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   ".",
				Usage:   "path used to run the find command",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "name used in the find -name flag",
			},
			&cli.StringFlag{
				Name:    "sep",
				Aliases: []string{"s"},
				Value:   "#",
				Usage:   "seperator used between pattern args",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "print the command before running it",
			},
			&cli.BoolFlag{
				Name:  "dryrun",
				Value: false,
				Usage: "print the command instead of running it",
			},
			&cli.BoolFlag{
				Name:    "no-recurse",
				Aliases: []string{"nr"},
				Value:   false,
				Usage:   "convenience flag for setting maxdepth=1",
			},
			&cli.IntFlag{
				Name:    "maxdepth",
				Aliases: []string{"d"},
				Value:   0,
				Usage:   "max depth used in the find -maxdepth flag",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 2 {
				return cli.ShowAppHelp(ctx)
			}

			path := ctx.String("path")
			name := ctx.String("name")
			sep := ctx.String("sep")
			p1, p2 := ctx.Args().Slice()[0], ctx.Args().Slice()[1]

			var b strings.Builder

			fmt.Fprintf(&b, "find %s", path)
			replace := []string{
				"find",
				path,
			}

			if ctx.Bool("no-recurse") {
				if err := ctx.Set("maxdepth", "1"); err != nil {
					return fmt.Errorf("failed to set maxdepth to 1")
				}
			}

			if depth := ctx.Int("maxdepth"); depth > 0 {
				fmt.Fprintf(&b, " -maxdepth %d", depth)
				replace = append(replace,
					[]string{
						"-maxdepth",
						strconv.Itoa(depth),
					}...,
				)
			}

			b.WriteString(" -type f")
			replace = append(replace,
				[]string{
					"-type",
					"f",
				}...,
			)
			if name != "" {
				fmt.Fprintf(&b, " -name %q", name)
				replace = append(replace, "-name", name)
			}

			pattern := fmt.Sprintf("s%s%s%s%s%sg", sep, p1, sep, p2, sep)
			replace = append(replace,
				[]string{
					"-exec",
					"sed",
					"-i",
					"",
					pattern,
					"{}",
					";",
				}...,
			)

			if ctx.Bool("verbose") || ctx.Bool("dryrun") {
				b.WriteString(" -exec sed -i '' ")
				b.WriteString(pattern)
				b.WriteString(" {} \\;")
				fmt.Fprintln(os.Stdout, b.String())
				if ctx.Bool("dryrun") {
					return nil
				}
			}

			find := exec.Command(replace[0], replace[1:]...)
			find.Stdout = os.Stdout
			find.Stderr = os.Stderr

			if err := find.Run(); err != nil {
				return fmt.Errorf("failed to run find command: %v", err)
			}

			return nil
		},
	}
}

func main() {
	app := New()

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
