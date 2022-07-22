package findjson

import (
	"fmt"
)

func ExampleFindJson() {
	s := []byte(`<script>
		a = [1, 2, true, [5, 8, 13], null, 21]
		b = {"one": 1, "two": 2, "three": 3}
	</script>`)

	i := 0
	for i < len(s) {
		start, end, err := FindJson(s, i, JsonValueBoolean|JsonValueObject)
		if err == nil {
			fmt.Println(string(s[start:end]))
		}

		i = end
	}

	// Output: true
	// {"one": 1, "two": 2, "three": 3}
}
