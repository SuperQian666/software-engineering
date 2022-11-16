package pack

import (
	"log"
	"testing"
)

func TestPack(t *testing.T) {
	filePath := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\1"
	dest := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\2"
	if err := Tar(filePath, dest); err != nil {
		log.Fatal(err)
	}
}

func TestUnPack(t *testing.T) {
	filePath := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\2"
	dest := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\tt"
	if err := UnTar(filePath, dest); err != nil {
		log.Fatal(err)
	}
}
