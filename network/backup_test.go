package network

import (
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	if err := Upload("C:\\Users\\Whisper\\GolandProjects\\software-engineering\\go.mod", "12345"); err != nil {
		log.Fatal(err)
	}
}

func TestDownload(t *testing.T) {
	if err := Download("系统设计", "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\testSpace\\testbackUp"); err != nil {
		log.Fatal(err)
	}
}
