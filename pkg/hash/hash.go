package hash

import (
	"crypto/sha256"
	"fmt"
)

func Encode(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
