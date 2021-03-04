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