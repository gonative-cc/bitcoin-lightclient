package btclightclient


type UTXOProof struct {
	
}

func (lc *BTC) verifyUTXOProof(utxoProof UTXOProof) bool {
	return true
}

func (lc *BTCLightClient) VerifyUTXOProofs(utxoProofs []UTXOProof) bool {
	return true
}
