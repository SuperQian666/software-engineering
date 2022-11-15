package network

import (
	"testing"
)

func TestConnect(t *testing.T) {
	Upload("C:\\Users\\Whisper\\GolandProjects\\software-engineering\\go.mod", "/var/file/")
}

func TestDownload(t *testing.T) {
	Download("C:\\Users\\Whisper\\GolandProjects\\software-engineering\\testSpace\\", "/var/file/")
}
