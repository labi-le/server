package file

import "strings"

const (
	alphabet    = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"
	alphabetLen = len(alphabet)
)

func Short(id int) string {
	var digits []int
	for id > 0 {
		digits = append(digits, id%alphabetLen)
		id /= alphabetLen
	}

	// reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	var b strings.Builder
	for _, digit := range digits {
		b.WriteString(string(alphabet[digit]))
	}

	return b.String()
}
