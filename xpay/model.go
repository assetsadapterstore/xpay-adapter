/*
* Copyright 2018 The OpenWallet Authors
* This file is part of the OpenWallet library.
*
* The OpenWallet library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The OpenWallet library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
 */

package xpay

import (
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/tidwall/gjson"
	"time"
)

const TimeLayout = `2006-01-02T15:04:05Z07:00`

type XIFAccount struct {
	Address   string `json:"address"`
	Publickey string `json:"publickey"`
	Symbol    string `json:"symbol"`
	Amount    string `json:"amount"`
	Nonce     uint64 `json:"nonce"`
	AccType   string `json:"type"`
}

func NewXIFAccount(result *gjson.Result) *XIFAccount {
	obj := XIFAccount{}
	obj.Address = result.Get("account.address").String()
	obj.Publickey = result.Get("account.publickey").String()
	obj.Symbol = result.Get("account.symbol").String()
	obj.Amount = result.Get("account.amount").String()
	obj.Nonce = result.Get("account.nonce").Uint()
	obj.AccType = result.Get("account.type").String()
	return &obj
}

type Block struct {
	Height    uint64    `json:"id"`
	Hash      string    `json:"hash"`
	LastHash  string    `json:"last_hash"`
	Txns      []string  `json:"txns"`
	Timestamp time.Time `json:"timestamp"`
}

func NewBlock(result *gjson.Result) *Block {
	obj := Block{}
	obj.Height = result.Get("id").Uint()
	obj.Hash = result.Get("hash").String()
	obj.LastHash = result.Get("last_hash").String()
	obj.Timestamp, _ = time.ParseInLocation(TimeLayout, result.Get("created").String(), time.UTC)
	obj.Txns = make([]string, 0)
	if txns := result.Get("txns"); txns.IsArray() {
		for _, tx := range txns.Array() {
			obj.Txns = append(obj.Txns, tx.String())
		}
	}
	return &obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Previousblockhash = b.LastHash
	obj.Height = b.Height
	obj.Time = uint64(b.Timestamp.Unix())
	obj.Symbol = symbol

	return &obj
}

type Transaction struct {
	/*
		{
		    "transaction": {
		        "key": "bfebbbbad1d2da76cfbcec246c8b8703c04cbcdf23144b02db033430319801c8",
		        "owner": "024ea4abd05b5d3a6ab3a877e1943c94ba71487dee1eea1a96b952b3ff99486bc7",
		        "created": "2020-03-11T03:05:13.200Z",
		        "sender_account": "02a826406c8b5e0b55f3484c6bd3fa40749a4fd5cc1739f1b547b798fbb4023e61",
		        "recipient_account": "024ea4abd05b5d3a6ab3a877e1943c94ba71487dee1eea1a96b952b3ff99486bc7",
		        "amount": "0.05",
		        "amount_num": "0.05",
		        "symbol": "XIF",
		        "type": "TRANSFER",
		        "hash": "86ea2aeb3dc84a16591c1ff68720822d5d14c2fb90484244a77f66c2d0eb01cf",
		        "block": "947684",
		        "signature": "304402202cc5dffe3de1b9c23cbee480eb7c2fb35e5fac577ff23ddf9d843d21e765ec41022010807232325ecd55eb5cad9e533a2acbc9ca594cd409fe1d89f8f89ca28d2947",
		        "notes": "Fee for f50d1b3e9ffdca74a8d6953627a68f3d852e37ecbe2e1f66469038657572b5ef",
		        "status": "COMPLETED"
		    }
		}
	*/

	Hash        string
	Owner       string
	From        string
	To          string
	Amount      string
	Symbol      string
	BlockHash   string
	BlockHeight uint64
	Status      string
	TxType      string
	Memo        string
	Timestamp   time.Time
}

func NewTransaction(result *gjson.Result) *Transaction {
	obj := Transaction{}
	obj.Hash = result.Get("transaction.key").String()
	obj.Owner = result.Get("transaction.owner").String()
	obj.From = result.Get("transaction.sender_account").String()
	obj.To = result.Get("transaction.recipient_account").String()
	obj.Amount = result.Get("transaction.amount").String()
	obj.Symbol = result.Get("transaction.symbol").String()
	obj.BlockHash = result.Get("transaction.hash").String()
	obj.BlockHeight = result.Get("transaction.block").Uint()
	obj.Status = result.Get("transaction.status").String()
	obj.TxType = result.Get("transaction.type").String()
	obj.Memo = result.Get("transaction.notes").String()
	obj.Timestamp, _ = time.ParseInLocation(TimeLayout, result.Get("transaction.created").String(), time.UTC)
	return &obj
}

type RawTransaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Symbol    string `json:"symbol"`
	Amount    string `json:"amount"`
	Nonce     uint64 `json:"nonce"`
	Signature string `json:"signature"`
}

func (rawTx *RawTransaction) Hash() []byte {
	message := fmt.Sprintf("%s%s%s%s%d", rawTx.Sender, rawTx.Recipient, rawTx.Symbol, rawTx.Amount, rawTx.Nonce)
	messageHash := owcrypt.Hash([]byte(message), 0, owcrypt.HASH_ALG_SHA256)
	return messageHash
}

func (rawTx *RawTransaction) FillSig(signature []byte) error {
	if len(signature) != 64 {
		return fmt.Errorf("signature length is not equal 64 bytes")
	}
	//DER-encoded, 30440220+前32字节+0220+后32字节
	lBytes := signature[:32]
	rBytes := signature[32:]
	der := append([]byte{0x30, 0x44, 0x02, 0x20}, lBytes...)
	der = append(der, 0x02, 0x20)
	der = append(der, rBytes...)
	rawTx.Signature = hex.EncodeToString(der)
	return nil
}
