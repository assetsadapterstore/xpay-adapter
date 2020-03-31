package xpay

import (
	"encoding/hex"
	"github.com/blocktree/go-owcrypt"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {

	p2pk, _ := hex.DecodeString("cf036690a6fbbcdd9dfdd6249e1d121bec1eacd8")
	p2pkAddr, _ := tw.Decoder.AddressEncode(p2pk)
	t.Logf("p2pkAddr: %s", p2pkAddr)
}

func TestDecompressPubKey(t *testing.T) {
	pub, _ := hex.DecodeString("03a2147994c34ec6ac3ba4d0737e672002a832f2cf050f0d2cffc7906ad8a1e7b2")
	uncompessedPublicKey := owcrypt.PointDecompress(pub, owcrypt.ECC_CURVE_SECP256K1)
	t.Logf("pub: %s", hex.EncodeToString(uncompessedPublicKey))
}