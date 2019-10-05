// Package csvcolumn reads specified columns from a csv style input.
package csvcolumn

import (
	"encoding/csv"
	"io"
	"strings"
)

type Column struct {
	Name  string
	Index int
	Value Value
}

const (
	// IndexUnbound is an init value of a Column.Index
	IndexUnbound = -2
	// IndexMissing is an Column.Index value used when Column.Name was not
	// found in the headers of the source
	IndexMissing = -1
)

// NewReader is a factory function to create a *Reader
func NewReader(src io.Reader) *Reader {
	return &Reader{
		Comma: ',',
		r:     src,
	}
}

// Reader keeps the configuration and state of reading from source
type Reader struct {
	Comma               rune
	Comment             rune
	FieldsPerRecord     int
	LazyQuotes          bool
	TrimLeadingSpace    bool
	CaseSensitiveHeader bool

	r       io.Reader
	err     error
	parser  *csv.Reader
	headers []string
	columns []Column
}

func (r *Reader) init() {
	r.parser = csv.NewReader(r.r)
	r.parser.Comma = r.Comma
	r.parser.Comment = r.Comment
	r.parser.FieldsPerRecord = r.FieldsPerRecord
	r.parser.LazyQuotes = r.LazyQuotes
	r.parser.TrimLeadingSpace = r.TrimLeadingSpace
	r.parser.ReuseRecord = true
	r.bind()
}

// bind parses the headers from source to validate columns of interest
// are present in a header.
func (r *Reader) bind() {
	if len(r.headers) == 0 {
		r.headers, r.err = r.parser.Read()
		if r.err != nil {
			return
		}
	}

	for i := range r.columns {
		column := &r.columns[i]
		if column.Index >= 0 || column.Index == IndexMissing {
			continue
		}

		if r.CaseSensitiveHeader {
			for k, header := range r.headers {
				if column.Name == header {
					column.Index = k
					break
				}
			}
		} else {
			for k, header := range r.headers {
				if strings.EqualFold(column.Name, header) {
					column.Index = k
					break
				}
			}
		}

		if column.Index == IndexUnbound {
			column.Index = IndexMissing
		}
	}
}

// Next parses a next line from a source.
func (r *Reader) Next() (ok bool) {
	if r.parser == nil {
		r.init()
		if r.err != nil {
			return false
		}
	}

	record, err := r.parser.Read()
	if err == io.EOF {
		return false
	}
	if err != nil && r.err == nil {
		r.err = err
		return false
	}

	for i := range r.columns {
		column := &r.columns[i]
		if column.Index < 0 {
			continue
		}

		err := column.Value.Scan(record[column.Index])
		if err != nil && r.err == nil {
			r.err = err
			return false
		}
	}

	return true
}

// Err returns the latest error of a reader.
func (r *Reader) Err() error { return r.err }

func (r *Reader) Bind(column string, value Value) {
	if r.parser != nil {
		panic("binding must be done before calling Next")
	}
	r.columns = append(r.columns, Column{Name: column, Index: IndexUnbound, Value: value})
}

func (r *Reader) StringColumn(column string) *String {
	value := &String{}
	r.Bind(column, value)
	return value
}

func (r *Reader) IntColumn(column string) *Int {
	value := &Int{}
	r.Bind(column, value)
	return value
}

// String returns a pointer to a string field value of a
//  * row where the reader currently is reading from
//  * column passed as an argument
func (r *Reader) String(column string) *string { return &r.StringColumn(column).Value }

// Int returns a pointer to an int field value of a
//  * row where the reader currently is reading from
//  * column passed as an argument
func (r *Reader) Int(column string) *int { return &r.IntColumn(column).Value }
