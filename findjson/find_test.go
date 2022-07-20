package findjson

import (
	"testing"
)

func TestFindJsonNumbersInArray(t *testing.T) {
	//           0         1         2         3
	//           0123456789012345678901234567890
	s := []byte(`[1, 1, 2, 3, 5, 8, 13, 21]`)

	got := make([]string, 0)
	i := 0
	for i < len(s) {
		start, end, err := FindJson(s, i, JsonValueNumber)
		if err == nil {
			got = append(got, string(s[start:end]))
		}

		i = end
	}

	exp := []string{
		"1", "1", "2", "3", "5", "8", "13", "21",
	}

	for i, v := range exp {
		if got[i] != v {
			t.Errorf("expected %s, got %s", v, got[i])
		}
	}
}

func TestFindJsonNumbersInTokenArray(t *testing.T) {
	//           0         1         2         3
	//           0123456789012345678901234567890
	s := []byte(`[1, 1, 2, null, 5, true, [], 21]`)

	got := make([]string, 0)
	i := 0
	for i < len(s) {
		start, end, err := FindJson(s, i, JsonValueNumber)
		if err == nil {
			got = append(got, string(s[start:end]))
		}

		i = end
	}

	exp := []string{
		"1", "1", "2", "5", "21",
	}

	for i, v := range exp {
		if got[i] != v {
			t.Errorf("expected %s, got %s", v, got[i])
		}
	}
}

func TestFindJsonBooleanInTokenArray(t *testing.T) {
	//           0         1         2         3
	//           0123456789012345678901234567890
	s := []byte(`[1, 1, 2, null, 5, true, [], 21, trust, flag]`)

	got := make([]string, 0)
	i := 0
	for i < len(s) {
		start, end, err := FindJson(s, i, JsonValueBoolean)
		if err == nil {
			got = append(got, string(s[start:end]))
		}

		i = end
	}

	exp := []string{
		"true",
	}

	for i, v := range exp {
		if got[i] != v {
			t.Errorf("exp[%d](%s) != got[%d](%s)", i, v, i, got[i])
		}
	}
}
