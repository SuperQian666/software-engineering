package copy

import (
	"testing"
)

func TestFileCopy(t *testing.T) {
	if err := FileCopy("C:/Users/Whisper/GolandProjects/software-engineering/src.txt", "../dec.txt"); err != nil {
		t.Error(err)
	}
}
