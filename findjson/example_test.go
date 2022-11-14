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

	// Output:
	// true
	// {"one": 1, "two": 2, "three": 3}
}

func ExampleFindJsonWithStyle_normative_style() {
	s := []byte(`<script>
		a = [1, 2, 3]
		b = [4, 5, 6,]
		c = {"one": 1, "two": 2, "three": 3}
	</script>`)

	i := 0
	for i < len(s) {
		start, end, err := FindJsonWithStyle(s, i, JsonValueArray, NormativeStyle)
		if err == nil {
			fmt.Println(string(s[start:end]))
		}

		i = end
	}

	// Output:
	// [1, 2, 3]
}

func ExampleFindJsonWithStyle_javascript_style() {
	s := []byte(`<script>
		a = [1, 2, 3]
		b = [4, 5, 6,]
		c = {"one": 1, "two": 2, "three": 3}
	</script>`)

	i := 0
	for i < len(s) {
		start, end, err := FindJsonWithStyle(s, i, JsonValueArray, JavaScriptStyle)
		if err == nil {
			fmt.Println(string(s[start:end]))
		}

		i = end
	}

	// Output:
	// [1, 2, 3]
	// [4, 5, 6,]
}
