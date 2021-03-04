# Go进阶编程之Go调用C++（Linux）

## 环境：
* Ubuntu 18.04
* golang 1.14
* linux amd64
## 一、了解调用流程
* c调用c++动态链接库
* go通过cgo调用c
* 从而达到go调用c++，其实是通过c来做中间层转换的功能

### go调用c动态库
number目录结构
```
root@ubuntu:/workspace/gospace/cppingo/callso/number# tree
.
├── libnumber.so
├── main.go
├── number.c
└── number.h
```
number.h
```
int number_add_mod(int a, int b, int mod);
number.c
#include "number.h"
int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}
```
main.go
```
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
```
动态链接库生成
```
gcc -shared -o libnumber.so number.c
```
运行
```
root@ubuntu:/workspace/gospace/cppingo/callso/number# go build -o main main.go
root@ubuntu:/workspace/gospace/cppingo/callso/number# ./main
3
```
如果报错：
```
./main: error while loading shared libraries: libnumber.so: cannot open shared object file: No such file or directory
```
是链接时没找到libnumber.so这个库，解决办法，指定链接时搜索的目录即可，其中一个办法如下:
```
步骤：
1）新建文件/etc/ld.so.conf.d/test.conf
2）在文件中指定链接搜索路径
3）保存后执行ldconfig
我的路径设置如下：
root@ubuntu:/etc/ld.so.conf.d# ldconfig
root@ubuntu:/etc/ld.so.conf.d# cat test.conf 
/workspace/gospace/cppingo/callso/number
/workspace/gospace/cppingo/callso/callcppso
```

### c调用c++动态链接库
ccallcpp目录
```
root@ubuntu:/workspace/gospace/cppingo/callso/ccallcpp# tree
.
├── libperson.so
├── main.go
├── Makefile
├── person.cpp
├── person.h
├── wrapper.cpp
└── wrapper.h
```
person.h
```
#ifndef PERSON_H_
#define PERSON_H_
#include <string>
class Person {
 public:
  Person(std::string name, int age);
  ~Person() {}
  const char *GetName() { return name_.c_str(); }
  int GetAge() { return age_; }
 private:
  std::string name_;
  int age_;
};
#endif // PERSON_H_
```
person.cpp
```
#include <iostream>
#include "person.h"
Person::Person(std::string name, int age)
    : name_(name), age_(age) {}
```
由于c没有类的概念，因此需要写一个wrapper把c++中的类转成c的函数来调用
wrapper.h
```
#ifndef WRAPPER_H_
#define WRAPPER_H_
#ifdef __cplusplus
extern "C"
{
#endif
void *call_Person_Create();
void call_Person_Destroy(void *);
int call_Person_GetAge(void *);
const char *call_Person_GetName(void *);
#ifdef __cplusplus
}
#endif
#endif // WRAPPER_H_
```
wrapper.cpp
```
#include "person.h"
#include "wrapper.h"
#ifdef __cplusplus
extern "C"{
#endif
void *call_Person_Create() {
  return new Person("urmsone", 18);
}
void call_Person_Destroy(void *p) {
  delete static_cast<Person *>(p);
}
int call_Person_GetAge(void *p) {
  return static_cast<Person *>(p)->GetAge();
}
const char *call_Person_GetName(void *p) {
  return static_cast<Person *>(p)->GetName();
}
#ifdef __cplusplus
}
#endif
```
main.go
```
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
```
生成动态库，动态库的位置和名字要于main函数中的`cgo LDFLAGS: -L./ -lperson`对应；
```
g++ person.cpp wrapper.cpp -fPIC -shared -o libperson.so
```

[Go语言高级编程](https://chai2010.cn/advanced-go-programming-book/ch2-cgo/ch2-02-basic.html)
提到：
```
CFLAGS部分，-I定义了头文件包含的检索目录。LDFLAGS部分，-L指定了链接时库文件检索目录，-l指定了链接时需要链接person库。
因为C/C++遗留的问题，C头文件检索目录可以是相对目录，但是库文件检索目录则需要绝对路径。在库文件的检索目录中可以通过${SRCDIR}变量表示当前包目录的绝对路径：
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
```
但我这里使用的是相对路径，也没出问题。

运行
```
root@ubuntu:/workspace/gospace/cppingo/callso/ccallcpp# g++ person.cpp wrapper.cpp -fPIC -shared -o libperson.so
root@ubuntu:/workspace/gospace/cppingo/callso/ccallcpp# go build -o main main.go 
root@ubuntu:/workspace/gospace/cppingo/callso/ccallcpp# ./main 
./main: error while loading shared libraries: libperson.so: cannot open shared object file: No such file or directory
root@ubuntu:/workspace/gospace/cppingo/callso/ccallcpp# ./main 
18
urmsone
```

**流程：go->cgo->c->wrapper-c++**

## 二、回调
###描述
主要记录一个常用的场景，参考官网[链接](https://github.com/golang/go/wiki/cgo#function-pointer-callbacks) ，描述如下：
go调用c/c++的某个函数A（c/c++语言实现），然后把go语言实现的函数B作为参数传递给函数A，在执行函数A的过程中调用函数B。这个流程包含了go调用c/c++，和c/c++调用go两种情况。
调用过程：Go.main -> C.A -> Go.B
###场景：
go需要调用原生C++动态链接库，而动态库中某个函数A需要接收《函数指针》变量最为参数，通过《函数指针》来调用回调函数。当go调用函数A时，需要传递一个go语言实现的函数，并把该函数的指针传递给A。

###实现
ccallbackgo目录：
```
root@ubuntu:/workspace/gospace/cppingo/callbackso/ccallbackgo# tree
.
├── cfunc.go
├── clibrary.c
├── clibrary.h
├── libclibrary.so
└── main.go
```

clibrary.h
```
// 定义函数some_c_func和
// 函数指针（简单理解为回调函数的地址，通过该地址可之间调用回调函数）
#ifndef CLIBRARY_H
#define CLIBRARY_H
typedef int (*callback_fcn)(int);
void some_c_func(callback_fcn);
#endif
```

clibrary.h
```
// 实现some_c_func函数
#include <stdio.h>
#include "clibrary.h"
void some_c_func(callback_fcn callback){
    int arg = 2;
    printf("C.some_c_func(): calling callback with arg =%d\n", arg);
    int response = callback(2);
    printf("C.some_c_func(): callback response with %d\n", response);
}
```
main.go
```
/*
@Author: urmsone urmsone@163.com
@Date: 3/4/21 7:37 PM
@Name: main.go
@Description:
1）在注释中声明回调函数callOnMeGo_cgo（该名字应与网关函数同名）
2）实现回调函数,通过cgo的export callOnMeGo注释导出
3）编写网关函数
*/
package main
/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lclibrary
#include "clibrary.h"
int callOnMeGo_cgo(int in); // Forward declaration.
*/
import "C"
import (
	"fmt"
	"unsafe"
)
//export callOnMeGo
func callOnMeGo(in int) int {
	fmt.Printf("Go.callOnMeGo(): called with arg = %d\n", in)
	return in + 1
}
func main() {
	fmt.Printf("Go.main(): calling C function with callback to us\n")
        // C.callOnMeGo_cgo为注释中声明的回调函数名（网关函数名）
        // some_c_func里面的涵义是，先把C.callOnMeGo_cgo函数转成指针
        // 再通过C.callback_fcn把指针类型转成回调指针函数的类型
	C.some_c_func((C.callback_fcn)(unsafe.Pointer(C.callOnMeGo_cgo)))
}
```

cfunc.go
```
/*
@Author: urmsone urmsone@163.com
@Date: 3/4/21 7:45 PM
@Name: cfunc.go
@Description:
*/
package main
/*
#include <stdio.h>
// The gateway function
int callOnMeGo_cgo(int in){
	printf("C.callOnMeGo_cgo(): called with arg = %d\n", in);
	int callOnMeGo(int);
	return callOnMeGo(in);
}
*/
import "C"
```

生成动态链接库
```
gcc -shared -o libclibrary.so clibrary.c

运行
root@ubuntu:/workspace/gospace/cppingo/callbackso/ccallbackgo# go build -o main *.go
root@ubuntu:/workspace/gospace/cppingo/callbackso/ccallbackgo# ./main
Go.main(): calling C function with callback to us
C.some_c_func(): calling callback with arg =2
C.callOnMeGo_cgo(): called with arg = 2
Go.callOnMeGo(): called with arg = 2
C.some_c_func(): callback response with 3
注意：要把网关函数一起build，否则会报错
root@ubuntu:/workspace/gospace/cppingo/callbackso/ccallbackgo# go build -o main main.go 
# command-line-arguments
/tmp/go-build285429389/b001/_cgo_main.o:/tmp/go-build/cgo-generated-wrappers:2: undefined reference to `callOnMeGo_cgo'
```

### 三、内存管理
* go调c++
    * 何时释放由c/c++开辟的内存空间
* c++调go
    * 何时释放go回调回来由go创建的对象（go有GC，可能c/c++ keep住的对象已被GC回收，导致内存地址invalid）

### 四、c++语法学习
将c++类转成c函数调用
```
void *call_Person_Create() {
  return new Person("urmsone", 18);
}
void call_Person_Destroy(void *p) {
  delete static_cast<Person *>(p);
}
```

`void *call_Person_Create() `是一个指针函数，返回c++中Person类对象的指针，指针类型为void（因为c中没有c++中Person类对应的类型，所以类型为void）。

这里涉及到c++类型转换符，static_cast,const_cast,reinterpret_cast,dynamic_cast
static_cast<Person *>(p)

回调类型：`typedef void (*EnvSDK_Callback)(int, const char*, void*);`

### 代码仓库
[源码仓库](https://github.com/UrmsOne/cppingo)

### 参考：
[Go语言高级编程-gitbook](https://chai2010.cn/advanced-go-programming-book/)  
