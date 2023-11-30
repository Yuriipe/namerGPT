package processing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
)

type Product struct {
	ID          int    `csv:"id"`
	Name        string `csv:"name"`
	Description string `csv:"opispl"`
}

var (
	NumRows          int
	responseFileName = "resp_" + Language + "_" + strconv.Itoa(time.Now().Minute()) + "." + strconv.Itoa(time.Now().Second()) + ".csv"
	Products         = []*Product{}
	Client           = resty.New()
	// "gpt-3.5-turbo-0301" - 3m 25s : 90 v : 30 t : 0.2 temp
	// "gpt-3.5-turbo-16k-0613" - 5m 14s : 90 v : 30 t : 0.2 temp
	// "gpt-3.5-turbo-1106" - 3m 25s : 90 v : 30 t : 0.1 temp
	Model            = "gpt-3.5-turbo-1106" //"gpt-3.5-turbo-1106"
	Language         = "french"
	SourceFile       = `E:\golang_projects\namerGPT\db\tested.csv`
	NameQuery        = "translate to " + Language
	DescriptQuery    = "generate produkt description in" + Language + "for:"
	TokenLimName     = 30
	TokenLimDescript = 200
	Temp             = 0.2
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
	// Use your API KEY here
	apiKey = "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"
)

func QuestionGPT(query string, tokenLimit int) string {

	response, err := Client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": Model,
			"messages": []interface{}{map[string]interface{}{
				"role":    "user",
				"content": query,
			}},
			"max_tokens":  tokenLimit,
			"temperature": Temp,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return ""
	}

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content

}

func changeName(n string, wg *sync.WaitGroup, result chan<- string) {
	defer wg.Done()

	q := NameQuery + n
	t := TokenLimName
	result <- QuestionGPT(q, t)
}

func generateDescription(n string, wg *sync.WaitGroup, result chan<- string) {
	defer wg.Done()

	q := DescriptQuery + n
	t := TokenLimDescript
	result <- QuestionGPT(q, t)
}

func countRows(s string) int {
	file, err := os.OpenFile(s, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// count the number of rows in csv file
	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		panic("Can't count the number of rows")
	}

	NumRows := len(rows)
	fmt.Println(NumRows)
	return NumRows
}

func processGPTQuery(product Product, csvWriter *csv.Writer, wg *sync.WaitGroup) {
	defer wg.Done()

	newName := make(chan string)
	newDesc := make(chan string)

	var innerWg sync.WaitGroup

	innerWg.Add(2)
	go changeName(product.Name, &innerWg, newName)
	go generateDescription(product.Name, &innerWg, newDesc)

	go func() {
		innerWg.Wait()
		close(newName)
		close(newDesc)
	}()

	name := <-newName
	desc := <-newDesc

	res := []string{strconv.Itoa(product.ID), name, desc}
	_ = csvWriter.Write(res)
}

// func speedGPTTest() string {}

func ProcessingAndReturn() {
	// execution timer
	start := time.Now()
	NumRows = countRows(SourceFile) - 1

	file, err := os.OpenFile(SourceFile, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// manipulate file for struct purposes
	if err := gocsv.UnmarshalFile(file, &Products); err != nil {
		panic(err)
	}

	// creating new file to write csv values
	intoFile, err := os.Create(responseFileName)
	if err != nil {
		panic("File was not created. Check settings")
	}
	defer intoFile.Close()

	// creating new csv writer
	csvWriter := csv.NewWriter(intoFile)
	headersCSV := []string{"ID", "name", "description"}
	csvWriter.Write(headersCSV)

	var wg sync.WaitGroup
	//processing query to chat GPT and writing product data to csv

	for _, product := range Products {
		wg.Add(1)
		go processGPTQuery(*product, csvWriter, &wg)
	}

	wg.Wait()

	csvWriter.Flush()

	execution := time.Since(start)
	fmt.Println(execution)
}
