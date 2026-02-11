package zipSm

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	encrypt, err := Encrypt([]byte(`myyixinfixkdsxdo`))
	fmt.Println(encrypt, err)
	fmt.Println(len(encrypt))
	decrypt, err := Decrypt(encrypt)
	fmt.Println(string(decrypt), err)
	fmt.Println(len(decrypt))

	for i := 0; i < 30; i++ {

		encrypt, err := Encrypt([]byte(`myyixinfixkdsxdo`))
		fmt.Println(encrypt, err)
	}

}
