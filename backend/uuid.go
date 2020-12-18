package backend

import (
	"math/rand"
	"strconv"
)

func uuid() string {
	uuid := ""
	for i := 0; i < 32; i++ {
		switch i {
		case 8:
		case 20:
			uuid += "-"
			uuid += strconv.Itoa(rand.Intn(100)*16 | 0)
			break
		case 12:
			uuid += "-"
			uuid += "4"
			break
		case 16:
			uuid += "-"
			uuid += strconv.Itoa(rand.Intn(100)*4 | 8)
			break
		default:
			uuid += strconv.Itoa(rand.Intn(100)*16 | 0)
		}
	}
	return uuid
}
