/* 
bmp driver
*/

package main

import (
	"os"
	"fmt"
	"github.com/Ysgard/bitmap"
)

func main() {
	b, err := bitmap.NewBitmap("CDtest.BMP")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Cannot load CDTest.BMP")
		return
	}
	b.Info()


}