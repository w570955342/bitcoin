package main

import (
	"bolt"
	"log"
	"fmt"
)

func main() {
	//1. 打开数据库
	db,err:=bolt.Open("test.db",0600,nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()


	db.Update(func(tx *bolt.Tx) error {
		//2. 找到抽屉,没有就创建
		bucket:=tx.Bucket([]byte("b"))
		if bucket==nil {//没有抽屉
			bucket,err=tx.CreateBucket([]byte("b"))
			if err != nil {
				log.Panic(err)
			}
		}

		//3. 往bucket写数据
		bucket.Put([]byte("乔布斯"),[]byte("苹果"))
		bucket.Put([]byte("中本聪"),[]byte("比特币"))
		return nil
	})

	//4. 从bucket读数据
	db.View(func(tx *bolt.Tx) error {
		//1. 找到抽屉
		bucket:=tx.Bucket([]byte("b"))
		if err != nil {
			log.Panic(err)
		}
		//2. 读取数据
		value1:=bucket.Get([]byte("乔布斯"))
		value2:=bucket.Get([]byte("中本聪"))
		fmt.Printf("value1: %s\n",value1)
		fmt.Printf("value2: %s\n",value2)
		return nil
	})
}
