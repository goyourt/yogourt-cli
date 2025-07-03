package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yogourt",
	Short: "yogourt CLI",
	Long:  "Ceci est un CLI pour le package yogourt.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bienvenue dans yogourt !")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
