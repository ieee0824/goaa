package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/ieee0824/goaa"
	"github.com/nfnt/resize"
)

func main() {
	file, err := os.Open("test.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	img, _, err := image.Decode(file)
	img = resize.Resize(0, 400, img, resize.Lanczos3)

	if err != nil {
		log.Fatalln(err)
	}
	imgstr, err := goaa.ConvertASCII(img)
	if err != nil {
		log.Fatalln(err)
	}

	s := strings.Join(imgstr, "\n")
	fmt.Println(s)
}
