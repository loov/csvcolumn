package csvcolumn_test

import (
	"fmt"
	"strings"

	"github.com/loov/csvcolumn"
)

func ExampleReader() {
	source := strings.NewReader(`Index,Age,Name
1,52,Alice
5,42,Bob
512,31,Charlie`)

	data := csvcolumn.NewReader(source)
	data.LazyQuotes = true
	name, age := data.String("Name"), data.Int("Age")

	for data.Next() && data.Err() == nil {
		fmt.Println(*name, *age)
	}

	if data.Err() != nil {
		fmt.Println(data.Err())
	}

	//Output:Alice 52
	//Bob 42
	//Charlie 31
}
