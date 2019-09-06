package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

func main() {
	for {
		fmt.Println("Hacking...")
		epo := strconv.Itoa(int(time.Now().Unix())) + ".virus"
		d1 := []byte("hehexd\n")
		ioutil.WriteFile(epo, d1, 0644)
		time.Sleep(3 * time.Second)
	}
}
