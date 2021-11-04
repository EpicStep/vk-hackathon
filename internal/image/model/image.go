package model

type Image struct {
	ID string
	Image []byte
	Hash uint64
	Height, Width int
}
