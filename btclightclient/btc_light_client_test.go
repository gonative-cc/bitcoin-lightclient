package btclightclient

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"gotest.tools/assert"
)

type commonTestCaseData struct {
	headers []string
	header  string
}

type testCaseMap map[string]commonTestCaseData

var HEADERS = []string{
	"0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff7f2002000000",
	"0000002006226e46111a0b59caaf126043eb5bbf28c34f3a5e332a1fc7b2b73cf188910f2fe76e709f3031b5ed684f098b5cd35a09633943d141a6d0525f34a1643dcf44e0b23c67ffff7f2002000000",
	"000000205dcd36bceabfc1816ab17503f753e8de66be7ceddee2b8b806b85e61cf9fdc68ce1efdd1cb457e408bc2984151f5ba8e7efcf8a64c2ca9f23de07a3e718a7e9de1b23c67ffff7f2001000000",
	"00000020866a367051ccd8849186173fe851392406fc1ac4b444776466260a34b58820351ab2820d880bf70bfe937525e3e7cdd48142186d7fefbef16c42f235f923aa70e1b23c67ffff7f2000000000",
	"00000020849665c354e016f4ada91e1bb4d77cfe85e8c40b466641455bd3df70d1e80d363862bfea99b4b68243fcc88e8166a89ace00f5a10a8720e9603684bed07a6278e2b23c67ffff7f2002000000",
	"000000203b05a05cb26319f421e41c1391ce11d94d7648396c74096560951c4bf412d7174339e31618fee9e9d0acad931c64d6a0785b66cea24a3ef630962f3bd272530be2b23c67ffff7f2000000000",
	"000000202cbb8366b0051b44757d626cb859631994adc82399ad37ce0b04edf01e829924ac2b58bde3db81c8aea571e4c90d4f638164ab8debd192c44767575bfc28edbbe2b23c67ffff7f2000000000",
	"00000020f8557187addf502f0bbf3fa3a7baa9c5021166315c5b62c7c5ab63f69726ce646e4393f90540542998e3934f0181df08d07afb0daa8e0bedb23e603523739aaae2b23c67ffff7f2004000000",
	"000000203890e056da2ec78018dc491b61cfbb3134c88ee0f76e838db758b7ce7a237327f97921cc0471669591969f0eb960a67d72ff58d19496a58c9bf48c477fb80242e3b23c67ffff7f2000000000",
	"000000204716da4899c77d5c64a0f17ff69c568f31d5c5429939ade626c6ef58cf972437ad21b69fc615dee7c8c940526a690bb81205b634ab0992ec465e7a2fb39df4a2e3b23c67ffff7f2001000000",
	"000000207b58138803dfe75a1a035a8519489acefde0768e1b3295fb988ea93dd8199926f6a63660b9a502ed33eacf7130b980208e73c63c623ca27cb6f4185550d0a262e3b23c67ffff7f2000000000",
	"000000209e26cb5b513b048b423dbdd6fed32620feb1bf8dba1e069473d3ad83a0c3477ae4e0efa71fc24467985b923badf291bb5cd62b81bfa4aeec335b86e091f964bae3b23c67ffff7f2001000000",
	"00000020e78f4506e5f9f958b41da8b467c84377efd02bad4ff7c214f98971bab4fc936399c8a0513e6e0f2c9f804d1644a779789f0b0a9b47c0b1e26f08ae08e5ff967ae3b23c67ffff7f2000000000",
	"00000020b084bde47ad161583492f9d397d6cfc87950bef4d80b5720684fe952fe189327d60d5179f5512c82cd32a46176260dee781e0c30f499f81ce7e2fcda8a50bfa9e3b23c67ffff7f2003000000",
	"0000002054775e67825f7a4caf14b3fabe95af99a4204f5c3960b935aa429966e1e143498af4ff5bc198b3ab10e9d542be5e0c0e0e8d249e7d4bfa19eb2517c978f06737e4b23c67ffff7f2007000000",
	"00000020f972e238a280d6a3f0eb7560deae346c29fdd0203722510e0a82f792f1020933a5f456f5967ce558ac2da01200814edd3b045f9f7f1a6f8235086260d3ea4a80e4b23c67ffff7f2000000000",
	"00000020b27998bce45300e48a19d64f57fb219c82e7edc283da74b73eeb4125db49d51bee3a55f8afeef128522c5a0ef633840e098d5288d077eb9c1bd6a1e01bf333b1e4b23c67ffff7f2001000000",
	"0000002018c821aeb4f94b3847dccc7946ff82a4d022a9e87162bc0601992b7dbaf12b432237b251c5f3473a3ece242cb340c6c7d69ec9db9579d57c38069341dfabcf98bdb43c67ffff7f2000000000",
}

func CommonTestCases() testCaseMap {
	tcs := map[string]commonTestCaseData{
		"Append a fork": {
			headers: HEADERS,
			header:  "000000201dc5f8e7cdb2dda12a307615bd0b9847c60f813baff61a591cb15c44f6a242205abc34309a00ac2983bb4cf1bc76ebc8800ec9965d423508fba603eb763b9215bdb43c67ffff7f2002000000",
		},
		"Create fork": {
			headers: HEADERS,
			header:  "0000002018c821aeb4f94b3847dccc7946ff82a4d022a9e87162bc0601992b7dbaf12b43b275a58e2107b29a5735d8c6bd63d674ef01d4596c78961f338e5a364693b03502b53c67ffff7f2000000000",
		},
		"Insert failed because fork too old": {
			headers: HEADERS,
			header:  "0000002006226e46111a0b59caaf126043eb5bbf28c34f3a5e332a1fc7b2b73cf188910f91bf7fc009e51a44f6c7b063e64d80b36af5cb8bc9879b9dadc7eebec779a70b437f3c67ffff7f2001000000",
		},
		"Unknown parent": {
			headers: HEADERS,
			header:  "0000002018c8213eb4f94b3847dccc7946ff82a4d022a9e87162bc0601992b7dbaf12b43b275a58e2107b29a5735d8c6bd63d674ef01d4596c78961f338e5a364693b03502b53c67ffff7f2000000000",
		},
	}

	return tcs
}

func initLightClient(t *testing.T, headers []string) *BTCLightClient {
	decodedHeaders := make([]wire.BlockHeader, len(headers))
	for id, str := range headers {
		h, err := BlockHeaderFromHex(str)
		assert.NilError(t, err)
		decodedHeaders[id] = h
	}
	lc := NewBTCLightClientWithData(&chaincfg.RegressionNetParams, decodedHeaders, 0)
	return lc
}

func TestInsertHeader(t *testing.T) {
	commonTestCase := CommonTestCases()
	testCases := map[string]error{
		"Append a fork":                      nil,
		"Create fork":                        nil,
		"Insert failed because fork too old": ErrForkTooOld,
		"Unknown parent":                     ErrParentBlockNotInChain,
	}

	run := func(t *testing.T, testcase string) {
		data := commonTestCase[testcase]

		lc := initLightClient(t, data.headers)
		btcHeader, _ := BlockHeaderFromHex(data.header)
		lcErr := lc.InsertHeader(btcHeader)

		expectedErr := testCases[testcase]
		if expectedErr == nil {
			assert.NilError(t, lcErr)
		} else {
			assert.Error(t, lcErr, expectedErr.Error())
		}

	}

	for testcase := range testCases {
		t.Run(testcase, func(t *testing.T) {
			run(t, testcase)
		})
	}
}

func TestLatestFinalizedBlock(t *testing.T) {
	commonTestCase := CommonTestCases()
	type TestCase struct {
		Error     error
		Height    int64
		BlockHash string
	}
	testCases := map[string]TestCase{
		// add new block to the longest chain, increasing finalised block height by 1
		"Append a fork": {nil, 11, "6393fcb4ba7189f914c2f74fad2bd0ef7743c867b4a81db458f9f9e506458fe7"},
		// add new block to a fork, keeping finalised block height same
		"Create fork": {nil, 10, "7a47c3a083add37394061eba8dbfb1fe2026d3fed6bd3d428b043b515bcb269e"},
		// fails new block addition, keeping finalised block height same
		"Insert failed because fork too old": {ErrForkTooOld, 10, "7a47c3a083add37394061eba8dbfb1fe2026d3fed6bd3d428b043b515bcb269e"},
		// no new block in forks, keeping finalised block height same
		"Unknown parent": {ErrParentBlockNotInChain, 10, "7a47c3a083add37394061eba8dbfb1fe2026d3fed6bd3d428b043b515bcb269e"},
	}

	run := func(t *testing.T, name string, tc TestCase) {
		data := commonTestCase[name]

		lc := initLightClient(t, data.headers)
		btcHeader, err := BlockHeaderFromHex(data.header)
		assert.NilError(t, err)

		lcErr := lc.InsertHeader(btcHeader)
		err = lc.CleanUpFork()
		assert.NilError(t, err)
		if tc.Error == nil {
			assert.NilError(t, lcErr)
		} else {
			assert.ErrorType(t, lcErr, tc.Error)
		}

		assert.Equal(t, lc.LatestFinalizedBlockHeight(), tc.Height)
		assert.Equal(t, lc.LatestFinalizedBlockHash().String(), tc.BlockHash)
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			run(t, name, tc)
		})
	}
}

func TestCleanup(t *testing.T) {
	var err error

	tcs := CommonTestCases()
	lc := initLightClient(t, HEADERS)

	// test-case1
	// b1 <-b2 <- b3  .... b9
	// run cleanup
	// b2 <- b3  .... b9
	data := tcs["Append a fork"]
	btcHeader, _ := BlockHeaderFromHex(data.header)
	beforeCheckpoint := lc.btcStore.LatestCheckPoint()
	lcErr := lc.InsertHeader(btcHeader)
	if lcErr != nil {
		t.Fatal("Should not return error when insert header")
	}
	err = lc.CleanUpFork()
	assert.NilError(t, err)
	afterCheckpoint := lc.btcStore.LatestCheckPoint()
	assert.Assert(t, beforeCheckpoint.Height+1 == afterCheckpoint.Height)

	// notthing change
	// test-case2
	// b3  ... b8 <- b9 <- b10
	//                  \- c10
	// run cleanup
	// b3  ... b8 <- b9 <- b10
	//                  \- c10

	err = lc.CleanUpFork()
	assert.NilError(t, err)
	notupdateCheckpoint := lc.btcStore.LatestCheckPoint()
	assert.Assert(t, notupdateCheckpoint == afterCheckpoint)
	data = tcs["Create fork"]
	btcHeader, _ = BlockHeaderFromHex(data.header)
	err = lc.InsertHeader(btcHeader)
	assert.NilError(t, err)

	// test-case3
	// b1 <- b2 <- b3  .... b9
	//          \- c3
	//    \- d2 <- d3 <- d4
	// run cleanup
	// b2 <- b3  .... b9
	//    \- c3
	headers := []string{
		"000000201bb3e1c443436f66b4cd58bad75748ceddd9d5737cf8c28bfffe3be786e21f6df83b5f5af28fbbac9795e543e2ef0f97bf9137305483e6c2fe59771a59303d89e1bd3d67ffff7f2002000000",
		"0000002029539cf7f719b9ab72fccace6bdf9429fa9e5b3a34338f674e82c697822ee0724dbbfdc4c44cccd0a450db68a4223e47c8ecd859c00759e98b7fd227d6ad7e21e1bd3d67ffff7f2000000000",
		"0000002008f3a2e5297fc16ac46bb17ff9f39d1f631836f9d83de21ed6cbecdf26334f15341daa558da062aa4af3020410d94df326a7e674a85b4201a8524e9098a5bdfbe1bd3d67ffff7f2000000000",
		"000000202d71344e38568cc35604fec0421107143482fafc945b3291eaab4758bcb4db11c0d28e3d5b868addec8bd15b6feaa883b867f672062f5099b720f6390565df3de1bd3d67ffff7f2005000000",
		"00000020f004fcfd82e23c77592bb920a817a9fc3c93854203b906ca728d240f986e6d4d9266200adb7bdef93d07e9c050a7b3e2d33f0d8324389999b21405d3e2ec77b7e1bd3d67ffff7f2000000000",
		"00000020f07faecc9a8a61cee690a7756f8112df5f5bfb94e5c64f39a3e61e4f8b8db35bbd276e146642e535fccf8bee242895bbdfefe92f2efcb5069e339366bf2c59e5e1bd3d67ffff7f200a000000",
		"000000203325626e050e9f17884eb04c08f89d4c689879e86bae5518ba0e77535de63976e9413b0e69182999433ea49fc52e8bb29c8d3ff7ff69c58da5f3d7978e23ea5fe2bd3d67ffff7f2000000000",
	}
	for _, h := range headers {
		btcHeader, _ := BlockHeaderFromHex(h)
		lcErr := lc.InsertHeader(btcHeader)

		if lcErr != nil {
			t.Fatal("Should not return error when insert header")
		}
		err := lc.CleanUpFork()
		assert.NilError(t, err)
	}
	listFork := lc.btcStore.LatestBlockHashOfFork()
	assert.Assert(t, len(listFork) == 1)
}
