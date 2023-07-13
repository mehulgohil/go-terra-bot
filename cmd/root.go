package cmd

import (
	"fmt"
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

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatal(err)
		}
		if dryRun {
			fmt.Println("Running as `dry-run`. Resources wont be created")
		}

		services.InitializeTerraformFolders()
		services.CreateResource(userPrompt, dryRun)
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
	rootFlags.StringP("prompt", "p", "", "Prompt to create a cloud resource")
	cobra.MarkFlagRequired(rootFlags, "prompt")

	// TODO: change the default value to false, once we have terraform cmds in place
	rootFlags.Bool("dry-run", true, "To perform a dry run. TF files will be created. Resources wont be created in cloud infra.")
}
