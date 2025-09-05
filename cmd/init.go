package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/goyourt/yogourt-cli/FileGenerator"
	"github.com/spf13/cobra"
)

var projectNameInterface map[string]string

/* Commande init */
var initCmd = &cobra.Command{
	Use:   "init [projectName]",
	Short: "Initialise un nouveau projet yogourt",
	Long:  "Crée la structure de base pour un nouveau projet yogourt",
	Args:  cobra.ExactArgs(1), //Attends un seul argument
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		projectNameInterface = map[string]string{"ProjectName": projectName}
		createConfigFile()
		initProject(projectName)
		createMiddlewareFile()
		initAuth()
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
func createConfigFile() {

	// Initialisation du fichier de logs
	InitLogsFile()

	configPath := "./config.yaml"
	configFileContent := FileGenerator.GetComplexFileStr("config", projectNameInterface)

	FileGenerator.GenerateFile(configPath, configFileContent)
}

/* --- Fin création du fichier config --- */

/* --- Création du fichier middleware --- */
func createMiddlewareFile() {
	InitLogsFile()

	middlewareFolder := "./middleware/"
	FileGenerator.CreateFolder(middlewareFolder)

	//Création du fichier middleware
	middlewareFile := middlewareFolder + "middleware.go"
	middlewareFileContent := FileGenerator.GetFileStr("middlewares")

	FileGenerator.GenerateFile(middlewareFile, middlewareFileContent)
}

/* --- Fin création du fichier middleware --- */

/* --- Fonction pour la commande "init" du package --- */
func initProject(projectName string) {
	InitLogsFile()

	FileGenerator.CreateFolder("./api/")

	modelFolder := "./models/"
	FileGenerator.CreateFolder(modelFolder)

	/* Fichier modelRegistry - présent dans le dossier models */
	modelRegistryFile := modelFolder + "registry.go"
	registryFileContent := FileGenerator.GetFileStr("registry")

	FileGenerator.GenerateFile(modelRegistryFile, registryFileContent)

	/* --- Création du fichier cmd/migrate.go --- */
	/* Dossier cmd */
	cmdFolder := "./cmd/"
	FileGenerator.CreateFolder(cmdFolder)

	/* Fichier migrate.go */
	migrateFile := cmdFolder + "/migrate.go"
	migrationFileContent := FileGenerator.GetComplexFileStr("migration", projectNameInterface)

	FileGenerator.GenerateFile(migrateFile, migrationFileContent)

	/* docker-compose */
	dockerComposeFileContent := FileGenerator.GetComplexFileStr("docker-compose", projectNameInterface)
	FileGenerator.GenerateFile("./docker-compose.yml", dockerComposeFileContent)

	/* Fichier main - présent dans le dossier principal */
	mainFileContent := FileGenerator.GetFileStr("main")
	FileGenerator.GenerateFile("./main.go", mainFileContent)

	fmt.Println("L'environnement a été initialisé avec succès.")

	/* --- Ecriture des logs --- */
	log.Printf("Environnement initialisé avec succès: %s\n", projectName)
}

func initAuth() {
	modelFolder := "./models/"
	roleFile, roleContent := modelFolder+"Role.go", FileGenerator.GetFileStr("role")
	securityFile, securityContent := modelFolder+"Security.go", FileGenerator.GetFileStr("security")
	securityRoleFile, securityRoleContent := modelFolder+"SecurityRole.go", FileGenerator.GetFileStr("securityRole")
	tokenFile, tokenContent := modelFolder+"Token.go", FileGenerator.GetFileStr("token")
	userFile, userContent := modelFolder+"User.go", FileGenerator.GetFileStr("user")

	routesFolder := "./api/auth/"
	loginFile, loginContent := routesFolder+"login/login.go", FileGenerator.GetComplexFileStr("login", projectNameInterface)
	signupFile, signupContent := routesFolder+"signup/signup.go", FileGenerator.GetComplexFileStr("signup", projectNameInterface)

	controllerFolder := "./controllers/"
	tokenControllerFile, tokenControllerContent := controllerFolder+"tokenController.go", FileGenerator.GetComplexFileStr("tokenController", projectNameInterface)
	userControllerFile, userControllerContent := controllerFolder+"userController.go", FileGenerator.GetComplexFileStr("userController", projectNameInterface)

	serviceFolder := "./services/"
	userServicesFile, userServicesContent := serviceFolder+"userServices.go", FileGenerator.GetComplexFileStr("userServices", projectNameInterface)

	FileGenerator.GenerateFile(roleFile, roleContent)
	FileGenerator.GenerateFile(securityFile, securityContent)
	FileGenerator.GenerateFile(securityRoleFile, securityRoleContent)
	FileGenerator.GenerateFile(tokenFile, tokenContent)
	FileGenerator.GenerateFile(userFile, userContent)

	FileGenerator.CreateFolder(routesFolder)
	FileGenerator.CreateFolder(routesFolder + "login")
	FileGenerator.CreateFolder(routesFolder + "signup")
	FileGenerator.GenerateFile(loginFile, loginContent)
	FileGenerator.GenerateFile(signupFile, signupContent)
	FileGenerator.CreateFolder(controllerFolder)
	FileGenerator.GenerateFile(tokenControllerFile, tokenControllerContent)
	FileGenerator.GenerateFile(userControllerFile, userControllerContent)
	FileGenerator.CreateFolder(serviceFolder)
	FileGenerator.GenerateFile(userServicesFile, userServicesContent)
}

/* --- Ajout de la commande init à la commande root --- */
func init() {
	rootCmd.AddCommand(initCmd)
}
