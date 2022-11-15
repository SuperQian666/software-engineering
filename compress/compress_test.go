package compress

import "testing"

func TestCompress(t *testing.T) {
	if err := Zip("C:\\Users\\Whisper\\GolandProjects\\software-engineering\\pack", "C:\\Users\\Whisper\\GolandProjects\\software-engineering\\testSpace\\testCompress\\compress\\copy.zip"); err != nil {
		t.Fatal(err)
	}
}

func TestUnCompress(t *testing.T) {
	if err := UnZip("", ""); err != nil {
		t.Fatal(err)
	}
}
