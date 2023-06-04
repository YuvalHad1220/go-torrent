package main

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
