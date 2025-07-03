package cmd

import (
	"fmt"
	"log"
	"os"

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

	configFileContent := `
app_name: "` + ProjectName + `"
version: "1.0.0"
mode: "development"

server:
  port: 8080
  cors: true

database:
  type: "postgres"
  user: "admin"
  password: "password"
  host: "localhost"
  port: 5432
  dbname: "mydb"

paths:
  model_folder: "./models/"
  project_name: ` + ProjectName + `
  main_file: "./main.go"
  api_folder: "./api/"
`

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

	middlewareFileContent := `package middleware

import (
	"github.com/gin-gonic/gin"
)

var Callbacks = map[string]func(*gin.Context){
	"/api":          base,
}

func base(c *gin.Context) {
	c.Next()
}

`

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

	registryFileContent := `package models

var Models = map[string]interface{}{
}
`

	file.WriteString(registryFileContent) //Ecriture du contenu dans le fichier registry.go

	/* --- Création du fichier cmd/migrate.go --- */
	/* Dossier cmd */
	CmdFolder := ProjectName + "/cmd/"

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

	migrationFileContent := `package main

import (
	"log"

	"github.com/goyourt/yogourt-cli/database"
	 ProjectName + "/models"
)

func main() {
	// Initialisation de la base de données
	database.InitDatabase("../config.yaml")

	for name, model := range models.Models {
		if err := database.DB.AutoMigrate(model); err != nil {
			log.Printf("❌ Échec de la migration du modèle '%s': %v", name, err)
		} else {
			log.Printf("✅ Migration réussie pour le modèle '%s'", name)
		}
	}
}
`
	file.WriteString(migrationFileContent) //Ecriture du contenu dans le fichier migrate.go

	/* Fichier main - présent dans le dossier principal */
	MainFile := ProjectName + "/main.go"

	file, mainFileError := os.Create(MainFile)
	if mainFileError != nil {
		fmt.Printf("Erreur lors de la création du fichier main: %v \n", mainFileError)
		log.Printf("ERROR: %s\n", mainFileError) // Ecriture des logs
		return
	}
	defer file.Close() //Fermeture du fichier main

	mainFileContent := `package main

import (
	// Imports
)

func main() {
	/* Bienvenue sur yogourt ! */
}
`
	file.WriteString(mainFileContent) //Ecriture du contenu dans le fichier main.go

	fmt.Println("L'environnement a été initialisé avec succès.")

	/* --- Ecriture des logs --- */
	log.Printf("Environnement initialisé avec succès: %s\n", ProjectName)
}

/* --- Ajout de la commande init à la commande root --- */
func init() {
	rootCmd.AddCommand(InitCmd)
}
