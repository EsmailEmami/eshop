package random_code

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomDigit(length uint) string {
	rand.Seed(time.Now().UnixNano())
	str := ""
	for i := 1; i <= int(length); i++ {
		str += fmt.Sprintf("%d", rand.Intn(10))
	}
	return str
}
