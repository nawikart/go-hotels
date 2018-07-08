package custom

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func SaveThumbFromUrl(fn string, url string) {
	img, _ := os.Create(fn)
	defer img.Close()

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	b, _ := io.Copy(img, resp.Body)
	fmt.Println(fn, "created, file size: ", b)
}
