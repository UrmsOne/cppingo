/*
@Author: urmsone urmsone@163.com
@Date: 3/4/21 10:29 AM
@Name: main.go
@Description:
*/
package main

//#cgo CFLAGS: -I./
//#cgo LDFLAGS: -L${SRCDIR}/ -lnumber
//
//#include "number.h"
import "C"
import "fmt"


// 生成动态链接库
// gcc -shared -o libnumber.so number.c
func main() {
	fmt.Println(C.number_add_mod(10, 5, 12))
}
