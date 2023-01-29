package customerimporter

import (
	"testing"
)

func Test_treemapDomainStore_ShouldAutoSort(t *testing.T) {
	tr := NewTreeMap()
	tr.Save("c", 1)
	tr.Save("a", 1)
	tr.Save("b", 1)

	data := tr.GetAll()
	if data[0].domain != "a" {
		t.Errorf("expected a in first position, but got %s", data[0].domain)
	}
	if data[1].domain != "b" {
		t.Errorf("expected b in second position, but got %s", data[1].domain)
	}
	if data[2].domain != "c" {
		t.Errorf("expected c in third position, but got %s", data[2].domain)
	}

}

func Test_treemapDomainStore_ShouldBeEmptyAfterClear(t *testing.T) {
	tr := NewTreeMap()
	tr.Save("c", 1)

	data := tr.GetAll()
	if len(data) != 1 {
		t.Errorf("expected 1 element in this store, but got %d", len(data))
	}

	tr.Clear()
	data = tr.GetAll()
	if len(data) != 0 {
		t.Errorf("expected 0 elements in this store, but got %d", len(data))
	}

}
