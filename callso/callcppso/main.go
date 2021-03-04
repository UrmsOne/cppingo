/*
@Author: urmsone urmsone@163.com
@Date: 3/4/21 1:43 PM
@Name: main.go.go
@Description:
*/
package main

/*
#cgo CFLAGS: -I./
#cgo LDFLAGS: -L./ -lstu
#include <stdlib.h>
#include <stdio.h>
#include "wrapper.h" //非标准c头文件，所以用引号
*/
import "C"

import "unsafe"

func main() {
	name := "test"
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	obj:= C.stuCreate()
	C.initName(obj, cStr)
	C.getStuName(obj)
}
