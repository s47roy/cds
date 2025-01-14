package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ovh/cds/cli"
	"github.com/ovh/cds/sdk"
)

var adminBroadcastsCmd = cli.Command{
	Name:    "broadcasts",
	Aliases: []string{"broadcast"},
	Short:   "Manage CDS broadcasts",
}

func adminBroadcasts() *cobra.Command {
	return cli.NewCommand(adminBroadcastsCmd, nil, []*cobra.Command{
		cli.NewListCommand(adminBroadcastListCmd, adminBroadcastListRun, nil),
		cli.NewGetCommand(adminBroadcastShowCmd, adminBroadcastShowRun, nil),
		cli.NewCommand(adminBroadcastDeleteCmd, adminBroadcastDeleteRun, nil),
		cli.NewCommand(adminBroadcastCreateCmd, adminBroadcastCreateRun, nil),
	})
}

var adminBroadcastCreateCmd = cli.Command{
	Name:  "create",
	Short: "Create a CDS broadcast",
	Args: []cli.Arg{
		{Name: "title"},
	},
	Flags: []cli.Flag{
		{
			Name:      "level",
			ShortHand: "l",
			Usage:     "Level of broadcast: info or warning",
			Default:   "info",
			IsValid: func(s string) bool {
				if s != "info" && s != "warning" {
					return false
				}
				return true
			},
		},
	},
	Example: `level info:

	cdsctl admin broadcasts create "the title" < content.md

level warning:

	cdsctl admin broadcasts create --level warning "the title" "the content"
	`,
	Aliases: []string{"add"},
}

func adminBroadcastCreateRun(v cli.Values) error {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	bc := &sdk.Broadcast{
		Level:   v.GetString("level"),
		Title:   v.GetString("title"),
		Content: string(content),
	}
	return client.BroadcastCreate(bc)
}

var adminBroadcastShowCmd = cli.Command{
	Name:  "show",
	Short: "Show a CDS broadcast",
	Args: []cli.Arg{
		{Name: "id"},
	},
}

func adminBroadcastShowRun(v cli.Values) (interface{}, error) {
	bc, err := client.BroadcastGet(v.GetString("id"))
	if err != nil {
		return nil, err
	}
	return bc, nil
}

var adminBroadcastDeleteCmd = cli.Command{
	Name:  "delete",
	Short: "Delete a CDS broadcast",
	Args: []cli.Arg{
		{Name: "id"},
	},
	Flags: []cli.Flag{
		{
			Name:  "force",
			Usage: "if true, do not fail if action does not exist",
			IsValid: func(s string) bool {
				if s != "true" && s != "false" {
					return false
				}
				return true
			},
			Default: "false",
			Type:    cli.FlagBool,
		},
	},
}

func adminBroadcastDeleteRun(v cli.Values) error {
	err := client.BroadcastDelete(v.GetString("id"))
	if v.GetBool("force") && sdk.ErrorIs(err, sdk.ErrNoBroadcast) {
		fmt.Println(err)
		return nil
	}
	return err
}

var adminBroadcastListCmd = cli.Command{
	Name:  "list",
	Short: "List CDS broadcasts",
}

func adminBroadcastListRun(v cli.Values) (cli.ListResult, error) {
	srvs, err := client.Broadcasts()
	if err != nil {
		return nil, err
	}
	return cli.AsListResult(srvs), nil
}
