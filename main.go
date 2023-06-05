package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"path"
)



func ParseTorrent(TorrentFileName string) (*Torrent, error) {
	content, err := os.ReadFile(TorrentFileName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n",TorrentFileName, err)
        return nil, err
    }
	decoded, _ := decode(content)

	decoded_as_map, ok := decoded.(map[any]any)
		if !ok {
			fmt.Println("error decoding as map")
			return nil, fmt.Errorf("%s is not a map",TorrentFileName)
		}

	res := GetAnnounceUrl(decoded_as_map)
	fmt.Println(res)

	return nil, nil
}

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
			ParseTorrent(path.Join(dir, file.Name()))
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