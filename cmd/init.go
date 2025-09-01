package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/goyourt/yogourt-cli/FileGenerator"
	"github.com/spf13/cobra"
)

/* Commande init */
var InitCmd = &cobra.Command{
	Use:   "init [projectName]",
	Short: "Initialise un nouveau projet yogourt",
	Long:  "Crée la structure de base pour un nouveau projet yogourt",
	Args:  cobra.ExactArgs(1), //Attends un seul argument
	Run: func(cmd *cobra.Command, args []string) {
		ProjectName := args[0]

		CreateConfigFile(ProjectName)
		InitProject(ProjectName)
		createMiddlewareFile(ProjectName)
	},
}

/* Fonction pour l'initialisation du fichier de logs */
func InitLogsFile() {

	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

/* --- Création du fichier config --- */
func CreateConfigFile(ProjectName string) {

	// Initialisation du fichier de logs
	InitLogsFile()

	configPath := "./config.yaml"

	// Création du fichier config
	file, configFileError := os.Create(configPath)
	if configFileError != nil {
		log.Printf("ERROR: %s\n", configFileError) // Ecriture des logs
		log.Fatalf("Erreur lors de la création du fichier de config: %v\n", configFileError)
		return
	}
	defer file.Close()

	configFileContent := FileGenerator.GetComplexFileStr("config", map[string]string{"ProjectName": ProjectName})

	file.WriteString(configFileContent) //Ecriture du contenu dans le fichier config
}

/* --- Fin création du fichier config --- */

/* --- Création du fichier middleware --- */
func createMiddlewareFile(ProjectName string) {

	// Initialisation du fichier de logs
	InitLogsFile()

	/* Dossier middleware - présent dans le dossier principal */
	MiddlewareFolder := "./middleware/"

	middlewareFolderError := os.Mkdir(MiddlewareFolder, os.ModePerm)

	if middlewareFolderError != nil {
		fmt.Printf("Erreur lors de la création du dossier middleware: %v \n", middlewareFolderError)
		log.Printf("ERROR: %s\n", middlewareFolderError) // Ecriture des logs
		return
	}

	//Création du fichier middleware
	MiddlewareFile := MiddlewareFolder + "middleware.go"

	file, middlewareFileError := os.Create(MiddlewareFile)
	if middlewareFileError != nil {
		fmt.Printf("Erreur lors de la création du fichier middleware: %v \n", middlewareFileError)
		log.Printf("ERROR: %s\n", middlewareFileError) // Ecriture des logs
		return
	}
	defer file.Close() //Fermeture du fichier config

	middlewareFileContent := FileGenerator.GetFileStr("middlewares")

	file.WriteString(middlewareFileContent) //Ecriture du contenu dans le fichier middleware
}

/* --- Fin création du fichier middleware --- */

/* --- Fonction pour la commande "init" du package --- */
func InitProject(ProjectName string) {

	// Initialisation du fichier de logs
	InitLogsFile()

	/* Dossier route - présent dans le dossier principal */
	ApiFolder := "./api/"

	apiFolderError := os.Mkdir(ApiFolder, os.ModePerm)

	if apiFolderError != nil {
		fmt.Printf("Erreur lors de la création du dossier api: %v \n", apiFolderError)
		log.Printf("ERROR: %s\n", apiFolderError) // Ecriture des logs
		return
	}

	/* Dossier model - présent dans le dossier principal */
	ModelFolder := "./models/"

	modelFolderError := os.Mkdir(ModelFolder, os.ModePerm)

	if modelFolderError != nil {
		fmt.Printf("Erreur lors de la création du dossier models: %v \n", modelFolderError)
		log.Printf("ERROR: %s\n", modelFolderError) // Ecriture des logs
		return
	}

	/* Fichier modelRegistry - présent dans le dossier models */
	ModelRegistryFile := ModelFolder + "registry.go"

	file, modelRegistryError := os.Create(ModelRegistryFile)

	if modelRegistryError != nil {
		fmt.Printf("Erreur lors de la création du fichier models/registry.go: %v \n", modelRegistryError)
		log.Printf("ERROR: %s\n", modelRegistryError) // Ecriture des logs
		return
	}
	defer file.Close() //Fermeture du fichier registry.go

	registryFileContent := FileGenerator.GetFileStr("models")

	file.WriteString(registryFileContent) //Ecriture du contenu dans le fichier registry.go

	/* --- Création du fichier cmd/migrate.go --- */
	/* Dossier cmd */
	CmdFolder := "./cmd/"

	cmdFolderError := os.Mkdir(CmdFolder, os.ModePerm)

	if cmdFolderError != nil {
		fmt.Printf("Erreur lors de la création du dossier cmd: %v \n", cmdFolderError)
		log.Printf("ERROR: %s\n", cmdFolderError) // Ecriture des logs
		return
	}

	/* Fichier migrate.go */
	MigrateFile := CmdFolder + "/migrate.go"

	file, migrateFileError := os.Create(MigrateFile)
	if migrateFileError != nil {
		fmt.Printf("Erreur lors de la création du fichier migrate: %v \n", migrateFileError)
		log.Printf("ERROR: %s\n", migrateFileError) // Ecriture des logs
		return
	}
	defer file.Close() //Fermeture du fichier migrate

	migrationFileContent := FileGenerator.GetComplexFileStr("migration", map[string]string{"ProjectName": ProjectName})
	file.WriteString(migrationFileContent) //Ecriture du contenu dans le fichier migrate.go

	/* Fichier main - présent dans le dossier principal */
	MainFile := "./main.go"

	file, mainFileError := os.Create(MainFile)
	if mainFileError != nil {
		fmt.Printf("Erreur lors de la création du fichier main: %v \n", mainFileError)
		log.Printf("ERROR: %s\n", mainFileError) // Ecriture des logs
		return
	}
	defer file.Close() //Fermeture du fichier main

	mainFileContent := FileGenerator.GetFileStr("main")
	file.WriteString(mainFileContent) //Ecriture du contenu dans le fichier main.go

	fmt.Println("L'environnement a été initialisé avec succès.")

	/* --- Ecriture des logs --- */
	log.Printf("Environnement initialisé avec succès: %s\n", ProjectName)
}

/* --- Ajout de la commande init à la commande root --- */
func init() {
	rootCmd.AddCommand(InitCmd)
}
