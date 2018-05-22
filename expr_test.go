package lucener

import (
	"testing"
)

func TestSort(t *testing.T) {

	e := NewExpr()
	e.SortBy("created", false)
	if e.String() != `{"sort":[{"field":"created"}]}` {
		t.Fatalf("invalid sort op: %s", e)
	}

	e.SortBy("created", true)
	if e.String() != `{"sort":[{"field":"created","reverse":true}]}` {
		t.Fatalf("invalid sort op: %s", e)
	}

	// t.Error(e.Filter(Contains("foo", "bar"), Regexp("re", "^foo")).String())
	// t.Error(e.Query(Contains("foo", "bar"), Regexp("re", "^foo")).String())
	// t.Error(e.Refresh(true).String())
}
