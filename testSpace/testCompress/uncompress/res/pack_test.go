package compress

import "testing"

func TestPack(t *testing.T) {
	filePath := "../oldDir"
	dest := "../zip"
	Tar2(filePath, dest, false)
}
