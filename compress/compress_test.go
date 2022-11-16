package compress

import "testing"

func TestCompress(t *testing.T) {
	if err := Zip("C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir", "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\2"); err != nil {
		t.Fatal(err)
	}
}

func TestUnCompress(t *testing.T) {
	if err := UnZip("C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\2", "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\oldDir\\3"); err != nil {
		t.Fatal(err)
	}
}
