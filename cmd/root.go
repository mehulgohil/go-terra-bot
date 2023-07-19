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
	PreRun: func(cmd *cobra.Command, args []string) {
		_, ok := os.LookupEnv("OPENAPI_KEY")
		if !ok {
			log.Fatal("missing OPENAPI_KEY environment variable")
		}
	},
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
			fmt.Println("Running as `dry-run`. Resources wont be created, only terraform files will be created")
		}

		services.InitializeTerraformFolders()
		services.CreateResource(userPrompt, dryRun)
	},
	Example: `
# Example 1: Create a resource with defaults
go-terra-bot -p "create an azure resource group"

# Example 2: Create a resource with arguments (1)
go-terra-bot -p "create azure resource group named test in westindia"

# Example 3: Create a resource with arguments (2)
go-terra-bot -p "create aws vpc with cidr 10.0.0.0/32"

# Example 4: Create a resource in dry-run
go-terra-bot -p "create an azure resource group --dry-run"
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

	rootFlags.Bool("dry-run", false, "To perform a dry run. TF files will be created. Resources wont be created in cloud infra.")
}
