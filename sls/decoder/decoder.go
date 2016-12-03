package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	tm "github.com/buger/goterm"
	"github.com/ieee0824/goaa/sls/util"
)

type downloadData struct {
	data []byte
	name string
}

var client = &http.Client{}

func getList() ([]string, error) {
	ret := []string{}
	listApi := "http://localhost:8080/list"
	resp, err := client.Get(listApi)
	if err != nil {
		return nil, err
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.Unmarshal(bin, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func getMov(wg *sync.WaitGroup, q chan string, dl chan downloadData) {
	movApi := "http://localhost:8080/static/"
	defer wg.Done()
	for {
		fileName, ok := <-q
		if !ok {
			return
		}
		resp, err := client.Get(movApi + fileName)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		bin, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		dl <- downloadData{bin, fileName}
	}
}

func deCompress(bin []byte) ([]byte, error) {
	reader := bytes.NewReader(bin)
	fz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer fz.Close()
	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func main() {

	var wg sync.WaitGroup
	list, _ := getList()

	q := make(chan string, 1)
	dl := make(chan downloadData, 1)

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go getMov(&wg, q, dl)
	}

	go func() {
		for _, name := range list {
			q <- name
		}
		close(q)
	}()

	counter := 0
	for {
		var c util.Container
		if m, ok := <-dl; ok {
			movRow, err := deCompress(m.data)
			if err != nil {
				log.Fatalln(err)
			}
			err = json.Unmarshal(movRow, &c)
			if err != nil {
				log.Fatalln(err)
			}
			for _, s := range c.Data {
				tm.MoveCursor(0, 0)
				fmt.Println(s)
				time.Sleep(10 * time.Millisecond)
				tm.Flush()
			}

			counter++

		} else {
			break
		}
		if counter == len(list) {
			break
		}
	}
}
