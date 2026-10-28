package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/goyourt/yogourt-cli/FileGenerator"
	"github.com/goyourt/yogourt/migrations"
	"github.com/spf13/cobra"
)

var upMapRegex = regexp.MustCompile(`var Up = map\[string\]func\(\*gorm\.DB\)\{[^}]*\}`)
var downMapRegex = regexp.MustCompile(`var Down = map\[string\]func\(\*gorm\.DB\)\{[^}]*\}`)
const migrationFolder = "./migrations/"


/* Commande Create Migration */
var CreateMigrationCmd = &cobra.Command{
	Use:   "make-migration",
	Short: "Generate a new migration file",
	Long:  "Generate a new migration file from the models found in the models folder",
	Run:   createMigration,
}

/* Commande Migrate */
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the migrations",
	Long:  "Run the migration files of the migrations folder",
	Args:  cobra.MaximumNArgs(1),
	Run:   migrate,
}

func init() {
	rootCmd.AddCommand(CreateMigrationCmd)
	MigrateCmd.Flags().Bool("down", false, "Run the Down function of the given migration")
	MigrateCmd.Flags().Bool("up", false, "Run the Up function of the given migration")
	rootCmd.AddCommand(MigrateCmd)
	InitLogsFile()
}

func initMigrations() {
	FileGenerator.CreateFolder(migrationFolder)

	migrateFile := migrationFolder + "migrate.go"
	if _, err := os.Stat(migrateFile); os.IsNotExist(err) {
		FileGenerator.GenerateFile(migrateFile, FileGenerator.GetFileStr("migrate"))
	}

	if err := migrations.InitMigrations(); err != nil {
		fmt.Printf("error while creating migrations table: %v\n", err)
		os.Exit(1)
	}
}

func createMigration(cmd *cobra.Command, args []string) {
	initMigrations()

	timeStamp := time.Now().Format("200601021504")

	migrationFile := migrationFolder + timeStamp + ".go"
	migrationFileContent := FileGenerator.GetComplexFileStr("newMigration", map[string]string{"timeStamp": timeStamp})

	FileGenerator.GenerateFile(migrationFile, migrationFileContent)

	registerMigration(timeStamp)
}

func registerMigration(timeStamp string) {
	migrateFile := migrationFolder + "migrate.go"
	migrateContent, err := os.ReadFile(migrateFile)
	if err != nil {
		fmt.Printf("error while reading migrate.go: %v \n", err)
		return
	}

	content := string(migrateContent)
	content = addToMigrationMap(content, upMapRegex, fmt.Sprintf(`"%s": Up%s, `, timeStamp, timeStamp))
	content = addToMigrationMap(content, downMapRegex, fmt.Sprintf(`"%s": Down%s, `, timeStamp, timeStamp))

	if err := os.WriteFile(migrateFile, []byte(content), os.ModePerm); err != nil {
		fmt.Printf("error while updating migrate.go: %v \n", err)
	}
}

func addToMigrationMap(content string, mapRegex *regexp.Regexp, entry string) string {
	return mapRegex.ReplaceAllStringFunc(content, func(m string) string {
		braceIndex := strings.Index(m, "{")
		header := m[:braceIndex+1]

		// Normalise le corps de la map : retire la virgule finale éventuelle
		// (gofmt peut aussi l'avoir supprimée) puis garantit le séparateur.
		inner := strings.TrimSpace(m[braceIndex+1 : len(m)-1])
		inner = strings.TrimSuffix(inner, ",")
		inner = strings.TrimSpace(inner)
		if inner != "" {
			inner += ", "
		}

		return header + inner + entry + "}"
	})
}

func migrate(cmd *cobra.Command, args []string) {
	initMigrations()

	if _, err := os.Stat(migrationFolder); err != nil {
		fmt.Printf("error while reading migrations folder: %v \n", err)
		return
	}

	down, _ := cmd.Flags().GetBool("down")
	if down && len(args) == 0 {
		fmt.Println("--down requires a migration timestamp")
		return
	}

	up, _ := cmd.Flags().GetBool("up")

	runArgs := append([]string{"run", migrationFolder}, args...)
	if down {
		runArgs = append(runArgs, "--down")
	}
	if up {
		runArgs = append(runArgs, "--up")
	}

	run := exec.Command("go", runArgs...)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	if err := run.Run(); err != nil {
		fmt.Printf("error while running migrations: %v \n", err)
	}
}
