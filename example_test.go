package csvcolumn_test

import (
	"fmt"
	"strings"

	"github.com/loov/csvcolumn"
)

func ExampleReader() {
	source := strings.NewReader(`Index,Age,Height,Name
1,52,1.7,Alice
5,42,1.88,Bob
512,31,1.82,Charlie`)

	data := csvcolumn.NewReader(source)
	data.LazyQuotes = true
	name, age, height := data.String("Name"), data.Int("Age"), data.Float64("Height")

	for data.Next() && data.Err() == nil {
		fmt.Println(*name, *age, *height)
	}

	if data.Err() != nil {
		fmt.Println(data.Err())
	}

	//Output:Alice 52 1.7
	//Bob 42 1.88
	//Charlie 31 1.82
}
