package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/magodo/tfpluginbcd/tfpluginbcd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "tfpluginbcd",
		Version:   getVersion(),
		Usage:     "Terraform plugin schema breaking change detector",
		UsageText: "tfpluginbcd <old schema file> <new schema file>",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 2 {
				return fmt.Errorf("expected two args")
			}
			out, err := tfpluginbcd.Run(ctx.Args().Get(0), ctx.Args().Get(1))
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
