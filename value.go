package csvcolumn

import "strconv"

// Value interface describes a way how different data types should be
// scanned.
type Value interface {
	// Scan assigns a value from text
	//
	// An error should be returned if the value cannot be stored
	// without loss of information.
	Scan(text string) error
}

// String is a struct implementing the Value interface
type String struct {
	Value string
}

func (value *String) Scan(text string) error {
	value.Value = text
	return nil
}

// Int is a struct implementing the Value interface
type Int struct {
	Value   int
	Default int
	Error   error
}

func (value *Int) Scan(text string) error {
	value.Value, value.Error = strconv.Atoi(text)
	return value.Error
}

func (value *Int) String() string {
	return strconv.Itoa(value.Value)
}

// Float64 is a struct implementing the Value interface
type Float64 struct {
	Value   float64
	Default float64
	Error   error
}

func (value *Float64) Scan(text string) error {
	value.Value, value.Error = strconv.ParseFloat(text, 64)
	return value.Error
}

func (value *Float64) String() string {
	return strconv.FormatFloat(value.Value, 'g', -1, 64)
}
