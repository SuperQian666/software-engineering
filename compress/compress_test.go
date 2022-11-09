package compress

import "testing"

func TestCompress(t *testing.T) {
	filePath := "../oldDir"
	dest := "../zip"
	Tar2(filePath, dest, false)
}
