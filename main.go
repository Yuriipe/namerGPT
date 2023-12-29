package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	"github.com/tkanos/gonfig"
)

type Product struct {
	ID          int    `csv:"id"`
	Name        string `csv:"name"`
	Description string `csv:"description"`
}

type NamerGPTConfig struct {
	NameQuery        string
	DescriptQuery    string
	Model            string  `json:"model"`
	Language         string  `json:"language"`
	SourceFile       string  `json:"sourceFile"`
	TokenLimName     int     `json:"tokenLimName"`
	TokenLimDescript int     `json:"tokenLimDescript"`
	Temp             float64 `json:"temp"`
	APIKey           string  `json:"APIKey"`
	APIEndpoint      string  `json:"APIEndpoint"`
	ResponseFileName string
	Debug            bool
}

type NamerGPT struct {
	cfg   NamerGPTConfig
	resty *resty.Client
}

func (gpt *NamerGPT) questionGPT(query string, tokenLimit int) string {
	response, err := gpt.resty.R().
		SetAuthToken(gpt.cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": gpt.cfg.Model,
			"messages": []interface{}{map[string]interface{}{
				"role":    "user",
				"content": query,
			}},
			"max_tokens":  tokenLimit,
			"temperature": gpt.cfg.Temp,
		}).
		Post(gpt.cfg.APIEndpoint)

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
	return strings.ReplaceAll(content, "\n", "")
}

func (gpt *NamerGPT) changeName(n string, wg *sync.WaitGroup, result chan<- string) {
	gpt.cfg.NameQuery = "translate to " + gpt.cfg.Language
	defer wg.Done()

	q := gpt.cfg.NameQuery + n
	t := gpt.cfg.TokenLimName
	result <- gpt.questionGPT(q, t)
}

func (gpt *NamerGPT) generateDescription(n string, wg *sync.WaitGroup, result chan<- string) {
	gpt.cfg.DescriptQuery = "generate produkt description in" + gpt.cfg.Language + "for:" + gpt.cfg.NameQuery
	defer wg.Done()

	q := gpt.cfg.DescriptQuery + n
	t := gpt.cfg.TokenLimDescript
	result <- gpt.questionGPT(q, t)
}

func (gpt *NamerGPT) Process(product Product, csvWriter *csv.Writer, wg *sync.WaitGroup) {

	defer wg.Done()

	newName := make(chan string)
	newDesc := make(chan string)

	var innerWg sync.WaitGroup

	innerWg.Add(2)
	// fmt.Println(gpt.cfg.NameQuery)
	go gpt.changeName(product.Name, &innerWg, newName)
	go gpt.generateDescription(product.Name, &innerWg, newDesc)

	go func() {
		innerWg.Wait()
		close(newName)
		close(newDesc)
	}()

	name := <-newName
	desc := <-newDesc

	res := []string{strconv.Itoa(product.ID), name, desc}
	_ = csvWriter.Write(res)
	if gpt.cfg.Debug {
		fmt.Println(res)
	}
}

func main() {
	if err := doMain(); err != nil {
		panic(err)
	}
}

func doMain() error {
	fmt.Println("Start processing")

	cfg := NamerGPTConfig{}

	// getting environmental variables
	err := gonfig.GetConf("config/cfg.json", &cfg)
	if err != nil {
		panic(err)
	}

	gpt := &NamerGPT{cfg: cfg, resty: resty.New()}

	start := time.Now()

	file, err := os.OpenFile(cfg.SourceFile, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	products := []*Product{}

	// manipulate file for struct purposes
	if err := gocsv.UnmarshalFile(file, &products); err != nil {
		return err
	}

	// creating new file to write csv values
	cfg.ResponseFileName = "resp_" + cfg.Language + "_" + strconv.Itoa(time.Now().Minute()) + "." + strconv.Itoa(time.Now().Second()) + ".csv"
	fmt.Printf("Creating new csv response file: %s\n", cfg.ResponseFileName)
	intoFile, err := os.Create(cfg.ResponseFileName)
	if err != nil {
		return fmt.Errorf("file was not created. Check settings")
	}
	defer intoFile.Close()

	// creating new csv writer
	csvWriter := csv.NewWriter(intoFile)
	headersCSV := []string{"ID", "name", "description"}
	csvWriter.Write(headersCSV)

	var wg sync.WaitGroup
	//processing query to chat GPT and writing product data to csv
	for _, product := range products {
		wg.Add(1)
		gpt.Process(*product, csvWriter, &wg)
	}

	wg.Wait()

	csvWriter.Flush()

	execution := time.Since(start)
	fmt.Println("Finished in: ", execution)
	return nil

}
