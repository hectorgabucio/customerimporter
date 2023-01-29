// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var errFileNotFound = errors.New("file not found")

var defaultConcurrency = 40
var defaultChunkSize = 64000

var expectedCsvFields = 5
var emailFieldPosition = 2
var emailSeparator = "@"
var csvLineSeparator = "\n"

// A CustomerImporter reads the lines from a csv given as input, and stores it in a data structure.
type CustomerImporter struct {
	logger      *log.Logger
	concurrency int
	chunkSize   int
	store       DomainStore
}

type options struct {
	logger      *log.Logger
	concurrency int
	chunkSize   int
	store       DomainStore
}

// New returns a CustomerImporter with the default options.
// The parsing errors will be silent, the concurrency is set to the defaults
// and the store backend is a treemap.
func New() CustomerImporter {
	return CustomerImporter{logger: nil, concurrency: defaultConcurrency,
		store: NewTreeMap(), chunkSize: defaultChunkSize}
}

// WithOptions returns a new CustomerImporter based on the options input.
// Setting a nil logger would mean silent errors on parsing.
func WithOptions(opt *options) CustomerImporter {
	c := New()
	if opt.concurrency > 0 {
		c.concurrency = opt.concurrency
	}
	if opt.chunkSize > 0 {
		c.chunkSize = opt.chunkSize
	}
	if opt.logger != nil {
		c.logger = opt.logger
	}
	if opt.store != nil {
		c.store = opt.store
	}
	return c
}

// Import reads the input file, extracts the domains and returns an ordered slice of domain entries.
// Can return error if it fails to open the file.
func (c CustomerImporter) Import(filepath string) ([]DomainEntry, error) {
	if ok := fileExists(filepath); !ok {
		return nil, errFileNotFound
	}
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return c.readAndImport(f), nil

}

func (c CustomerImporter) readAndImport(reader io.Reader) []DomainEntry {
	buffer := make([]byte, c.chunkSize)

	var lastLine string

	c.store.Clear()
	var wg sync.WaitGroup

	chunks := make(chan string)
	domains := make(chan string)
	done := make(chan bool)

	for i := 0; i < c.concurrency; i++ {
		wg.Add(1)

		go func() {
			c.extractDomains(chunks, domains)
			wg.Done()
		}()
	}

	go func() {
		for domain := range domains {
			c.save(domain)
		}
		done <- true
	}()

	for {

		n, err := reader.Read(buffer)
		if n == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			c.logError(fmt.Sprintf("err reading csv line %s", err))
			continue
		}

		buffer = append([]byte(lastLine), buffer[:n]...)

		lines := strings.Split(string(buffer), csvLineSeparator)

		if len(lines[len(lines)-1]) > 0 {
			lastLine = lines[len(lines)-1]
			lines = lines[:len(lines)-1]
		} else {
			lastLine = ""
		}

		if len(lines) > 0 {
			chunks <- strings.Join(lines, csvLineSeparator)
		}

	}
	if lastLine != "" {
		chunks <- lastLine
	}

	close(chunks)
	wg.Wait()

	close(domains)
	<-done

	return c.store.GetAll()
}

func (c CustomerImporter) extractDomains(chunks <-chan string, domains chan<- string) {

	for chunk := range chunks {
		r := csv.NewReader(strings.NewReader(chunk))
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.logError(fmt.Sprintf("err reading chunk: %s", err))
				continue
			}
			if len(rec) != expectedCsvFields {
				c.logError(fmt.Sprintf("line does not have correct size %s", err))
				continue
			}

			parts := strings.Split(rec[emailFieldPosition], emailSeparator)

			if len(parts) != 2 {
				c.logError(fmt.Sprintf("failed to extract domain from email %s", rec[emailFieldPosition]))
				continue
			}
			domain := parts[1]
			domains <- domain
		}

	}
}

func (c *CustomerImporter) save(domain string) {
	previousValue, ok := c.store.Get(domain)
	if !ok {
		c.store.Save(domain, 1)
		return
	}
	c.store.Save(domain, previousValue+1)
}

func (c CustomerImporter) logError(msg string) {
	if c.logger == nil {
		return
	}
	c.logger.Println(msg)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
