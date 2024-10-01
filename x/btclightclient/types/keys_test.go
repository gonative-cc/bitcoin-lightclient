package types

import (
	"bytes"
	"fmt"
	"testing"
)


func TestHeaderKey(t *testing.T) {
	headerKeyBytes, err := HeaderKey(0)

	if err != nil {
		t.Fatal("Shouldn't return error when transform header key")
	} else {
		if !bytes.Equal(headerKeyBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0}) {
			t.Fatal("Wrong header key transform")
		}
	}
}
