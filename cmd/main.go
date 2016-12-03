package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/ieee0824/goaa"
)

func main() {
	file, err := os.Open("test3.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	img, _, err := image.Decode(file)
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
