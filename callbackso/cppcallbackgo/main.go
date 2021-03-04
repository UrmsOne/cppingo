/*
@Author: urmsone urmsone@163.com
@Date: 3/4/21 1:53 PM
@Name: main.go
@Description:
*/
package main

/*
#cgo CFLAGS: -I ./
#cgo LDFLAGS: -L./ -lperson
#include <stdlib.h>
#include <stdio.h>
#include "wrapper.h"
*/
import "C"

import (
	"fmt"
)

// export callback
func callback(code int) {
	fmt.Println("Go.callback")
}

// 生成动态库
// g++ person.cpp wrapper.cpp -fPIC -shared -o libperson.so
func main() {
	person := C.call_Person_Create()
	defer C.call_Person_Destroy(person)
	age := C.call_Person_GetAge(person)
	fmt.Println(age)
	//defer C.free(unsafe.Pointer(age))
	name := C.call_Person_GetName(person)
	//defer C.free(unsafe.Pointer(name))
	fmt.Println(C.GoString(name))
}
