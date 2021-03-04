/*
@Author: urmsone urmsone@163.com
@Date: 3/4/21 10:29 AM
@Name: main.go
@Description:
*/
package main

//#include <stdio.h>
import "C"
import "unsafe"

// 生成动态链接库
// g++ student.cpp wrapper.cpp -fPIC -shared -o libstu.so
func main() {
	buf := NewMyBuffer(1024)
	defer buf.Delete()

	copy(buf.Data(), []byte("hello\x00"))
	C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}
