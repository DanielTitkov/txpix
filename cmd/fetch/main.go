package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/DanielTitkov/txpix/airtable"
	"github.com/DanielTitkov/txpix/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to the configuration file")
	outputDir := flag.String("dir", "", "directory to save the text files")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	res, err := airtable.FetchFromAirtable(
		cfg.Airtable.APIKey,
		cfg.Airtable.BaseID,
		cfg.Airtable.TableName,
		cfg.Airtable.NameID,
		cfg.Airtable.BodyID,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, textfile := range res {
		// Create the file name with the .txt extension
		fileName := fmt.Sprintf("%s.txt", textfile.Name)
		// Join the output directory with the file name
		filePath := filepath.Join(*outputDir, fileName)

		// Write the file
		err := ioutil.WriteFile(filePath, []byte(textfile.Body), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			continue
		}

		fmt.Printf("File %s saved successfully\n", filePath)
	}
}
