package rtnetlink

import (
	"fmt"
	"strings"
)

func bPrint(b []byte) string {
	return strings.ReplaceAll(fmt.Sprintf("%# x", b), " ", ", ")
}
