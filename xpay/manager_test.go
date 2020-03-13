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
	"github.com/astaxie/beego/config"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/hdkeystore"
	"github.com/blocktree/openwallet/log"
	"path/filepath"
	"testing"
)

var (
	tw *WalletManager
)

func init() {
	tw = testNewWalletManager()
}

func testNewWalletManager() *WalletManager {
	wm := NewWalletManager()

	//读取配置
	absFile := filepath.Join("conf", "XIF.ini")
	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return nil
	}
	wm.LoadAssetsConfig(c)
	wm.client.Debug = true

	return wm
}

func TestWalletManager_GetWalletDetails(t *testing.T) {
	result, err := tw.GetWalletDetails("027a4522cfa1c35b51aa1a0f881ab1ccb19deae6277917b7d641490d4492083711")
	if err != nil {
		t.Errorf("GetWalletDetails failed, err: %v", err)
	}
	log.Infof("result: %+v", result)
}

func TestWalletManager_NewWallet(t *testing.T) {
	result, err := tw.NewWallet("AUSD")
	if err != nil {
		t.Errorf("NewWallet failed, err: %v", err)
	}
	log.Infof("result: %+v", result)
}

func TestWalletManager_CreateLocalWallet(t *testing.T) {
	prv, _ := hdkeystore.GenerateSeed(32)
	pub, _ := owcrypt.GenPubkey(prv, owcrypt.ECC_CURVE_NIST_P256)
	log.Infof("prv: %s", hex.EncodeToString(prv))
	comPub := owcrypt.PointCompress(pub, owcrypt.ECC_CURVE_NIST_P256)
	log.Infof("pub: %s", hex.EncodeToString(comPub))
}

func TestWalletManager_GetLatestBlock(t *testing.T) {
	result, err := tw.GetLatestBlock()
	if err != nil {
		t.Errorf("GetLatestBlock failed, err: %v", err)
	}
	log.Infof("result: %+v", result)
	//947694
}

func TestWalletManager_GetBlock(t *testing.T) {
	result, err := tw.GetBlock(947692)
	if err != nil {
		t.Errorf("GetBlock failed, err: %v", err)
	}
	log.Infof("result: %+v", result)
}

func TestWalletManager_GetTransaction(t *testing.T) {
	result, err := tw.GetTransaction("f022671089c7bbb169705822096709818946369ce1a11d8858852e0c67cb9bb2")
	if err != nil {
		t.Errorf("GetTransaction failed, err: %v", err)
	}
	log.Infof("result: %+v", result)
}

func TestWalletManager_SignRawTxOnline(t *testing.T) {
	sender := "027a4522cfa1c35b51aa1a0f881ab1ccb19deae6277917b7d641490d4492083711"
	privateKey, _ := hex.DecodeString("672c6012ef49a30d8b9b7501706ef3769aaefc38d72d0f048dfade7a850100b4")
	w, err := tw.GetWalletDetails(sender)
	if err != nil {
		t.Errorf("GetWalletDetails failed, err: %v", err)
	}
	nonce := w.Nonce + 3
	rawTx := &RawTransaction{
		Sender:    sender,
		Recipient: "033e379d467f0cb36b30b068f5fd9c81bd4ae7d2dbb93a5e08bad7cf2671eb6f46",
		Symbol:    "XIF",
		Amount:    "0.01",
		Nonce:     nonce,
	}
	err = tw.SignRawTxOnline(rawTx, privateKey)
	if err != nil {
		t.Errorf("SignRawTxOnline failed, err: %v", err)
	}
	log.Infof("signature: %s", rawTx.Signature)
}

func TestWalletManager_SignRawTxOffline(t *testing.T) {
	sender := "027a4522cfa1c35b51aa1a0f881ab1ccb19deae6277917b7d641490d4492083711"
	privateKey, _ := hex.DecodeString("672c6012ef49a30d8b9b7501706ef3769aaefc38d72d0f048dfade7a850100b4")
	w, err := tw.GetWalletDetails(sender)
	if err != nil {
		t.Errorf("GetWalletDetails failed, err: %v", err)
	}
	nonce := w.Nonce + 1
	rawTx := &RawTransaction{
		Sender:    sender,
		Recipient: "033e379d467f0cb36b30b068f5fd9c81bd4ae7d2dbb93a5e08bad7cf2671eb6f46",
		Symbol:    "XIF",
		Amount:    "0.01",
		Nonce:     nonce,
	}
	err = tw.SignRawTxOffline(rawTx, privateKey)
	if err != nil {
		t.Errorf("SignRawTxOffline failed, err: %v", err)
	}
	log.Infof("signature: %s", rawTx.Signature)
}

func TestWalletManager_SendRaw(t *testing.T) {
	sender := "027a4522cfa1c35b51aa1a0f881ab1ccb19deae6277917b7d641490d4492083711"
	privateKey, _ := hex.DecodeString("672c6012ef49a30d8b9b7501706ef3769aaefc38d72d0f048dfade7a850100b4")
	w, err := tw.GetWalletDetails(sender)
	if err != nil {
		t.Errorf("GetWalletDetails failed, err: %v", err)
		return
	}
	nonce := w.Nonce + 1
	rawTx := &RawTransaction{
		Sender:    sender,
		Recipient: "0336b1d33588a4e71ae83386174b52201c8f6a1cdd3c77af852794a9e5f55530c7",
		Symbol:    "XIF",
		Amount:    "3.7",
		Nonce:     nonce,
	}
	//03a0f4450a7ebc06354caa601601290e90ad950b775df3bee621a763a358dfe2a7
	err = tw.SignRawTxOffline(rawTx, privateKey)
	if err != nil {
		t.Errorf("SignRawTxOffline failed, err: %v", err)
		return
	}
	result, err := tw.Sendraw(rawTx)
	if err != nil {
		t.Errorf("SignRawTxOffline failed, err: %v", err)
		return
	}

	log.Infof("result: %+v", result)
}
