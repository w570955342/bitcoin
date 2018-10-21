package main

import (
	"encoding/gob"
	"io/ioutil"
	"bytes"
	"log"
	"crypto/elliptic"
)

//定一个 Wallet结构，它保存所有的Key以及它的地址
type Wallet struct {
	//map[地址]秘钥
	WalletMap map[string]*Key
}

//生成钱包
func NewWallet() *Wallet{

	var wallet Wallet
	//wallet.WalletsMap = make(map[string]*Key)
	wallet.loadFile()
	return &wallet
}

func (wallet *Wallet)CreateWallet()string  {
	key:=NewKey()
	address:=key.NewAddress()

	wallet.WalletMap[address]=key
	wallet.saveToFile()
	return address
}
//保存方法，把新建的key添加进去
func (wallet *Wallet) saveToFile() {

	var buffer bytes.Buffer
	//panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err:=encoder.Encode(wallet)
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile("wallet.dat", buffer.Bytes(), 0600)
}

//读取文件方法，把所有的key读出来,存在内存中
func (wallet *Wallet)loadFile()  {
	//读取内容
	content, err := ioutil.ReadFile("wallet.dat")
	if err != nil {
		log.Panic(err)
	}

	//解码
	//panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))

	var walletLocal Wallet

	err = decoder.Decode(&walletLocal)
	if err != nil {
		log.Panic(err)
	}

	//wallet = &walletLocal
	//对于结构来说，里面有map的，要指定赋值，不要再最外层直接赋值
	wallet.WalletMap = walletLocal.WalletMap
}
