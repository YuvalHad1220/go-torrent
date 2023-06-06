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

	info_entry, ok := decoded_as_map["info"].(map[any]any)
	if !ok {
		return nil, fmt.Errorf("%s does not have info entry",TorrentFileName)	
	}


	Size, err := GetSize(info_entry)
	// AnnounceUrl := GetAnnounceUrl(decoded_as_map)
	// InfoHash := GetInfoHash(decoded_as_map)

	fmt.Println(Size)

	return nil, nil
}

func MoveFile(oldDir string, newDir string) error {
	err := os.Rename(oldDir, newDir)
	if err != nil {
		return err
	}
	return nil

}

func MonitorDirectory(listenIn string, moveTo string, erroredTorrents string, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := os.Stat(listenIn)
	if os.IsNotExist(err) {
		fmt.Println("Listening-in directory does not exist, creating...")
		err := os.MkdirAll(listenIn, os.ModePerm)
        if err!= nil {
            fmt.Println(err)
            return
        }
		fmt.Println("Directory created")
    }

	_, err = os.Stat(moveTo)
	if os.IsNotExist(err) {
		fmt.Println("Moving-to directiory does not exist, creating...")
		err := os.MkdirAll(listenIn, os.ModePerm)
        if err!= nil {
            fmt.Println(err)
            return
        }
		fmt.Println("Directory created")
    }

	for {
		files, err := os.ReadDir(listenIn)
		if err!= nil {
			fmt.Println(err)
			return
		}
	
		for _, file := range files {
			old_path := path.Join(listenIn, file.Name())
			fmt.Println("New File:", file.Name())
			_, err := ParseTorrent(old_path)

			if err!= nil {
				error_path := path.Join(erroredTorrents, file.Name())
				fmt.Println("Error parsing torrent file:", err)
				err = MoveFile(old_path, error_path)

            } else {
				new_path := path.Join(moveTo, file.Name())
                err = MoveFile(old_path, new_path)
                if err!= nil {
					fmt.Println("Error moving file:", err)
                }
            }
		}
		
		time.Sleep(time.Second * 10)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go MonitorDirectory("NewTorrentFiles", "AddedTorrentFiles", "ErrorTorrents", &wg)
	wg.Wait()

}