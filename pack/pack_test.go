package pack

import (
	"log"
	"testing"
)

func TestPack(t *testing.T) {
	filePath := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir"
	dest := "../testSpace/testPackDir/t"
	if err := Tar(filePath, dest); err != nil {
		log.Fatal(err)
	}
}

func TestUnPack(t *testing.T) {
	filePath := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\testSpace\\testPackDir\\1"
	dest := "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\testSpace\\testPackDir"
	if err := UnTar(filePath, dest); err != nil {
		log.Fatal(err)
	}
}
