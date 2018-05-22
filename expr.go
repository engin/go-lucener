package lucener

import "github.com/gocql/gocql"
import "encoding/json"

type Rule interface {
}

// type Query struct {
// }
// type Filter struct {
// 	field string
// 	Type  string
// 	value interface{}
// }

type Sort struct {
	Field   string `json:"field,omitempty"`
	Reverse bool   `json:"reverse,omitempty"`
}

type Expr struct {
	Q []Rule  `json:"query,omitempty"`
	F []Rule  `json:"filter,omitempty"`
	S []*Sort `json:"sort,omitempty"`
	R bool    `json:"refresh,omitempty"`
}

type value struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type values struct {
	Type   string      `json:"type"`
	Field  string      `json:"field"`
	Values interface{} `json:"values"`
}

// Match query
func Match(f string, v interface{}) Rule {
	return &value{
		Type:  "match",
		Field: f,
		Value: v,
	}
}

// Prefix query
func Prefix(f string, v interface{}) Rule {
	return &value{
		Type:  "prefix",
		Field: f,
		Value: v,
	}
}

type allQuery struct {
	Type string `json:"type"`
}

// All selects all indexed rows
func All() Rule {
	return &allQuery{
		Type: "all",
	}
}

type fuzzyQuery struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// Fuzzy query
// https://github.com/Stratio/cassandra-lucene-index/blob/branch-3.0.10/doc/documentation.rst#fuzzy-search
func Fuzzy(f string, v interface{}) Rule {
	return &fuzzyQuery{
		Type:  "fuzzy",
		Field: f,
		Value: v,
	}
}

// Regexp query
// https://github.com/Stratio/cassandra-lucene-index/blob/branch-3.0.10/doc/documentation.rst#regexp-search
func Regexp(f string, v string) Rule {
	return &value{
		Type:  "regexp",
		Field: f,
		Value: v,
	}
}

// Wildcard query <field> <value>
// https://github.com/Stratio/cassandra-lucene-index/blob/branch-3.0.10/doc/documentation.rst#wildcard-search
func Wildcard(f string, v string) Rule {
	return &value{
		Type:  "wildcard",
		Field: f,
		Value: v,
	}
}

type phraseQuery struct {
	Type   string      `json:"type"`
	Field  string      `json:"field"`
	Values interface{} `json:"values"`
	Slop   int         `json:"slop,omitempty"`
}

// Phrase query
// see: https://github.com/Stratio/cassandra-lucene-index/blob/branch-3.0.10/doc/documentation.rst#phrase-search
// slop (default = 0): number of words permitted between words.
func Phrase(f string, v interface{}, slop int) Rule {
	return &phraseQuery{
		Type:   "phrase",
		Field:  f,
		Values: v,
		Slop:   slop,
	}
}

// Contains query
func Contains(f string, v interface{}) Rule {
	return &values{
		Type:   "contains",
		Field:  f,
		Values: v,
	}
}

// rangeValue struct
type rangeValue struct {
	Type         string      `json:"type"`
	Field        string      `json:"field"`
	Lower        interface{} `json:"lower,omitempty"`
	Upper        interface{} `json:"upper,omitempty"`
	IncludeLower bool        `json:"include_lower,omitempty"`
	IncludeUpper bool        `json:"include_upper,omitempty"`
}

// RangeAll query
func RangeAll(f string, lv, uv interface{}, incl, incu bool) Rule {
	return &rangeValue{
		Type:         "range",
		Field:        f,
		Lower:        lv,
		Upper:        uv,
		IncludeLower: incl,
		IncludeUpper: incu,
	}
}

// RangeLower query
func RangeLower(f string, v interface{}, inc bool) Rule {
	return &rangeValue{
		Type:         "range",
		Field:        f,
		Lower:        v,
		IncludeLower: inc,
	}
}

// RangeUpper query
func RangeUpper(f string, v interface{}, inc bool) Rule {
	return &rangeValue{
		Type:         "range",
		Field:        f,
		Upper:        v,
		IncludeUpper: inc,
	}
}

// BooleanQuery struct
type BooleanQuery struct {
	Type   string `json:"type"`
	Must   []Rule `json:"must,omitempty"`
	Should []Rule `json:"should,omitempty"`
	Not    []Rule `json:"not,omitempty"`
}

// // Lucene to match Rule interface
// func (v *BooleanQuery) Lucene() bool {
// 	return true
// }

// BooleanMust query
func BooleanMust(rs ...Rule) Rule {
	v := &BooleanQuery{
		Type: "boolean",
		Must: make([]Rule, len(rs)),
	}
	for i, r := range rs {
		v.Must[i] = r
	}
	return v
}

// BooleanShould query
func BooleanShould(rs ...Rule) Rule {
	v := &BooleanQuery{
		Type:   "boolean",
		Should: make([]Rule, len(rs)),
	}
	for i, r := range rs {
		v.Should[i] = r
	}
	return v
}

// BooleanNot query
func BooleanNot(rs ...Rule) Rule {
	v := &BooleanQuery{
		Type: "boolean",
		Not:  make([]Rule, len(rs)),
	}
	for i, r := range rs {
		v.Not[i] = r
	}
	return v
}

// Refresh enable / disable refreshing index
func (x *Expr) Refresh(e bool) *Expr {
	x.R = e
	return x
}

// MarshalCQL enables gocql library marshal into cql
func (x *Expr) MarshalCQL(_ gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(x)
}

// String for stringer interface
func (x *Expr) String() string {
	r, err := json.Marshal(x)
	if err != nil {
		return ""
	}
	return string(r)
}

// NewExpr creates a new expression
func NewExpr() *Expr {
	return &Expr{
		F: make([]Rule, 0),
		Q: make([]Rule, 0),
		S: make([]*Sort, 0),
	}
}

// ResetQuery removes pre-configured query expression
func (x *Expr) ResetQuery() *Expr {
	x.Q = make([]Rule, 0)
	return x
}

// ResetFilter removes pre-configured filter expression
func (x *Expr) ResetFilter() *Expr {
	x.F = make([]Rule, 0)
	return x
}

// ResetSort removes pre-configured sort expression
func (x *Expr) ResetSort() *Expr {
	x.S = make([]*Sort, 0)
	return x
}

// Reset all pre-configured options
func (x *Expr) Reset() *Expr {
	return x.ResetFilter().ResetQuery().ResetSort().Refresh(false)
}

// Query add queries to expression
func (x *Expr) Query(q ...Rule) *Expr {
	for _, o := range q {
		switch {
		case o == nil:
			continue
		default:
			x.Q = append(x.Q, o)
		}
	}
	return x
}

// Filter add filters to expression
func (x *Expr) Filter(q ...Rule) *Expr {
	for _, o := range q {
		switch {
		case o == nil:
			continue
		default:
			x.F = append(x.F, o)
		}
	}
	return x
}

// SortBy with field
func (x *Expr) SortBy(f string, r bool) *Expr {
	if len(x.S) > 0 {
		for _, s := range x.S {
			if s.Field == f { // modify if exists
				s.Reverse = r
				return x
			}
		}
	}
	x.S = append(x.S, &Sort{Field: f, Reverse: r})
	return x
}
