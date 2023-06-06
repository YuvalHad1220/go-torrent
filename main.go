package main

import (
	"fmt"
	"os"
	"time"
	"sync"
	"path"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func handleTorrent(t *Torrent){
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	  panic("failed to connect database")
	}
  
	// Migrate the schema
	db.AutoMigrate(&Torrent{})
  
	// Create
	db.Create(t)
  
}


func ParseTorrent(TorrentFileName string) (*Torrent, error) {
	content, err := os.ReadFile(TorrentFileName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n",TorrentFileName, err)
        return nil, err
    }
	decoded, _ := decode(content)

	t, err := ParseTorrentFromBencoded(decoded)

	if err!= nil {
        return nil, err
    }

	return t, nil
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
		err := os.MkdirAll(moveTo, os.ModePerm)
        if err!= nil {
            fmt.Println(err)
            return
        }
		fmt.Println("Directory created")
    }

	_, err = os.Stat(erroredTorrents)
	if os.IsNotExist(err) {
		fmt.Println("Errored Torrent directiory does not exist, creating...")
		err := os.MkdirAll(erroredTorrents, os.ModePerm)
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
			t, err := ParseTorrent(old_path)
			handleTorrent(t)
			if err!= nil {
				error_path := path.Join(erroredTorrents, file.Name())
				fmt.Println("Error parsing torrent file:", err)
				err = MoveFile(old_path, error_path)
				if err!= nil {
					fmt.Println("Error moving file:", err)
                }

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
	go MonitorDirectory("NewTorrentFiles", "AddedTorrentFiles", "ErroredTorrents", &wg)
	wg.Wait()

}