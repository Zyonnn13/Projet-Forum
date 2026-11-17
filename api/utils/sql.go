package utils

import (
	"strconv"
	"strings"
)

// Rebind convertit les ? (MySQL) en $1, $2… (PostgreSQL).
func Rebind(query string) string {
	var b strings.Builder
	n := 0
	for _, c := range query {
		if c == '?' {
			n++
			b.WriteString("$" + strconv.Itoa(n))
			continue
		}
		b.WriteRune(c)
	}
	return b.String()
}
