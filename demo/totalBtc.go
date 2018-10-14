package main

import "fmt"

func main() {
	fmt.Println("hello")
	//1. 首先每21万个块减半
	//2. 最初一次奖励50比特币
	//3. 用一个循环来判断，累加

	total := 0.0
	blockInterval := 21.0 //单位是万
	currentReward := 50.0

	//currentReward一直接近0，达到一定精度后，变为0
	for currentReward > 0 {
		//每一个区间内的总量
		amount1 := blockInterval * currentReward
		//currentReward /= 2
		currentReward *= 0.5 //除效率低，使用等价的乘法
		total += amount1
	}

	fmt.Println("比特币总量: ", total, "万")
}
