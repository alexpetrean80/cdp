package cmd

import (
	"fmt"
	"log"

	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
)

var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Show the last project opened with cdp.",
	Run: func(cmd *cobra.Command, args []string) {
		lastProject, err := lib.ReadLastProject()
		if err != nil {
			fmt.Println("No project found.")
		}
		fmt.Printf("Last project path: %s\n", lastProject)
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
	lastCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		hiddenFlags := []string{"last", "config"}
		for _, f := range hiddenFlags {
			err := c.Flags().MarkHidden(f)
			if err != nil {
				log.Fatal(err)
			}
		}

		c.Parent().HelpFunc()(c, s)
	})
}
