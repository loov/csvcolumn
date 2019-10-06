// Package csvcolumn reads specified columns from a csv style input.
package csvcolumn

import (
	"encoding/csv"
	"io"
	"strings"
)

type field struct {
	Name   string
	Column int
	Value  Value
}

const (
	// ColumnUnbound is an init value of a field.Column
	ColumnUnbound = -2
	// ColumnMissing is an field.Column value used when field.Name was not
	// found in the headers of the source
	ColumnMissing = -1
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
	fields  []field
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

// bind initializes the Reader reading state.
// It parses the headers from source to validate columns of interest
// are present in a header.
func (r *Reader) bind() {
	if len(r.headers) == 0 {
		r.headers, r.err = r.parser.Read()
		if r.err != nil {
			return
		}
	}

	for i := range r.fields {
		field := &r.fields[i]
		if field.Column != ColumnUnbound {
			continue
		}

		if r.CaseSensitiveHeader {
			for k, header := range r.headers {
				if field.Name == header {
					field.Column = k
					break
				}
			}
		} else {
			for k, header := range r.headers {
				if strings.EqualFold(field.Name, header) {
					field.Column = k
					break
				}
			}
		}

		if field.Column == ColumnUnbound {
			field.Column = ColumnMissing
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

	for i := range r.fields {
		field := &r.fields[i]
		if field.Column < 0 {
			continue
		}

		cell := record[field.Column]
		err := field.Value.Scan(cell)
		if err != nil && r.err == nil {
			r.err = err
			return false
		}
	}

	return true
}

// Err returns the latest error of a reader.
func (r *Reader) Err() error { return r.err }

// Bind binds a column in a csv to its matching Value struct.
// It's useful when you are interested in adding your own Value types.
func (r *Reader) Bind(columnName string, value Value) {
	if r.parser != nil {
		panic("binding must be done before calling Next")
	}
	r.fields = append(r.fields,
		field{
			Name:   columnName,
			Column: ColumnUnbound,
			Value:  value,
		},
	)
}

// String returns a pointer to a string field value of a given columnName
// which is reassigned with every Next() call.
func (r *Reader) String(columnName string) *string {
	value := &String{}
	r.Bind(columnName, value)
	return &value.Value
}

// Int returns a pointer to an int field value of a given columnName which
// is reassigned with every Next() call.
func (r *Reader) Int(columnName string) *int {
	value := &Int{}
	r.Bind(columnName, value)
	return &value.Value
}
