package consts

type MessageType uint8

const (
	TEXT MessageType = iota
	IMAGE
	VIDEO
	AUDIO
	DOCUMENT
	SYSTEM
)
