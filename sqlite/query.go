package sqlite

import (
	"fmt"
	"strings"
)

// GenerateMultiLikes genrates segments of sql for using same files but with diffrent LIKE value
func GenerateMultiLikes(field string, n int, and bool) (string, []string) {
	var sb strings.Builder
	likes := make([]string, 0)

	sb.WriteString("(")

	for i := 0; i < n; i++ {
		value := fmt.Sprintf("$like%d", i)
		likes = append(likes, value)

		if i != 0 {
			if and {
				sb.WriteString(" AND ")
			} else {
				sb.WriteString(" OR ")
			}
		}

		sb.WriteString(field)
		sb.WriteString(" LIKE ")
		sb.WriteString(value)
	}

	sb.WriteString(")")

	return sb.String(), likes
}

// GenerateInValues generates type safe values for IN ( ... )
func GenerateInValues(name string, n int) (string, []string) {
	var sb strings.Builder
	params := make([]string, 0)

	for i := 0; i < n; i++ {
		param := fmt.Sprintf("$%s%d", name, i)
		params = append(params, param)
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(param)
	}

	return sb.String(), params
}
