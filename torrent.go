package main
import (
	"fmt"
	"crypto/sha1"
	
)
type Torrent struct {
	ID          uint `gorm:"primarykey"`
	TorrentName string
	Infohash    []byte
	Size        uint
	Downloaded  uint
	Uploaded    uint

	AnnounceUrl    string
	TimeToAnnounce int
	Peers          []string
	Seeders        uint
	Leechers       uint

	PiecesData  []byte
	PieceLength uint
	PieceCount  uint
}

func GetInfoHash(info_entry map[any]any) []byte {
	hash := sha1.New()
	// Write the data to the hash object
	hash.Write(encode(info_entry))
	// Get the SHA-1 hash sum
	hashSum := hash.Sum(nil)
	return hashSum
}

func GetAnnounceUrl(decoded map[any]any) string {
	AnnounceUrl := decoded["announce"]
	return string(AnnounceUrl.([]byte))
}

func GetSize(info_entry map[any]any) (int, error) {
	value, ok := info_entry["files"]

	if ok {
		files, ok := value.([]any)
		if !ok {
			return 0, fmt.Errorf("invalid files entry")
			}

		size := int(0)

		for _, file := range files {
			file_as_dict, ok := file.(map[any]any)
			if !ok {
				return 0, fmt.Errorf("invalid files entry")
			}
			size += file_as_dict["length"].(int)
		}
		return size, nil
	}

	return info_entry["length"].(int), nil

}

// func GetName(info_entry map[any]any) (string, error) {
// }

// func GetPiecesData(info_entry map[any]any) ([]byte, error) {
// }

// func GetPiecesLength(info_entry map[any]any) (uint, error){

// }

// func GetPieceCount(info_entry map[any]any) (uint, error){

// }



// func ParseTorrentFromBencoded(decoded any) (*Torrent, error) {	
// 	decoded_as_map, ok := decoded.(map[any]any)
// 	if !ok {
//         return nil, fmt.Errorf("invalid bencoded torrent")
//     }

// 	info_entry, ok := decoded_as_map["info"]
// 	if !ok {
//         return nil, fmt.Errorf("Cant get info entry")
//     }




// 	return &Torrent{}, nil

// }