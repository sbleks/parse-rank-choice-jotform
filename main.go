package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func usage(program string) {
	log.Printf("[Usage]: %s <input file path> <output file path>\n", program)
}

type rankChoiceItem struct {
	rank int
	data string
}

type rankChoiceData []rankChoiceItem

func TransformDataTo2DSlice(d rankChoiceData) [][]string {

	numRows := len(d)
	result := make([][]string, numRows+1)

	// Add header row
	result[0] = []string{"Rank", "Item"}

	// Add data rows
	for i := 0; i < numRows; i++ {
		result[i+1] = []string{
			fmt.Sprintf("%d", d[i].rank),
			d[i].data,
		}
	}

	return result
}

func main() {

	program := os.Args[0]

	// Checks for a file path to be provided
	if len(os.Args) < 2 {
		log.Println("An input file path was not provided")
		usage(program)
		os.Exit(1)
	}

	filePath := os.Args[1]

	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.File
	file, err := os.Open(filePath)

	// Checks for the error
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	// Closes the file
	defer file.Close()

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		fmt.Println("Error reading records", err)
	}

	var rankingData rankChoiceData

	// Loop to iterate through
	// and parse the ranks
	for i, eachrecord := range records {
		if i > 0 {
			item := strings.Split(strings.ReplaceAll(eachrecord[20], "\r\n", "\n"), "\n")

			for _, it := range item {
				splitItem := strings.Split(it, ": ")
				rank, err := strconv.Atoi(splitItem[0])
				if err != nil {
					log.Fatal("Could not parse rank", err)
				}

				rankingData = append(rankingData, rankChoiceItem{
					rank: rank,
					data: splitItem[1],
				})
			}
		}
	}

	// Transform the data from an array of rankChoiceItem structs
	// to a 2DSlice that can be used in csv.WriteAll
	data := TransformDataTo2DSlice(rankingData)

	// Check if an output path was provided or use a default
	var outputPath string
	if len(os.Args) != 3 {
		outputPath = os.Args[2]
	} else {
		outputPath = "./output/rank-choice.csv"
	}

	// Create file to be written to
	csvFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	// Create writer and write data to file
	w := csv.NewWriter(csvFile)
	w.WriteAll(data)

	// Check if data when writing csv
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	os.Exit(0)
}
