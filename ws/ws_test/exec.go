package main

import (
	"log"
	"os/exec"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}
var ch = make(chan int,100)

func main()  {

	for i:=500 ;i<=1500; i++ {
		wg.Add(1)
		go funcName(i)
	}

	wg.Wait()
}

func funcName(i int) {
	ch <- i
	strI := strconv.Itoa(i)
	cmd := exec.Command("./mock_ws_client_coon.exe", strI)
	err := cmd.Start()
	o, _ := cmd.Output()
	log.Println(string(o))

	if err != nil {
		log.Println(err)
	}

	<-ch
	wg.Done()
}