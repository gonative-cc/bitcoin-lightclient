package main

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"gotest.tools/assert"
)

func TestInsertHeader(t *testing.T) {

	headers := []string{
		"0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff7f20020000000101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000",
"0000002006226e46111a0b59caaf126043eb5bbf28c34f3a5e332a1fc7b2b73cf188910f2fe76e709f3031b5ed684f098b5cd35a09633943d141a6d0525f34a1643dcf44e0b23c67ffff7f200200000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025100ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000205dcd36bceabfc1816ab17503f753e8de66be7ceddee2b8b806b85e61cf9fdc68ce1efdd1cb457e408bc2984151f5ba8e7efcf8a64c2ca9f23de07a3e718a7e9de1b23c67ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025200ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020866a367051ccd8849186173fe851392406fc1ac4b444776466260a34b58820351ab2820d880bf70bfe937525e3e7cdd48142186d7fefbef16c42f235f923aa70e1b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025300ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020849665c354e016f4ada91e1bb4d77cfe85e8c40b466641455bd3df70d1e80d363862bfea99b4b68243fcc88e8166a89ace00f5a10a8720e9603684bed07a6278e2b23c67ffff7f200200000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025400ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000203b05a05cb26319f421e41c1391ce11d94d7648396c74096560951c4bf412d7174339e31618fee9e9d0acad931c64d6a0785b66cea24a3ef630962f3bd272530be2b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025500ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000202cbb8366b0051b44757d626cb859631994adc82399ad37ce0b04edf01e829924ac2b58bde3db81c8aea571e4c90d4f638164ab8debd192c44767575bfc28edbbe2b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025600ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020f8557187addf502f0bbf3fa3a7baa9c5021166315c5b62c7c5ab63f69726ce646e4393f90540542998e3934f0181df08d07afb0daa8e0bedb23e603523739aaae2b23c67ffff7f200400000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025700ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000203890e056da2ec78018dc491b61cfbb3134c88ee0f76e838db758b7ce7a237327f97921cc0471669591969f0eb960a67d72ff58d19496a58c9bf48c477fb80242e3b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025800ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000204716da4899c77d5c64a0f17ff69c568f31d5c5429939ade626c6ef58cf972437ad21b69fc615dee7c8c940526a690bb81205b634ab0992ec465e7a2fb39df4a2e3b23c67ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025900ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000207b58138803dfe75a1a035a8519489acefde0768e1b3295fb988ea93dd8199926f6a63660b9a502ed33eacf7130b980208e73c63c623ca27cb6f4185550d0a262e3b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025a00ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"000000209e26cb5b513b048b423dbdd6fed32620feb1bf8dba1e069473d3ad83a0c3477ae4e0efa71fc24467985b923badf291bb5cd62b81bfa4aeec335b86e091f964bae3b23c67ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025b00ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020e78f4506e5f9f958b41da8b467c84377efd02bad4ff7c214f98971bab4fc936399c8a0513e6e0f2c9f804d1644a779789f0b0a9b47c0b1e26f08ae08e5ff967ae3b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025c00ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020b084bde47ad161583492f9d397d6cfc87950bef4d80b5720684fe952fe189327d60d5179f5512c82cd32a46176260dee781e0c30f499f81ce7e2fcda8a50bfa9e3b23c67ffff7f200300000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025d00ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"0000002054775e67825f7a4caf14b3fabe95af99a4204f5c3960b935aa429966e1e143498af4ff5bc198b3ab10e9d542be5e0c0e0e8d249e7d4bfa19eb2517c978f06737e4b23c67ffff7f200700000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025e00ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020f972e238a280d6a3f0eb7560deae346c29fdd0203722510e0a82f792f1020933a5f456f5967ce558ac2da01200814edd3b045f9f7f1a6f8235086260d3ea4a80e4b23c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025f00ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"00000020b27998bce45300e48a19d64f57fb219c82e7edc283da74b73eeb4125db49d51bee3a55f8afeef128522c5a0ef633840e098d5288d077eb9c1bd6a1e01bf333b1e4b23c67ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff026000ffffffff0200f2052a01000000160014f1734bfdd31f315548e1a654e686d3b6367ecec70000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
"0000002018c821aeb4f94b3847dccc7946ff82a4d022a9e87162bc0601992b7dbaf12b432237b251c5f3473a3ece242cb340c6c7d69ec9db9579d57c38069341dfabcf98bdb43c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff03011100ffffffff0200f2052a010000001600145bda3166f8ae9ff21d249197b66acf08326af84b0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
	}
	type insertHeaderTestCase struct {
		name    string
		headers []string
		header  string
		err     error
	}

	run := func(t *testing.T, tc insertHeaderTestCase) {
		decodedHeader := make([]wire.BlockHeader, len(tc.headers))
		btcHeader, _ := BlockHeaderFromHex(tc.header)
		for id, str := range tc.headers {
			h, _ := BlockHeaderFromHex(str)
			decodedHeader[id] = h
		}
		lc := NewBTCLightClientWithData(&chaincfg.RegressionNetParams, decodedHeader, 0)

		lcErr := lc.InsertHeader(btcHeader)

		if tc.err == nil {
			assert.Assert(t, lcErr == nil)
		} else if errors.Is(lcErr, tc.err) {
			t.Fatalf("Error not match")
		}

	}

	tcs := []insertHeaderTestCase{
		{
			name:    "Success Insert",
			headers: headers,
			header:  "000000201dc5f8e7cdb2dda12a307615bd0b9847c60f813baff61a591cb15c44f6a242205abc34309a00ac2983bb4cf1bc76ebc8800ec9965d423508fba603eb763b9215bdb43c67ffff7f200200000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff03011200ffffffff0200f2052a010000001600145bda3166f8ae9ff21d249197b66acf08326af84b0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
			err:     nil,
		},
		{
			name:    "Insert failed",
			headers: headers,
			header:  "0000002006226e46111a0b59caaf126043eb5bbf28c34f3a5e332a1fc7b2b73cf188910f91bf7fc009e51a44f6c7b063e64d80b36af5cb8bc9879b9dadc7eebec779a70b437f3c67ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025100ffffffff0200f2052a01000000160014db62d5ead43bde6defb99c188151bf9c9d37f6c30000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
			err:     errors.New("fork too old"),
		},
		{
			name:    "Create fork",
			headers: headers,
			header:  "0000002018c821aeb4f94b3847dccc7946ff82a4d022a9e87162bc0601992b7dbaf12b43b275a58e2107b29a5735d8c6bd63d674ef01d4596c78961f338e5a364693b03502b53c67ffff7f200000000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff03011100ffffffff0200f2052a01000000160014c6aee90d5d3120a43a090b1e3720aea464248c5f0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000",
			err:     nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
func TestCleanup(t *testing.T) {
	// test-case1
	// b1 <-b2 <- b3  .... b9
	// run cleanup
	// b2 <- b3  .... b9

	// test-case2
	// b3  ... b8 <- b9 <- b10
	//                  \- c10
	// run cleanup
	// b3  ... b8 <- b9 <- b10
	//                  \- c10

	// test-case3
	// b1 <- b2 <- b3  .... b9
	//          \- c3
	//    \- d2 <- d3 <- d4
	// run cleanup
	// b2 <- b3  .... b9
	//    \- c3

}
