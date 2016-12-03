package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/ieee0824/goaa"
	"github.com/ieee0824/goaa/sls/util"
)

var (
	rate  = flag.Int("r", 15, "frame rate")
	chank = flag.Int("c", 10, "chank size")
	input = flag.String("i", "", "input path")
	out   = flag.String("o", "out", "output path")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	infos, err := ioutil.ReadDir(*input)
	if err != nil {
		log.Fatalln(err)
	}

	queue := make(chan string, *rate**chank*2)
	conf := utli.SlsConf{*rate, *chank}
	bin, err := json.Marshal(conf)
	if err != nil {
		log.Fatalln(err)
	}
	if err = ioutil.WriteFile(*out+"/conf.json", bin, 0644); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for _, info := range infos {
			file, err := os.Open(strings.TrimSuffix(*input, "/") + "/" + info.Name())
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
			img, _, err := image.Decode(file)
			if err != nil {
				log.Fatalln(err)
			}
			imgstrs, err := goaa.ConvertASCII(img)
			if err != nil {
				log.Fatalln(err)
			}
			queue <- strings.Join(imgstrs, "\n")
		}
		close(queue)
	}()

	var c util.Container
	counter := 0
	for {
		if frame, ok := <-queue; ok {
			c.Data = append(c.Data, frame)
		} else {
			if len(c.Data) != 0 {
				bin, err := json.Marshal(c)
				if err != nil {
					log.Fatalln(err)
				}
				err = ioutil.WriteFile(fmt.Sprintf(*out+"/%05d.json.gz", counter), compress(bin), 0644)
				if err != nil {
					log.Fatalln(err)
				}
			}
			break
		}
		if len(c.Data) == *chank {
			bin, err := json.Marshal(c)
			if err != nil {
				log.Fatalln(err)
			}
			err = ioutil.WriteFile(fmt.Sprintf(*out+"/%05d.json.gz", counter), compress(bin), 0644)
			if err != nil {
				log.Fatalln(err)
			}
			c = util.Container{}
			counter++
		}
	}
}

func compress(bin []byte) []byte {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(bin)
	w.Close()
	return b.Bytes()
}
