package network

import (
	"log"
	"testing"
)

func TestUpload(t *testing.T) {
	if err := Upload("C:\\Users\\Whisper\\Pictures\\soft_engine", "test1111"); err != nil {
		log.Fatal(err)
	}
}

func TestDownload(t *testing.T) {
	if err := Download("test1111", "C:\\Users\\Whisper\\Pictures\\soft_engine"); err != nil {
		log.Fatal(err)
	}
}
