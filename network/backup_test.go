package network

import (
	"testing"
)

func TestConnect(t *testing.T) {
	if _, err := connect(user, passwd, hostAddr); err != nil {
		t.Error(err)
	}
}
