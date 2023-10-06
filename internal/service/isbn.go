package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/emiliogozo/xyz-books/internal/api"
	"github.com/emiliogozo/xyz-books/internal/util"
)

type ISBNService struct {
	apiBasePath string
	client      *http.Client
	csvWriter   *util.CsvWriter
}

// NewService creates a new instance of the Service
func NewISBNService(serverAddress string, outputPath string) *ISBNService {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	outputCSV := fmt.Sprintf("%s/isbn.csv", outputPath)

	cw, _ := util.NewCsvWriter(outputCSV)

	return &ISBNService{
		apiBasePath: fmt.Sprintf("http://%s%s", serverAddress, config.APIBasePath),
		client:      &http.Client{},
		csvWriter:   cw,
	}
}

func (s *ISBNService) Run() {
	indexChan := make(chan api.Book)
	convertChan := make(chan util.ISBN)
	updateChan := make(chan util.ISBN)
	csvChan := make(chan bool)

	var wg sync.WaitGroup

	// Fetch books from the index endpoint
	wg.Add(1)
	go s.fetchBooks(indexChan, &wg)

	// Convert ISBN-10 <=> ISBN-13
	numConverters := 4
	for i := 0; i < numConverters; i++ {
		wg.Add(1)
		go s.convertISBN(indexChan, convertChan, &wg)
	}

	// Update missing ISBNs via the update endpoint
	numUpdaters := 2
	for i := 0; i < numUpdaters; i++ {
		wg.Add(1)
		go s.updateISBN(convertChan, updateChan, &wg)
	}

	// Append new ISBNs to a CSV file
	wg.Add(1)
	go s.appendToCSV(updateChan, csvChan, &wg)

	for range csvChan {
		s.csvWriter.Flush()
	}
	wg.Wait()
}

// fetchBooks Fetch books from the index endpoint
func (s *ISBNService) fetchBooks(ch chan<- api.Book, wg *sync.WaitGroup) {
	defer wg.Done()

	nextPage := int32(1)

	for nextPage != -1 {
		res, err := s.client.Get(fmt.Sprintf("%s/books?page=%d", s.apiBasePath, nextPage))
		if err != nil {
			log.Println("Error making HTTP request:", err)
			return
		}
		defer res.Body.Close()

		var data api.PaginatedBooks
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			log.Println("Error decoding JSON:", err)
			return
		}

		for _, book := range data.Items {
			ch <- book
		}

		nextPage = data.NextPage
		if data.NextPage == 0 {
			nextPage = -1
		}
	}
	close(ch)
}

// convertISBN Convert ISBN-10 <=> ISBN-13
func (s *ISBNService) convertISBN(inputChan <-chan api.Book, outputChan chan<- util.ISBN, wg *sync.WaitGroup) {
	defer wg.Done()

	for book := range inputChan {
		var isbn util.ISBN
		if len(book.ISBN13) != 13 {
			isbn = *util.NewISBN(book.ISBN10)
		} else if len(book.ISBN10) != 10 {
			isbn = *util.NewISBN(book.ISBN13)
		} else {
			continue
		}

		outputChan <- isbn
	}
}

// updateISBN Update missing ISBNs via the update endpoint
func (s *ISBNService) updateISBN(inputChan <-chan util.ISBN, outputChan chan<- util.ISBN, wg *sync.WaitGroup) {
	defer wg.Done()

	for isbn := range inputChan {
		url := fmt.Sprintf("%s/books/%s", s.apiBasePath, isbn.ISBN13)
		data, err := json.Marshal(isbn)
		if err != nil {
			log.Printf("error encoding data: %v\n", err)
			return
		}
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
		if err != nil {
			log.Printf("error making HTTP PUT request: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := s.client.Do(req)
		if err != nil {
			log.Printf("error making HTTP PUT request: %v\n", err)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Printf("received non-OK status code: %d\n", res.StatusCode)
			return
		}

		outputChan <- isbn
	}
}

// appendToCSV Append new ISBNs to a CSV file
func (s *ISBNService) appendToCSV(inputChan <-chan util.ISBN, outputChan chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for isbn := range inputChan {
		var record []string
		if isbn.SourceType == util.ISBN13 {
			record = []string{isbn.ISBN13}
		} else if isbn.SourceType == util.ISBN10 {
			record = []string{isbn.ISBN10}
		}
		s.csvWriter.Write(record)

		outputChan <- true
	}
}
