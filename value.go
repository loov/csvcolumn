package csvcolumn

import "strconv"

// Value interface defines functions for
//  * creating a value from string
//  * returning a string representation of a value
type Value interface {
	Scan(text string) error
	String() string
}

// String is a struct implementing the Value interface
type String struct {
	Value string
}

func (value *String) Scan(text string) error {
	value.Value = text
	return nil
}

func (value *String) String() string {
	return value.Value
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
