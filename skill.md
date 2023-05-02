# 学习心得

## string

``
优先使用for遍历 make([]byte, 0, n) append追加字符串，减少内存，强烈不能使用 + 和 fmt.Sprintf
``

## 切片

```
切片容量 cap < 2048，以2倍增长，
当超过2048, 不再以2倍增长，节约内存：
    newcap += (newcap + 3*256) / 4
    
不适合大量删除（使用链表替代）

切片删除元素后，最后面的空余的位置 置空提高垃圾回收

不建议在头部追加元素，时间和空间复杂度均为 O(N)

// lastNumsBySlice 占用大量内存，得不到释放
func lastNumsBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

// lastNumsByCopy 内存占用极少
func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	copy(result, origin[len(origin)-2:])
	return result
}

```

## for 和 range

```
多用for 少用 range
每次迭代的值占用很小的内存，for range 性能相同，
迭代的值占用内存大，直接使用 for遍历，性能比range高达上千倍

range 遍历中修改结构体中的值无效（返回的是拷贝），for遍历中修改值有效
使用指针 在range 遍历中也可以修改值

```

## reflect 反射

```

FieldByName("Name")  Field(0)
反射使用下标遍历比使用字段名查找性能高，按照字段名需要遍历所有字段名

在使用json序列化时，go 自带的 json库 是通过反射实现的
可以用 `easyjson` 库代替，性能提升5倍左右

将FieldByName 对应的 字段名和索引 缓存起来，避免每次使用 FieldByName时反复遍历，可以此方式提高查找性能

```

## struct 结构体

```
可以使用 空结构体 当作占位符，不占内存， 字节数：0

使用空结构体：
    比如在构造一个 set 结合，可以将map 的值设置为 空结构体
    不发送消息的信道设置为空结构体
```

## 内存对齐

```
减少cpu读取内存次数，提高吞吐量
在结构体中 字段的类型，按照从小到大字节数排，int8，int16，int32， int64=int
可以减小占用内存
```

## 协程如何退出

[超时:](!https://geektutu.com/post/hpg-timeout-goroutine.html)

```
make(chan bool, 1) 设置带有缓存的通道：

    done := make(chan bool, 1)
	go f(done)
	select {
	case <-done:
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):  // 超时退出，并报错
		return fmt.Errorf("timeout")
	}

分段执行
	
```

其他：

```
close(ch)
chan 多次被关闭，会 panic
```

[优雅的关闭chan](!https://gfw.go101.org/article/channel-closing.html)

## 控制并发数量

```
使用带缓存的chan

利用第三方库 协程池
Jeffail/tunny    https://github.com/Jeffail/tunny
panjf2000/ants   https://github.com/panjf2000/ants

调整系统上限
ulimit

虚拟内存 // 短时间内需要较大内存时，可以使用，一般不推荐（使用硬盘临时充当内存）
```

## sync.Pool

``减少内存分配，降低GC压力``

## sync.Once

```
使函数只执行一次，例如单例模式
并发情况下，是线程安全的
用途初始化
var once sync.Once
once.Do(func(){})
```

## sync.Cond

``协程等待通知``
