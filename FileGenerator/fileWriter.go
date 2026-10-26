package FileGenerator

import (
	"fmt"
	"log"
	"os"
)

func GenerateFile(fileName string, fileContent string) {
	file, fileError := os.Create(fileName)
	if fileError != nil {
		fmt.Printf("error while creating file: %v \n", fileError)
		log.Printf("ERROR: %s\n", fileError)
		return
	}
	defer file.Close()

	file.WriteString(fileContent)
}

func CreateFolder(folderPath string) {
	folderError := os.Mkdir(folderPath, os.ModePerm)
	if folderError != nil {
		fmt.Printf("error while creating folder: %v \n", folderError)
		log.Printf("ERROR: %s\n", folderError)
		return
	}
}
