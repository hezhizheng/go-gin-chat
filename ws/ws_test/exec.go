package main

import (
	"log"
	"os/exec"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}
var ch = make(chan int, 20)

func main() {

	for i := 500; i <= 600; i++ {
		wg.Add(1)
		go execCommand(i)
	}

	wg.Wait()

	log.Println("okkkkkkkkkkkkkkkkkk")
}

func execCommand(i int) {

	defer func() {
		//捕获read抛出的panic
		if err := recover();err!=nil{
			log.Println("execCommand",err)
		}
	}()

	ch <- i
	strI := strconv.Itoa(i)
	//cmd := exec.Command("./mock_ws_client_coon.exe", strI)
	cmd := exec.Command("E:\\go1.15.2.windows-386\\go\\bin\\go.exe",
		"run",
		"D:\\phpstudy_pro\\WWW\\org\\public-go-gin-chat\\ws\\ws_test\\mock_ws_client_coon.go",
		strI)

	err := cmd.Start()

	if err != nil {
		log.Println(err)
	}

	log.Println(i)

	//time.Sleep(time.Second * 1)
	<-ch
	wg.Done()
}
