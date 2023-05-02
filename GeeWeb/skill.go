package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// todo builderConcat 优先使用次方法追加字符串，减少内存，强烈不能使用 + 和 fmt.Sprintf
func builderConcat(n int, str string) string {
	a := []int{1}
	fmt.Println(cap(a))
	var builder strings.Builder
	builder.Cap()
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}

func format(a []int) {
	fmt.Printf("a: %v， ptr(a) = %x, len(a) = %d, cap(a) = %d\n", a, &a, len(a), cap(a))
}

type Student struct {
	Name string
	Age  int
}

func sam_reflect() {
	student := Student{}
	stu := reflect.TypeOf(student)
	fmt.Println(stu.Name())
	fmt.Println(stu.String())
	fmt.Println(stu.Kind())
	fmt.Println(stu.NumField()) // 字段数量
	fmt.Println(stu)

	value := reflect.Indirect(reflect.ValueOf(&student))
	fmt.Println(value)

}

type Class struct {
	Name    string `json:"Name"`
	Monitor string `json:"Monitor"`
}

func to_json() {
	//cls := Class{"", ""}
	//c := easyjson.Class{"", ""}
	//c.

}

type demo1 struct {
	a int8
	b int16
	c int32
}

type demo2 struct {
	a int8
	b int16 //对齐倍数2
	c int32 //对齐倍数4
	d int8  //对齐倍数4
}

func demo() {
	fmt.Println(unsafe.Sizeof(demo1{})) // 2+2+4 = 8  占用的字节数
	fmt.Println(unsafe.Sizeof(demo2{})) // 2+2+4+4 =12

	fmt.Println(unsafe.Sizeof(struct{}{})) //0
	fmt.Println(unsafe.Sizeof(true))       // 1
	buf := make([]byte, 0, 100000001)
	for i := 0; i < 100000000; i++ {
		buf = append(buf, "a"...)
	}
	fmt.Println(unsafe.Sizeof(string(buf)))

}

func main() {
	fmt.Println("**********************************")
	//a := []int{1,2,3}
	//format(a)
	//a = append(a, 4,5)
	//format(a)

	//sam_reflect()

	demo()
	fmt.Println("**********************************")
}
