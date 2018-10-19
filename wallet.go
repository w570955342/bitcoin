package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

//定义秘钥结构体，每个结构体包含一对秘钥（公钥和私钥）
type Key struct {
	//私钥
	Private *ecdsa.PrivateKey
	//PubKey *ecdsa.PublicKey
	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	PubKey []byte
}

//创建钱包
func NewKey() *Key {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}

	//生成公钥
	pubKeyOrig := privateKey.PublicKey

	//拼接X, Y
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	return &Key{Private: privateKey, PubKey: pubKey}
}

//生成地址
