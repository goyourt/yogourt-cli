package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yogourt",
	Short: "Yogourt CLI",
	Long:  "Command line interface for the Yogourt framework.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bienvenue dans yogourt !")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
