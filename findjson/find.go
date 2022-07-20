package findjson

func FindJson(s []byte, i int, kind JsonValueKind) (int, int, error) {
	l := len(s)
	j := i
	for j < l {
		if isWhiteSpace(s[j]) {
			j = jumpNextNonWhiteSpace(s, j)
		}

		c := s[j]
		if scanner := kind.GetScanner(c); scanner != nil {
			start, end, err := scanner(s, j)
			return start, end, err
		}

		j++
	}

	return i, j, NewJsonError(j, "no JSON string found in %s", kind)
}