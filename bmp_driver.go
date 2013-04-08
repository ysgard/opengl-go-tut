/* 
bmp driver
*/

package main

import (
	"fmt"
	"os"
)

func main() {
	b, err := NewBitmap("CDtest.BMP")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Cannot load CDTest.BMP")
		return
	}
	b.Info()

}
