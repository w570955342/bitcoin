package main

//定一个 Wallet结构，它保存所有的Key以及它的地址
type Wallet struct {
	//map[地址]秘钥
	WalletsMap map[string]*Key
}

//生成钱包
func NewWallet() *Wallet {
	key := NewKey()
	address := key.NewAddress()

	var wallet Wallet
	wallet.WalletsMap = make(map[string]*Key)
	wallet.WalletsMap[address] = key

	return &wallet

}

//保存方法，把新建的key添加进去

//读取文件方法，把所有的key读出来
