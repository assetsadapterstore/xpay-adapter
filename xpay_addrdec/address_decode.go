package xpay_addrdec

import (
	"encoding/hex"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/openwallet"
)

var (
	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	*openwallet.AddressDecoderV2Base
}

//NewAddressDecoder 地址解析器
func NewAddressDecoderV2() *AddressDecoderV2 {
	decoder := AddressDecoderV2{}
	return &decoder
}

//AddressDecode 地址解析
func (dec *AddressDecoderV2) AddressDecode(addr string, opts ...interface{}) ([]byte, error) {

	return hex.DecodeString(addr)
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {
	if len(hash) >= 64 {
		//压缩公钥
		hash = owcrypt.PointCompress(hash, owcrypt.ECC_CURVE_NIST_P256)
	}
	return hex.EncodeToString(hash), nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {
	pub, err := hex.DecodeString(address)
	if err != nil {
		return false
	}

	if len(pub) != 33 {
		return false
	}
	return true
}
