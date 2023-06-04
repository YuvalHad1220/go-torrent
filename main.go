package main

import (
	"fmt"
	"os"
	"time"
	"sync"
)

func MonitorDirectory(dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println("Directory does not exist, creating...")
		err := os.MkdirAll(dir, os.ModePerm)
        if err!= nil {
            fmt.Println(err)
            return
        }
		fmt.Println("Directory created")
    }


	for {
		files, err := os.ReadDir(dir)
		if err!= nil {
			fmt.Println(err)
			return
		}
	
		for _, file := range files {
			fmt.Println("New File:", file.Name())
		}
		time.Sleep(time.Second * 10)
	}


}








func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go MonitorDirectory("NewTorrentFiles", &wg)
	wg.Wait()

}