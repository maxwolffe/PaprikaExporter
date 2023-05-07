package main

import (
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Recipe map[string]interface{}

func processGzipFile(file io.Reader, fieldList []string) ([]Recipe, error) {
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	jsonBytes, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	var recipe, extractedRecipe Recipe
	if err := json.Unmarshal(jsonBytes, &recipe); err != nil {
		return nil, err
	}

	extractedRecipe = make(Recipe)
	for _, field := range fieldList {
		if value, ok := recipe[field]; ok {
			extractedRecipe[field] = value
		}
	}

	return []Recipe{extractedRecipe}, nil
}

func main() {
	var inputPath, fields, outputFile string
	flag.StringVar(&inputPath, "input", "", "Input path containing gzipped JSON files or a zip archive")
	flag.StringVar(&fields, "fields", "name,categories", "Comma-separated list of fields to extract from the JSON files")
	flag.StringVar(&outputFile, "outputFile", "", "File to write output to, stdout if not specified")
	flag.Parse()

	if inputPath == "" {
		fmt.Println("Please provide the input path using the -input flag")
		return
	}

	if fields == "" {
		fmt.Println("Please provide the fields to extract using the -fields flag")
		return
	}

	fieldList := strings.Split(fields, ",")
	recipes := []Recipe{}
	var err error

	fi, err := os.Stat(inputPath)
	if err != nil {
		fmt.Println("Error checking input path:", err)
		return
	}

	// If the input path is a directory, walk the directory and process each file
	if fi.IsDir() {
		fmt.Println("Processing directory")
		err = filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == ".paprikarecipe" {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				extractedRecipes, err := processGzipFile(file, fieldList)
				if err != nil {
					return err
				}
				recipes = append(recipes, extractedRecipes...)
			}
			return nil
		})
	} else if filepath.Ext(inputPath) == ".paprikarecipes" {
		fmt.Println("Processing zip file")
		// The paprika output is a zip archive, so we need to extract the files from the archive
		zr, err := zip.OpenReader(inputPath)
		if err != nil {
			fmt.Println("Error opening zip file:", err)
			return
		}
		defer zr.Close()

		for _, zf := range zr.File {
			if filepath.Ext(zf.Name) == ".paprikarecipe" {
				file, err := zf.Open()
				if err != nil {
					fmt.Println("Error opening file in zip:", err)
					return
				}
				defer file.Close()

				extractedRecipes, err := processGzipFile(file, fieldList)
				if err != nil {
					fmt.Println("Error processing gzip file in zip:", err)
					return
				}
				recipes = append(recipes, extractedRecipes...)
			}
		}
	} else {
		fmt.Println("Invalid input path.")
	}

	if err != nil {
		fmt.Println("Error processing files:", err)
		return
	}

	outputJSON, err := json.MarshalIndent(recipes, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	if outputFile == "" {
		_, err = os.Stdout.Write(outputJSON)
		if err != nil {
			fmt.Println("Error writing output to stdout:", err)
			return
		}
	} else {
		err = ioutil.WriteFile(outputFile, outputJSON, 0644)
		if err != nil {
			fmt.Println("Error writing output file:", err)
			return
		}
	}

}
