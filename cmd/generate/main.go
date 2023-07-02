package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/DanielTitkov/txpix/config"
	"github.com/DanielTitkov/txpix/pix"
)

func main() {
	// Parse the command line flags for the directory path and config file path
	dirPath := flag.String("dir", ".", "Directory path to process text files")
	configPath := flag.String("config", "config.yaml", "Configuration file path")
	flag.Parse()

	// Load the configuration from the specified YAML file
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	// Check if the output directory exists, if not create it
	if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
		os.MkdirAll(cfg.OutputDir, 0755)
	}

	// Walk through all the files in the specified directory and its subdirectories
	err = filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the file has a .txt extension, read it and generate the images
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			// Read the file
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			// Use the file name (without the extension) as the output file name
			base := filepath.Base(path)
			output := strings.TrimSuffix(base, filepath.Ext(base)) // remove .txt extension
			output = filepath.Join(cfg.OutputDir, output)          // prepend the output directory path

			// Generate the images
			err = pix.Build(data, output, cfg)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}
