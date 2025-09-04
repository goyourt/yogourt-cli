package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

/* Commande Migration */
var MigrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Effectue une migration des modeles",
	Long:  "Effectue une migration des modeles vers la base de donnée configurée",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

// Cette fonction sert à exécuter le fichier migrate.go du dossier cmd
func executeMigration() {

	// Initialisation du fichier de logs
	InitLogsFile()

	// Vérification et lecture du fichier config

	// Définir le chemin du projet de l'utilisateur (chemin vers le projet de l'utilisateur)
	pathToUserProject := "./"

	// Crée une commande pour exécuter `go run cmd/migrate/main.go`
	cmd := exec.Command("go", "run", "./cmd/migrate.go")
	cmd.Dir = pathToUserProject // Définit le répertoire de travail du projet
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Exécution de la commande
	cmd.Run()
}

func migrate() {
	InitLogsFile()
	executeMigration()
}

/* --- Ajout de la commande migration à la commande root --- */
func init() {
	rootCmd.AddCommand(MigrationCmd)
}
