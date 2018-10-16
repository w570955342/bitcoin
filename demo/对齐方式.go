package main

import (
	"fmt"
	"unsafe" //go语言的sizeof
)
// 内存增加方向:-------->
// 地址：0x100 0x101 0x102 0x103
// 12 34 56 78 -> 大端 -> 高尾端
// 78 56 34 12 -> 小端 -> 低尾端

func main() {
	s := uint32(0x12345678)
	b := uint8(s)//从低位开始截取
	fmt.Println("uint32字节大小为", unsafe.Sizeof(s)) //结果为2
	if b == 0x78 {
		fmt.Println("本机器的字节序方式为：little endian(低尾端)")
	} else {
		fmt.Println("本机器的字节序方式为：big endian(高尾端)")
	}
	//fmt.Printf("%x",s)
}
