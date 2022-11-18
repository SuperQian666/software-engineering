package compress

import "testing"

func TestCompress(t *testing.T) {
	if err := Zip("C:\\Users\\Whisper\\Pictures\\test\\2.jpg", "C:\\Users\\Whisper\\Pictures\\test\\1119"); err != nil {
		t.Fatal(err)
	}
}

func TestUnCompress(t *testing.T) {
	if err := UnZip("C:\\Users\\Whisper\\Pictures\\test\\1119", "C:\\Users\\Whisper\\Pictures\\test\\1119new"); err != nil {
		t.Fatal(err)
	}
}
