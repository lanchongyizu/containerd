package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli"
)

var containersCommand = cli.Command{
	Name:    "containers",
	Usage:   "manage containers (metadata)",
	Aliases: []string{"c"},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "print only the container id",
		},
	},
	Subcommands: []cli.Command{
		containersDeleteCommand,
		containersSetLabelsCommand,
		containerInfoCommand,
	},
	ArgsUsage: "[filter, ...]",
	Action: func(context *cli.Context) error {
		var (
			filters     = context.Args()
			quiet       = context.Bool("quiet")
			ctx, cancel = appContext(context)
		)
		defer cancel()

		client, err := newClient(context)
		if err != nil {
			return err
		}
		containers, err := client.Containers(ctx, filters...)
		if err != nil {
			return err
		}
		if quiet {
			for _, c := range containers {
				fmt.Printf("%s\n", c.ID())
			}
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "CONTAINER\tIMAGE\tRUNTIME\tLABELS\t")
		for _, c := range containers {
			var labelStrings []string
			for k, v := range c.Info().Labels {
				labelStrings = append(labelStrings, strings.Join([]string{k, v}, "="))
			}
			labels := strings.Join(labelStrings, ",")
			if labels == "" {
				labels = "-"
			}
			imageName := c.Info().Image
			if imageName == "" {
				imageName = "-"
			}
			record := c.Info()
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%v\t\n",
				c.ID(),
				imageName,
				record.Runtime.Name,
				labels,
			); err != nil {
				return err
			}
		}
		return w.Flush()
	},
}
