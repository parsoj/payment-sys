package debug

import (
	"fmt"
)

func PrintWithLength(s string) {

	fmt.Printf("%s :: %d\n", s, len(s))

}
