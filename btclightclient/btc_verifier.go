package btclightclient

type UTXOProof struct {
	
}

func (lc *BTCLightClient) verifyUTXOProof(utxoProof UTXOProof) bool {

	return true
}

func (lc *BTCLightClient) VerifyUTXOProofs(utxoProofs []UTXOProof) (uint64, error) {
	for i := 0; i < len(utxoProofs); i++ {
		
	}
	return 0, nil
}
