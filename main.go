package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/magodo/tfpluginbcd/tfpluginbcd"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func main() {
	var (
		flagAll         bool
		flagRules       string
		flagCustomRules cli.StringSlice
	)

	app := &cli.App{
		Name:    "tfpluginbcd",
		Version: getVersion(),
		Usage:   "Terraform plugin schema breaking change detector",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list pre-defined rules",
				Action: func(ctx *cli.Context) error {
					var names []string
					for name := range tfpluginbcd.Rules {
						names = append(names, name)
					}
					sort.StringSlice(names).Sort()

					for _, name := range names {
						rule := tfpluginbcd.Rules[name]
						fmt.Printf("%s: %s\n", rule.ID, rule.Description)
					}
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run the breaking change detector and show breaking changes (all changes will be shown if no option is specified).",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "all",
						EnvVars:     []string{"TFPLUGINBCD_ALL"},
						Usage:       "Enable all pre-defined rules",
						Destination: &flagAll,
					},
					&cli.StringFlag{
						Name:        "rules",
						EnvVars:     []string{"TFPLUGINBCD_RULES"},
						Usage:       "One or more pre-defined rule names (separated by comma)",
						Destination: &flagRules,
					},
					&cli.StringSliceFlag{
						Name:        "custom-rule",
						EnvVars:     []string{"TFPLUGINBCD_CUSTOM_RULE"},
						Usage:       "Custom breaking change rule expression",
						Destination: &flagCustomRules,
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() != 2 {
						return fmt.Errorf("expected two args")
					}

					var opt tfpluginbcd.Opt
					if flagAll {
						var allRules []string
						for name := range tfpluginbcd.Rules {
							allRules = append(allRules, name)
						}
						slices.Sort(allRules)
						opt.Rules = allRules
					} else {
						if flagRules != "" {
							var rules []string
							for _, rule := range strings.Split(flagRules, ",") {
								rules = append(rules, strings.TrimSpace(rule))
							}
							opt.Rules = rules
						}
					}
					opt.CustomRuleExprs = flagCustomRules.Value()

					out, err := tfpluginbcd.Run(ctx.Context, ctx.Args().Get(0), ctx.Args().Get(1), opt)
					if err != nil {
						return err
					}
					fmt.Println(out)
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
