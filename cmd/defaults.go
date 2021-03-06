package cmd

import (
	"fmt"

	"github.com/mightymatth/arcli/config"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mightymatth/arcli/utils"
	"github.com/spf13/cobra"
)

func newDefaultsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "defaults",
		Short: "User session defaults",
	}

	c.AddCommand(newDefaultsListCmd())
	c.AddCommand(newDefaultsAddCmd())
	return c
}

func newDefaultsListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "all"},
		Short:   "List of all user session defaults",
		Run: func(cmd *cobra.Command, args []string) {
			drawDefaults(config.Defaults())
		},
	}

	return c
}

func newDefaultsAddCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "add [defaultName] [value]",
		Aliases: []string{"set"},
		Args:    validDefaultsAddArgs(),
		Short:   "Add default value",
		Run: func(cmd *cobra.Command, args []string) {
			err := config.SetDefault(config.DefaultsKey(args[0]), args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("'%v: %v' has been successfully added to defaults.\n", args[0], args[1])
		},
	}

	return c
}

func validDefaultsAddArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(2)(cmd, args)
		if err != nil {
			return err
		}

		if !contains(config.AvailableDefaultsKeys, args[0]) {
			return fmt.Errorf("invalid default (allowed ones: [%v])",
				utils.PrintWithDelimiter(config.AvailableDefaultsKeys))
		}

		if args[0] == string(config.Activity) {
			activities, err := RClient.GetActivities()
			if err != nil {
				return fmt.Errorf("cannot get time entry activities: %v", err)
			}

			_, exists := activities.Valid(args[1])
			if !exists {
				return fmt.Errorf("invalid activity (allowed ones: [%v])",
					utils.PrintWithDelimiter(activities.Names()))
			}
		}

		return nil
	}
}

func drawDefaults(defaults map[string]string) {
	if len(defaults) == 0 {
		fmt.Println("You have no previously defaults set.")
		fmt.Printf("These can be set with: '%v'\n", newDefaultsAddCmd().UseLine())
		return
	}

	t := utils.NewTable()
	t.AppendHeader(table.Row{"Default entity", "Value"})
	for key, val := range defaults {
		t.AppendRow(table.Row{key, val})
	}

	t.Render()
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
