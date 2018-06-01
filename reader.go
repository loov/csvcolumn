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
	IndexUnbound = -2
	IndexMissing = -1
)

func NewReader(src io.Reader) *Reader {
	return &Reader{
		Comma: ',',
		r:     src,
	}
}

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

func (r *Reader) Next() bool {
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
	}

	for i := range r.columns {
		column := &r.columns[i]
		if column.Index < 0 {
			continue
		}

		err := column.Value.Scan(record[column.Index])
		if err != nil && r.err == nil {
			r.err = err
		}
	}

	return true
}

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

func (r *Reader) String(column string) *string { return &r.StringColumn(column).Value }
func (r *Reader) Int(column string) *int       { return &r.IntColumn(column).Value }
