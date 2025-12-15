package cmd

import (
	"github.com/goyourt/yogourt-cli/FileGenerator"
	"github.com/goyourt/yogourt/services"
	"github.com/goyourt/yogourt/services/providers"
	"github.com/spf13/cobra"
)

const signupFolder = "signup/"
const loginFolder = "login/"
const newPasswordFolder = "new-password/"

/* Command Auth */
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Init auth system",
	Long:  "Create basic authentication system",
	Run: func(cmd *cobra.Command, args []string) {
		initAuth()
	},
}

func init() {
	rootCmd.AddCommand(AuthCmd)
}

func initAuth() {
	projectNameInterface = map[string]string{"ProjectName": providers.GetConfig().AppName}

	modelFolder := "./models/"
	roleFile, roleContent := modelFolder+"Role.go", FileGenerator.GetFileStr("role")
	securityFile, securityContent := modelFolder+"Security.go", FileGenerator.GetFileStr("security")
	securityRoleFile, securityRoleContent := modelFolder+"SecurityRole.go", FileGenerator.GetFileStr("securityRole")
	tokenFile, tokenContent := modelFolder+"Token.go", FileGenerator.GetFileStr("token")
	userFile, userContent := modelFolder+"User.go", FileGenerator.GetFileStr("user")
	registryFile, registryContent := modelFolder+"registry.go", FileGenerator.GetFileStr("authRegistry")

	routesFolder := "./api/auth/"
	loginFile, loginContent := routesFolder+loginFolder+"login.go", FileGenerator.GetComplexFileStr("login", projectNameInterface)
	signupFile, signupContent := routesFolder+signupFolder+"signup.go", FileGenerator.GetComplexFileStr("signup", projectNameInterface)
	newPasswordFile, newPasswordContent := routesFolder+newPasswordFolder+"newPassword.go", FileGenerator.GetComplexFileStr("newPassword", projectNameInterface)

	controllerFolder := "./controllers/"
	tokenControllerFile, tokenControllerContent := controllerFolder+"tokenController.go", FileGenerator.GetComplexFileStr("tokenController", projectNameInterface)
	userControllerFile, userControllerContent := controllerFolder+"userController.go", FileGenerator.GetComplexFileStr("userController", projectNameInterface)

	middlewareFile, middlewareFileContent := "./middleware/middleware.go", FileGenerator.GetFileStr("authMiddlewares")

	services.GenerateFile(roleFile, roleContent)
	services.GenerateFile(securityFile, securityContent)
	services.GenerateFile(securityRoleFile, securityRoleContent)
	services.GenerateFile(tokenFile, tokenContent)
	services.GenerateFile(userFile, userContent)
	services.GenerateFile(registryFile, registryContent)

	services.CreateFolder(routesFolder)
	services.CreateFolder(routesFolder + loginFolder)
	services.CreateFolder(routesFolder + signupFolder)
	services.CreateFolder(routesFolder + newPasswordFolder)
	services.GenerateFile(loginFile, loginContent)
	services.GenerateFile(signupFile, signupContent)
	services.GenerateFile(newPasswordFile, newPasswordContent)

	services.GenerateFile(tokenControllerFile, tokenControllerContent)
	services.GenerateFile(userControllerFile, userControllerContent)

	services.GenerateFile(middlewareFile, middlewareFileContent)
}
