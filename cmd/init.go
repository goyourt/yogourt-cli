package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/goyourt/yogourt-cli/FileGenerator"
	"github.com/goyourt/yogourt/services"
	"github.com/spf13/cobra"
)

var projectNameInterface map[string]string

/* init command */
var initCmd = &cobra.Command{
	Use:   "init [projectName]",
	Short: "Initialize a new yogourt project",
	Long:  "Create base structure for a new yogourt project with the given project name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		projectNameInterface = map[string]string{"ProjectName": projectName}
		createConfigFile()
		initProject(projectName)
		createMiddlewareFile()
	},
}

func InitLogsFile() {

	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func createConfigFile() {
	InitLogsFile()

	configPath := "./config.yaml"
	configFileContent := FileGenerator.GetComplexFileStr("config", projectNameInterface)

	services.GenerateFile(configPath, configFileContent)
}

func createMiddlewareFile() {
	InitLogsFile()

	middlewareFolder := "./middleware/"
	services.CreateFolder(middlewareFolder)

	middlewareFile := middlewareFolder + "middleware.go"
	middlewareFileContent := FileGenerator.GetFileStr("middlewares")

	services.GenerateFile(middlewareFile, middlewareFileContent)
}

func initProject(projectName string) {
	InitLogsFile()

	services.CreateFolder("./api/")
	services.CreateFolder("./public/")
	services.CreateFolder("./public/files/")

	modelFolder := "./models/"
	services.CreateFolder(modelFolder)

	modelRegistryFile := modelFolder + "registry.go"
	registryFileContent := FileGenerator.GetFileStr("registry")

	services.GenerateFile(modelRegistryFile, registryFileContent)

	cmdFolder := "./cmd/"
	services.CreateFolder(cmdFolder)

	migrateFile := cmdFolder + "/migrate.go"
	migrationFileContent := FileGenerator.GetComplexFileStr("migration", projectNameInterface)

	services.GenerateFile(migrateFile, migrationFileContent)

	dockerComposeFileContent := FileGenerator.GetComplexFileStr("docker-compose", projectNameInterface)
	services.GenerateFile("./docker-compose.yml", dockerComposeFileContent)

	mainFileContent := FileGenerator.GetFileStr("main")
	services.GenerateFile("./main.go", mainFileContent)

	controllerFolder := "./controllers/"
	services.CreateFolder(controllerFolder)

	fileRouteFolder := "./api/file/"
	fileRouteContent := FileGenerator.GetComplexFileStr("file", projectNameInterface)
	fileControllerContent := FileGenerator.GetComplexFileStr("fileController", projectNameInterface)
	fileModelContent := FileGenerator.GetFileStr("fileModel")
	services.CreateFolder(fileRouteFolder)
	services.GenerateFile(fileRouteFolder+"file.go", fileRouteContent)
	services.GenerateFile(controllerFolder+"fileController.go", fileControllerContent)
	services.GenerateFile(modelFolder+"File.go", fileModelContent)

	fmt.Println("Environment successfully created.")

	log.Printf("Environment successfully created: %s\n", projectName)
}

/* --- Add init command to root --- */
func init() {
	rootCmd.AddCommand(initCmd)
}
