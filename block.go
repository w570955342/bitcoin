package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

//1. 定义区块结构
type Block struct {
	//1. 版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. 梅克尔根
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数
	Nonce uint64

	//a. 当前区块哈希,正常比特币区块中没有当前区块的哈希！
	Hash []byte
	//b. 交易数据
	Data []byte
}

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

//2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{}, //先填空，后面再计算 //TODO
		Data:       []byte(data),
	}

	block.SetHash()
	return &block
}

//3. 为Block绑定方法生成哈希
func (block *Block) SetHash() {
	var blockInfo []byte
	//1. 拼接数据
	/*blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, block.PrevHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	//blockInfo = append(block.PrevHash, block.Data...)//错误
	blockInfo = append(blockInfo, block.Data...)*/

	//简化一下
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}

	blockInfo=bytes.Join(tmp,[]byte{})

	//2. sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}