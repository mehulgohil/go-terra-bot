package cmd

import (
	"github.com/mehulgohil/go-terra-bot/services"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-terra-bot",
	Short: "Create cloud resource with the help open ai and terraform with just one line",
	Long:  `Create cloud resource with the help open ai and terraform with just one line`,
	Run: func(cmd *cobra.Command, args []string) {
		userPrompt, err := cmd.Flags().GetString("prompt")
		if err != nil {
			log.Fatal(err)
		}

		services.CreateResource(userPrompt)
	},
	Example: `
# Example 1: Add a task
go-terra-bot -p "create an azure resource group"
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootFlags := rootCmd.Flags()
	rootFlags.StringP("prompt", "p", "", "prompt to create a cloud resource")
	cobra.MarkFlagRequired(rootFlags, "prompt")
}
