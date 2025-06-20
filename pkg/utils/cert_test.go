package utils

import (
	"fmt"
	"testing"
)

func TestCertGen(t *testing.T) {
	cert, key, err := DefaultGenCert.GenCert()
	if err != nil {
		return
	}
	fmt.Println("cert:", string(cert))
	fmt.Println("key:", string(key))
}
