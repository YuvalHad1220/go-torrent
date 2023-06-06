package main
import (
	"fmt"
	"crypto/sha1"
	
)
type Torrent struct {
	ID          int `gorm:"primarykey"`
	TorrentName string
	Infohash    []byte
	Size        int
	Downloaded  int
	Uploaded    int

	AnnounceUrl    string
	TimeToAnnounce int
	Peers          []string `gorm:"type:text[]"`
	Seeders        int
	Leechers       int

	PiecesData  []byte
}

func GetInfoHash(info_entry map[any]any) ([]byte, error) {
	hash := sha1.New()
	// Write the data to the hash object
	_, err := hash.Write(encode(info_entry))
	if err!= nil {
        return nil, err
    }
	// Get the SHA-1 hash sum
	hashSum := hash.Sum(nil)
	return hashSum, nil
}

func GetAnnounceUrl(decoded map[any]any) (string, error) {
	announceURL, ok := decoded["announce"]
	if !ok {
		return "", fmt.Errorf("error finding announcement url key")
	}


	announceURLStr, ok := announceURL.([]byte)

	if !ok {
		return "", fmt.Errorf("error decoding announce url")
	}

	return string(announceURLStr), nil
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

func GetName(info_entry map[any]any) (string, error) {
	name, ok := info_entry["name"]
	if !ok {
		return "", fmt.Errorf("error finding name key in info entry")
	}

	name_as_str, ok := name.([]byte)

	if !ok{
		return "", fmt.Errorf("error decoding name in info entry")
	}


	return string(name_as_str), nil
}

func GetPiecesData(info_entry map[any]any) ([]byte, error) {
	pieces_hashes, ok := info_entry["pieces"]
	if !ok {	
		return nil, fmt.Errorf("error decoding pieces in info entry")
	}

	pieces, ok := pieces_hashes.([]byte)
	if!ok {
        return nil, fmt.Errorf("error converting pieces in info entry to []byte")
    }

	return pieces, nil

}


func ParseTorrentFromBencoded(decoded any) (*Torrent, error) {	
	decoded_as_map, ok := decoded.(map[any]any)
	if !ok {
        return nil, fmt.Errorf("invalid bencoded torrent")
    }

	info_entry, ok := decoded_as_map["info"]
	if !ok {
        return nil, fmt.Errorf("no info entry in bencoded torrent")
    }

	info_as_map, ok := info_entry.(map[any]any)

	if !ok {
		return nil, fmt.Errorf("cant convert info entry to map[any]any")
	}

	tname, err := GetName(info_as_map)

	if err!= nil {
        return nil, err
    }

	info_hash, err := GetInfoHash(info_as_map)
	if err!= nil {
        return nil, err
    }

	size, err := GetSize(info_as_map)

	if (err!= nil) {
		return nil, err
	}

	announce_url, err := GetAnnounceUrl(decoded_as_map)
	if (err!= nil) {
		return nil, err
	}

	pieces_data, err := GetPiecesData(info_as_map)
	if (err!= nil) {
        return nil, err
    }

	return &Torrent{
        TorrentName: tname,
        Infohash:    info_hash,
        Size:        size,
        Downloaded:  0,
        Uploaded:    0,

        AnnounceUrl:    announce_url,
        TimeToAnnounce: 0,
        Peers:          []string{},
        Seeders:        0,
        Leechers:		0,

		PiecesData:  pieces_data,
	}, nil
}