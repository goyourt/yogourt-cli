package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

/* migrate command */
var MigrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the models",
	Long:  "Migrate the models to the configured database",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func executeMigration() {
	InitLogsFile()
	pathToUserProject := "./"

	cmd := exec.Command("go", "run", "./cmd/migrate.go")
	cmd.Dir = pathToUserProject // define working directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func migrate() {
	InitLogsFile()
	executeMigration()
}

/* --- Add migration command to root --- */
func init() {
	rootCmd.AddCommand(MigrationCmd)
}
