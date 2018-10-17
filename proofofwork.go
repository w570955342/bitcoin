package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//1. 定义工作量证明结构体 ProofOfWork
type ProofOfWork struct {
	block  *Block   //当前区块
	target *big.Int //目标值，判断当前区块是否满足要求
}

//2. 创建 ProofOfWork 结构体对象的方法
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//指定难度值
	targetStr := "0000f00000000000000000000000000000000000000000000000000000000000"
	//转化类型，见demo/big.Int.go文件
	target := big.Int{}
	target.SetString(targetStr, 16) //16指的是 指定targetStr是16进制的字符
	pow.target = &target
	return &pow
}

//3. 根据随机数和交易信息找到满足要求的哈希值
func (pow *ProofOfWork) Run() ([]byte, uint64) {

	var nonce uint64
	var hash [32]byte
	block := pow.block
	for {
		//1. 拼接数据
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			//block.Data,
			//只对区块头做哈希，区块体通过梅克尔根影响交易ID
		}

		blockInfo := bytes.Join(tmp, []byte{})

		//2. sha256
		hash = sha256.Sum256(blockInfo)
		//3. 与pow中的target进行比较
		hashInt := big.Int{}
		hashInt.SetBytes(hash[:])
		//fmt.Printf("big.Int 中的对象以十六进制输出:\n%x\n",hashInt.Abs(&hashInt))
		//4. 判断hashInt中的数是否小于pow.target中的数
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash: %x  nonce: %d\n", hash, nonce)
			return hash[:], nonce
		} else {
			nonce++
		}
	}

}
