package generathion

import (
	"math/rand"
	"time"
)

func GenarathionCode(lenght int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, lenght)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}

	return string(code)
}
