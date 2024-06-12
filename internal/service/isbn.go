package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/util"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

type ISBNService struct {
	apiBasePath string
	client      HTTPClient
	csvWriter   util.Writer
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
	bookChan := make(chan models.Book)
	isbnChan := make(chan util.ISBN)
	updateErrorChan := make(chan error)
	csvWriteSuccessChan := make(chan bool)

	// Fetch books from the index endpoint
	go s.fetchBooks(bookChan)

	// Convert ISBN-10 <=> ISBN-13
	go s.convertISBN(bookChan, isbnChan)

	// Update missing ISBNs via the update endpoint
	go s.updateISBN(isbnChan, updateErrorChan)

	// Append new ISBNs to a CSV file
	go s.appendToCSV(isbnChan, csvWriteSuccessChan)

	successfulWrites := 0
	for isSuccess := range csvWriteSuccessChan {
		if isSuccess {
			successfulWrites++
		}
	}

	s.csvWriter.Flush()
	if err := s.csvWriter.Error(); err != nil {
		fmt.Println("Error closing CSV writer:", err)
	}
}

// fetchBooks Fetch books from the index endpoint
func (s *ISBNService) fetchBooks(outChan chan<- models.Book) {
	defer close(outChan)

	nextPage := int32(1)

	for nextPage != 0 {
		res, err := s.client.Get(fmt.Sprintf("%s/books?page=%d", s.apiBasePath, nextPage))
		if err != nil {
			log.Println("Error making HTTP request:", err)
			return
		}
		defer res.Body.Close()

		var data models.PaginatedBooks
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			log.Println("Error decoding JSON:", err)
			return
		}

		for _, book := range data.Items {
			outChan <- book
		}

		nextPage = data.NextPage
	}
}

// convertISBN Convert ISBN-10 <=> ISBN-13
func (s *ISBNService) convertISBN(inChan <-chan models.Book, outChan chan<- util.ISBN) {
	defer close(outChan)

	for book := range inChan {
		var isbn util.ISBN
		if len(book.ISBN13) != 13 {
			isbn = *util.NewISBN(book.ISBN10)
		} else if len(book.ISBN10) != 10 {
			isbn = *util.NewISBN(book.ISBN13)
		} else {
			continue
		}

		outChan <- isbn
	}
}

// updateISBN Update missing ISBNs via the update endpoint
func (s *ISBNService) updateISBN(inChan <-chan util.ISBN, outChan chan<- error) {
	defer close(outChan)

	for isbn := range inChan {
		url := fmt.Sprintf("%s/books/%s", s.apiBasePath, isbn.ISBN13)
		data, err := json.Marshal(isbn)
		if err != nil {
			log.Printf("error encoding data: %v\n", err)
			outChan <- err
			continue
		}
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
		if err != nil {
			log.Printf("error making HTTP PUT request: %v\n", err)
			outChan <- err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := s.client.Do(req)
		if err != nil {
			log.Printf("error making HTTP PUT request: %v\n", err)
			outChan <- err
			continue
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Printf("received non-OK status code: %d\n", res.StatusCode)
			outChan <- fmt.Errorf("error: received non-OK status code: %d", res.StatusCode)
			continue
		}
	}
}

// appendToCSV Append new ISBNs to a CSV file
func (s *ISBNService) appendToCSV(inChan <-chan util.ISBN, outChan chan<- bool) {
	defer close(outChan)
	for isbn := range inChan {
		var record []string
		if isbn.SourceType == util.ISBN13 {
			record = []string{isbn.ISBN13}
		} else if isbn.SourceType == util.ISBN10 {
			record = []string{isbn.ISBN10}
		}
		s.csvWriter.Write(record)

		outChan <- true
	}
}
