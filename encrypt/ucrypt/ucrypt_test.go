package ucrypt

import (
	"fmt"
	"testing"
)

func TestEnc(t *testing.T) {

	fmt.Println(Encrypt2("2025-12-15", "1901547647767793664", 5728))

}

func TestDec(t *testing.T) {

	fmt.Println(Decrypt2("NmeFIVQuLvkq3B", "1901547647767793664", 5728))

}

func TestJf(t *testing.T) {
	//1405070
	fmt.Println(14 + 2688 + 94208 + 1310720)
}
