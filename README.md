# csvcolumn

[![GoDoc](https://godoc.org/github.com/loov/csvcolumn?status.svg)](http://godoc.org/github.com/loov/csvcolumn)

csvcolumn package implements convenient CSV reading with column access.

Internally it uses https://golang.org/pkg/encoding/csv/,
which means it inherits all the same restrictions.


``` go
const CSV = `Index,Age,Name
1,52,Alice
5,42,Bob
512,31,Charlie
`

func Example() {
	source := strings.NewReader(CSV)

	data := csvcolumn.NewReader(source)
	data.LazyQuotes = true
	name, age := data.String("Name"), data.Int("Age")

	for data.Next() && data.Err() == nil {
		fmt.Println(*name, *age)
	}

	if data.Err() != nil {
		fmt.Println(data.Err())
	}
}
```
