package customerimporter

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"example.com/interview/test_util"
)

func TestCustomerImporter_ShouldWorkCorrectly(t *testing.T) {

	c := New()
	data, _ := c.Import("customers.csv")
	first := data[0]
	last := data[len(data)-1]

	expectedFirstDomain := "123-reg.co.uk"
	expectedFirstOccurrences := 8

	expectedLastDomain := "zimbio.com"
	expectedLastOccurrences := 3

	if first.domain != expectedFirstDomain {
		t.Errorf("first expected: %s but got: %s", expectedFirstDomain, first.domain)
		return
	}
	if first.occurrences != expectedFirstOccurrences {
		t.Errorf("first expected occurrences: %d but got: %d", expectedFirstOccurrences, first.occurrences)
		return
	}

	if last.domain != expectedLastDomain {
		t.Errorf("last expected: %s but got: %s", expectedLastDomain, last.domain)
		return
	}
	if last.occurrences != expectedLastOccurrences {
		t.Errorf("last expected occurrences: %d but got: %d", expectedLastOccurrences, last.occurrences)
	}

}

func TestCustomerImporter_ShouldFailWhenFileDoesNotExist(t *testing.T) {
	c := New()
	_, err := c.Import("doesnt-exist.txt")
	if err != errFileNotFound {
		t.Errorf("expected err: %s but got: %s", errFileNotFound, err)
	}
}

func TestCustomerImporter_ShouldReturnEmptyDataForEmptyFile(t *testing.T) {
	c := New()
	buf := test_util.BuildBufferFile(0)
	data := c.readAndImport(buf)

	length := len(data)
	if length != 0 {
		t.Errorf("expected empty data structure but got length: %d", length)
	}
}

func TestCustomerImporter_ShouldSkipIncorrectLines(t *testing.T) {

	input := `a, b, c, d, e
	b
	c`

	c := New()
	data := c.readAndImport(strings.NewReader(input))
	length := len(data)
	if length != 0 {
		t.Errorf("expected empty data structure but got length: %d", length)
	}
}

func TestCustomerImporter_ShouldParseOnlyLinesWithExpectedFieldsNumber(t *testing.T) {

	expected := "hectorgabucio.com"

	input :=
		`a, b, c, d, e
hola@gmail.com, b, c, d, d
a, b, me@hectorgabucio.com, a, b`

	c := New()
	data := c.readAndImport(strings.NewReader(input))
	length := len(data)
	if length != 1 {
		t.Errorf("expected data structure with length 1 but got length: %d", length)
		return
	}
	key := data[0].domain
	if key != expected {
		t.Errorf("expected %s but got %s", expected, key)
		return
	}
}

func TestCustomerImporter_ShouldSkipEmailsWithoutAtSign(t *testing.T) {

	input := `a, b, gmail.com, d, e
	`

	c := New()
	data := c.readAndImport(strings.NewReader(input))
	length := len(data)
	if length != 0 {
		t.Errorf("expected data structure with length 1 but got length: %d", length)
	}

}

type StoreSpy struct {
	data map[string]int
}

func (s *StoreSpy) Clear() {

}
func (s *StoreSpy) Get(domain string) (int, bool) {
	value, ok := s.data[domain]
	return value, ok
}
func (s *StoreSpy) GetAll() []DomainEntry {
	entries := make([]DomainEntry, 0)
	for k, v := range s.data {
		entries = append(entries, DomainEntry{domain: k, occurrences: v})
	}
	return entries
}
func (s *StoreSpy) Save(domain string, value int) {
	s.data[domain] = value
}

func TestCustomerImporter_ShouldUseTheProvidedStore(t *testing.T) {

	input := `a, b, a@gmail.com, d, e
	`
	storeSpy := StoreSpy{data: map[string]int{}}

	c := WithOptions(&options{logger: nil, concurrency: 1, store: &storeSpy})
	c.readAndImport(strings.NewReader(input))
	data := storeSpy.GetAll()
	length := len(data)
	if length != 1 {
		t.Errorf("expected data structure with length 1 but got length: %d", length)
		return
	}

	entry := data[0]
	if entry.domain != "gmail.com" || entry.occurrences != 1 {
		t.Errorf("expected gmail.com entry with 1 occurrence but got: %v", entry)
	}

}

func TestCustomerImporter_ShouldUseTheProvidedLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := WithOptions(&options{logger: logger})
	input :=
		`a,b,c,d,e`
	c.readAndImport(strings.NewReader(input))
	logs := buf.String()
	if !strings.HasPrefix(logs, "failed to extract domain") {
		t.Errorf("expected to log failure to extract domain, but got %s", logs)
	}
}

func TestCustomerImporter_ShouldWorkWithSmallChunkSize(t *testing.T) {
	expected := "hectorgabucio.com"

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := WithOptions(&options{logger: logger, concurrency: 10, chunkSize: 10})
	input :=
		`a, b, c, d, e
hola@gmail.com, b, c, d, d
a, b, me@hectorgabucio.com, a, b
a, hi@hotmail.com, i`
	data := c.readAndImport(strings.NewReader(input))
	length := len(data)
	if length != 1 {
		t.Errorf("expected data structure with length 1 but got length: %d", length)
		return
	}
	key := data[0].domain
	if key != expected {
		t.Errorf("expected %s but got %s", expected, key)
		return
	}

}

func TestCustomerImporter_ShouldWorkWithBigChunkSize(t *testing.T) {
	expected := "hectorgabucio.com"

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := WithOptions(&options{logger: logger, concurrency: 10, chunkSize: 10000})
	input :=
		`a, b, c, d, e
hola@gmail.com, b, c, d, d
a, b, me@hectorgabucio.com, a, b
a, hi@hotmail.com, i`
	data := c.readAndImport(strings.NewReader(input))
	length := len(data)
	if length != 1 {
		t.Errorf("expected data structure with length 1 but got length: %d", length)
		return
	}
	key := data[0].domain
	if key != expected {
		t.Errorf("expected %s but got %s", expected, key)
		return
	}

}
