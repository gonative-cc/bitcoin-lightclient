package btclightclient

import (
	"encoding/hex"
	"testing"

	"gotest.tools/assert"
)

func TestBTCHeaderFromHex(t *testing.T) {
	type testCase struct {
		name              string
		headerHex         string
		expectedBlockHash string
		expectedError     error
	}

	run := func(t *testing.T, tc testCase) {
		header, err := BlockHeaderFromHex(tc.headerHex)
		blockHash := header.BlockHash()
		if tc.expectedError != nil {
			assert.Error(t, err, tc.expectedError.Error())
		} else {
			assert.NilError(t, err)
			assert.Equal(t, tc.expectedBlockHash, blockHash.String())
		}
	}

	testCases := []testCase{
		{
			name:              "happy test case",
			headerHex:         "020000004cdba1415b2c6e7808c1b3c18df1374238454f7104203475bf01000000000000c17ea9d06015dc83902911cd24837a8ba4bdc0c1d72b873f906d921e06e48d2f984a8250ef75051a72a8061a",
			expectedError:     nil,
			expectedBlockHash: "0000000000000363d7f5f3341fb0b5b69949103e2d681591c9f737e4ea67e2a7",
		},
		{
			name:              "invalid header length",
			headerHex:         "020000004cdba1415b2c6e7808c1b3c18df1374238454f7104203475bf01000000000000c17ea9d06015dc83902911cd24837a8ba4bdc0c1d72b873f906d921e06e48d2f984a8250ef75051a72a8061ab",
			expectedError:     ErrInvalidHeaderSize,
			expectedBlockHash: "0000000000000363d7f5f3341fb0b5b69949103e2d681591c9f737e4ea67e2a7",
		},
		{
			name:              "non-hex character",
			headerHex:         "020000004cdba1415b2c6e7808c1b3c18df1374238454f7104203475bf01000000000000c17ea9d06015dc83902911cd24837a8ba4bdc0c1d72b873f906d921e06e48d2f984a8250ef75051a72a8061Q",
			expectedError:     hex.InvalidByteError('Q'),
			expectedBlockHash: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
