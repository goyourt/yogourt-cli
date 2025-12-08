package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yogourt",
	Short: "yogourt CLI",
	Long:  "This is the Command Line Interface for yogourt.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to yogourt !")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
