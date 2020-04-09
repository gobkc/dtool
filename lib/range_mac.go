package lib

import (
	"crypto/rand"
	"fmt"
)

func GetRandMac() string {
	b := make([]byte, 6)
	rand.Read(b)
	mac := fmt.Sprintf("00:%02x:%02x:%02x:%02x:%02x", b[1], b[2], b[3], b[4], b[5])
	return mac
}