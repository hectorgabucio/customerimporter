package customerimporter

import "github.com/igrmk/treemap/v2"

// A DomainEntry represents the number of occurrences for a specific domain
type DomainEntry struct {
	domain      string
	occurrences int
}

// A DomainStore takes care of persisting the domain entries.
type DomainStore interface {
	// Get returns the number of occurrences stored for a given domain, and a bool to indicate if exists.
	Get(domain string) (int, bool)
	// Save sets value to the specified domain.
	Save(domain string, value int)
	// Clear will clear the state of the data structure.
	Clear()
	// GetAll returns an array of domain entries, sorted.
	GetAll() []DomainEntry
}

type treemapDomainStore struct {
	data *treemap.TreeMap[string, int]
}

// NewTreeMap returns a domain store with a treemap backend.
func NewTreeMap() DomainStore {
	return &treemapDomainStore{data: treemap.New[string, int]()}
}

func (t *treemapDomainStore) Get(domain string) (int, bool) {
	return t.data.Get(domain)
}

func (t *treemapDomainStore) Save(domain string, value int) {
	t.data.Set(domain, value)
}

func (t *treemapDomainStore) Clear() {
	t.data.Clear()
}

func (t *treemapDomainStore) GetAll() []DomainEntry {
	iterator := t.data.Iterator()
	entries := make([]DomainEntry, 0)
	for iterator.Valid() {
		entries = append(entries, DomainEntry{domain: iterator.Key(), occurrences: iterator.Value()})
		iterator.Next()
	}

	return entries
}
