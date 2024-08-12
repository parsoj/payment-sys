package dict

import (
	"fmt"
	"testing"
)

func TestRandomWordsDictionary(t *testing.T) {
	t.Skip("Not using this feature")

	d, err := NewRandomWordsDictionary()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i <= 20; i++ {
		s, err := d.GetRandomString()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(s)

		}

	}

}
