/* --- Wizard permettant la création d'un model pour la base de données --- */

package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/goyourt/yogourt-cli/FileGenerator"
	"github.com/goyourt/yogourt/services/providers"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

/* model command */
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Create a model",
	Long:  "Create a table from the model in the database",
	Run: func(cmd *cobra.Command, args []string) {
		CreateModel()
	},
}

func CreateModel() {
	// Init text color
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	type Field struct {
		Name       string
		Type       string
		Constraint string
	}

	type Model struct {
		Name   string
		Fields []Field
	}

	InitLogsFile()
	cfg := providers.GetConfig()

	ModelFolder := cfg.Paths.ModelFolder

	var validName = regexp.MustCompile(`^[a-zA-Z_]+$`)

	/* --- wizard --- */
	// Model name
	fmt.Printf("%s ", blue("What is the model name ?\n"))
	var modelName string
	fmt.Scanln(&modelName)

	ModelNamerunes := []rune(modelName)
	ModelNamerunes[0] = unicode.ToUpper(ModelNamerunes[0])

	if string(ModelNamerunes) == "" || !validName.MatchString(string(ModelNamerunes)) {
		fmt.Println("❌ Invalid name, only letters authorized.")

		return
	}

	fmt.Println(green("Model name: ", string(ModelNamerunes)))

	idType := ""
	survey.AskOne(&survey.Select{
		Message: blue("Choose id type ?"),
		Options: []string{"Int (auto-increment)", "UUID"},
	}, &idType)

	var fields []Field

	for moreFields() {
		fmt.Printf("%s ", blue("Enter the field name ?\n"))
		var fieldName string
		fmt.Scanln(&fieldName)

		// Field name
		FieldNamerunes := []rune(fieldName)
		FieldNamerunes[0] = unicode.ToUpper(FieldNamerunes[0])

		if string(FieldNamerunes) == "" || !validName.MatchString(string(FieldNamerunes)) {
			fmt.Println("❌ Invalid field name, only letters authorized.")
			return
		}

		// Field type
		fieldType := ""
		TypePrompt := &survey.Select{
			Message: blue("Enter field type :"),
			Options: []string{"string", "int", "bool", "float", "datetime"},
		}
		survey.AskOne(TypePrompt, &fieldType)

		// Field constraint
		fieldConstraint := ""
		ConstraintPrompt := &survey.Select{
			Message: blue("Add field constraint :"),
			Options: []string{"NOT NULL", "UNIQUE", "UNIQUE & NOT NULL", "None"},
		}
		survey.AskOne(ConstraintPrompt, &fieldConstraint)

		field := Field{
			Name:       string(FieldNamerunes),
			Type:       fieldType,
			Constraint: fieldConstraint,
		}

		fields = append(fields, field)

		fmt.Println(green("Champ " + string(FieldNamerunes) + " de type " + fieldType + " à été créé avec succès."))
	}

	// Generate model
	newModelFile := ModelFolder + modelName + ".go"

	modelFile, modelFileError := os.Create(newModelFile)
	if modelFileError != nil {
		fmt.Printf("❌ Error while creating model: %s\n", modelFileError)
		log.Printf("ERROR: %s\n", modelFileError)
		return
	}
	defer modelFile.Close()

	modelFileContent := FileGenerator.GetFileStr("newModel")

	modelFileContent += fmt.Sprintf("type %s struct {\n", modelName)
	if idType == "UUID" {
		modelFileContent += "\tID string \t`gorm:\"type:uuid;default:gen_random_uuid();primaryKey\" json:\"id\"`\n"
	} else {
		modelFileContent += "\tID int \t`gorm:\"primaryKey;autoIncrement;not null;unique\" json:\"id\"`\n"
	}

	for _, field := range fields {
		modelFileContent += fmt.Sprintf("\t%s %s ", field.Name, field.Type)
		if field.Constraint != "None" {
			switch field.Constraint {
			case "UNIQUE & NOT NULL":
				modelFileContent += fmt.Sprintf("\t`gorm:\"not null;unique\" json:\"%s\"`\n", field.Name)
			case "NOT NULL":
				modelFileContent += fmt.Sprintf("\t`gorm:\"not null\" json:\"%s\"`\n", field.Name)
			case "UNIQUE":
				modelFileContent += fmt.Sprintf("\t`gorm:\"unique\" json:\"%s\"`\n", field.Name)
			}
		} else {
			modelFileContent += fmt.Sprintf("\t`json:\"%s\"`\n", field.Name)
		}
	}

	modelFileContent += "}\n\n"

	for _, field := range fields {
		modelFileContent += fmt.Sprintf("// Field %s\n", field.Name)
		// Getter
		modelFileContent += fmt.Sprintf("func (u *%s) Get%s() %s {\n", modelName, field.Name, field.Type)
		modelFileContent += fmt.Sprintf("\treturn u.%s\n}\n\n", field.Name)
		// Setter
		modelFileContent += fmt.Sprintf("func (u *%s) Set%s(%s %s) error {", modelName, field.Name, field.Name, field.Type)
		switch field.Type {
		case "int":
			modelFileContent += `
	if ` + field.Name + ` < 0 {
		return fmt.Errorf("value cannot be negative")
	}
	u.` + field.Name + ` = ` + field.Name + `
	return nil
}

`
		case "string":
			modelFileContent += `
	if len(` + field.Name + `) == 0 {
		return fmt.Errorf("field cannot be empty")
	}
	u.` + field.Name + ` = ` + field.Name + `
	return nil
}

`
		case "float":
			modelFileContent += `
	if ` + field.Name + ` < 0.00 {
		return fmt.Errorf("value cannot be negative")
	}
	u.` + field.Name + ` = ` + field.Name + `
	return nil
}

`
		case "bool":
			modelFileContent += `
	u.` + field.Name + ` = ` + field.Name + `
	return nil
}

`
		}

	}

	_, err := modelFile.WriteString(modelFileContent)
	if err != nil {
		fmt.Printf("❌ Error while writing file : %v\n", err)
		log.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Println(green("✅ Model " + modelName + " successfully created."))

	registryFile := ModelFolder + "registry.go"

	registryFileContent, err := os.ReadFile(registryFile)
	if err != nil {
		fmt.Printf("Erreur while reading registry.go: %v \n", err)
		log.Printf("ERROR: %s\n", err)
	}

	contentStr := string(registryFileContent)
	if strings.Contains(contentStr, `"`+modelName+`": &`+modelName+`{}`) {
		fmt.Println("✅ Model already in registry.go")
		return
	}

	mapMarker := "map[string]interface{}{"
	index := strings.Index(contentStr, mapMarker)
	if index == -1 {
		fmt.Println("❌ Cannot find model map in registry.go")
		return
	}

	insertionPoint := index + len(mapMarker)
	newEntry := `"` + modelName + `": &` + modelName + `{},`
	newContent := contentStr[:insertionPoint] + "\n\t" + newEntry + contentStr[insertionPoint:]

	err = os.WriteFile(registryFile, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error writing registry.go: %v\n", err)
		log.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Println("✅ Model successfully added to registry.go")
	log.Printf("Model successfully created: %s\n", modelName)
}

func moreFields() bool {
	var response string
	fmt.Println("Do you want to add another field? (Y/n)")
	fmt.Scanln(&response)

	return strings.ToLower(response) != "n"
}

/* --- Add model command to root --- */
func init() {
	rootCmd.AddCommand(ModelCmd)
}
