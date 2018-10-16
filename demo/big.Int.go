package main

import (
	"math/big"
	"fmt"
)

//将十六进制显示的字符串转化成 big.Int 存储
func main() {
	targetStr:="1f"
	//targetStr:="0000f00000000000000000000000000000000000000000000000000000000000"
	target:=big.Int{}
	target.SetString(targetStr,16)//指定字符串所代表的进制，即"1f"是以16进制的形式展现出来的
	fmt.Println("string:\n",targetStr)
	fmt.Println("==============================================")
	fmt.Println("big.Int:\n",target)
	fmt.Println("==============================================")
	fmt.Println("big.Int 中的对象:\n",target.Abs(&target))
	fmt.Printf("big.Int 中的对象类型：%T\n",target.Abs(&target))
	fmt.Println("==============================================")
	fmt.Printf("big.Int 中的对象以十六进制输出:\n%x\n",target.Abs(&target))
	fmt.Printf("big.Int 中的对象以二进制输出:\n%b\n",target.Abs(&target))
	fmt.Printf("big.Int 中的对象以十进制输出:\n%d\n",target.Abs(&target))
}
