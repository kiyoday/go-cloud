package main

import (
	"fmt"
	"unsafe"
)

const limit = 512
const top uint16 = 1421
const Pi float64 = 3.1415926
const x,y int = 1,3 //多重赋值
//自定义错误类型
type myError struct {
	arg  int
	errMsg string
}

func printA(){
	var i int = 7
	fmt.Println("length of top: ", unsafe.Sizeof(top))
	var j string = "kiyo"
	fmt.Println("length of Pi: ", unsafe.Sizeof(Pi))
	fmt.Println(i,j,"Hello，world！")
	t1 := "\"hello\""
	fmt.Println(t1)
}

func main() {

}