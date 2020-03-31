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
	"github.com/blocktree/openwallet/v2/hdkeystore"
	"github.com/blocktree/openwallet/v2/log"
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
	result, err := tw.GetWalletDetails("02b7b468f6e653c798b151a7f7dee454b10b540b6e518616ed4ca40f2ac8262223")
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
	result, err := tw.GetTransaction("c8cceba27d8ad7cc2a52cff30229ec5618d4ea3bc7b1b5432598ce5f168ffb68")
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
	sender := "036be8d5d120331c73c6e67990c1ce49e2240a892087bd77a75eb3692a110eeac8"
	privateKey, _ := hex.DecodeString("b7446244ddf3cb0bac88ea35751cca1a179358fa8e6202ecfa20187626da398a")
	w, err := tw.GetWalletDetails(sender)
	if err != nil {
		t.Errorf("GetWalletDetails failed, err: %v", err)
		return
	}
	nonce := w.Nonce + 1
	rawTx := &RawTransaction{
		Sender:    sender,
		Recipient: "03a94747ce9fb236b90298cd7bdb7fe71b9183032e9706d548ef64059ca1714c29",
		Symbol:    "XIF",
		Amount:    "4.75",
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

func TestWalletManager_InformWallet(t *testing.T) {
	prv, _ := hdkeystore.GenerateSeed(32)
	pub, _ := owcrypt.GenPubkey(prv, owcrypt.ECC_CURVE_NIST_P256)
	log.Infof("prv: %s", hex.EncodeToString(prv))
	comPub := owcrypt.PointCompress(pub, owcrypt.ECC_CURVE_NIST_P256)
	log.Infof("pub: %s", hex.EncodeToString(comPub))
	err := tw.InformWallet(hex.EncodeToString(comPub), "XIF")
	if err != nil {
		t.Errorf("InformWallet failed, err: %v", err)
		return
	}
}
