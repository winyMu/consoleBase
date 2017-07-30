# 功能描述
通过console端口，查看进程内部数据

# 接口使用

1.初始化 <br>
初始化一个console端口，如果被占用会尝试最多20次，仍然失败会panic <br>
func `InitConsolePort`(consolePort int) <br>
<br>
2.函数注册 <br>
支持 **普通函数** 和 **strcut的成员函数** <br>
函数签名如下：一个string类型参数，一个返回值 <br>
type `ConsoleCallBack` func(args string) string <br>
**参数说明**：<br>
**args**    -> 注册函数的各种参数， 以分隔符分割， 例如空格" "，冒号:等<br>
<br>
func `RegConsoleCallBack`(fp ConsoleCallBack, cmdStr string, helpStr string) bool <br>
**参数说明**：<br>
**fp**      -> ConsoleCallBack 类型的函数，strcut的成员函数一样处理 <br>
**cmdStr**  -> 命令参数，唯一标识被注册的函数fp，连上console端口后，可以被访问 <br>
**helpStr** -> 函数功能帮助说明，一般写成 cmdStr + 函数参数 形式 <br>

```
import (
	"nemoCommon/consoleBase"
)
func testConsolFunc1(params string) string {
	return "testConsolFunc1!!!"
}

func testConsolFunc2(params string) string {
	return "testConsolFunc2!!!"
}

type consoleDemo struct {
}

func (*consoleDemo) demoStructFunc(params string) string {
	return "test demoStructFunc"
}

func main() {
	consoleDemo := &consoleDemo{}

	// 初始化端口
	consoleBase.InitConsolePort(8888)
	// 注册普通方法
	consoleBase.GetRegCbInstance().RegConsoleCallBack(testConsolFunc1, "test1", "test1 param1 param2")
	consoleBase.GetRegCbInstance().RegConsoleCallBack(testConsolFunc2, "test2", "test2 param1 param2")
	// 注册成员方法
	consoleBase.GetRegCbInstance().RegConsoleCallBack(consoleDemo.demoStructFunc, "testStructFunc", "testStructFunc param1 param2 param3")

	time.Sleep(time.Minute * 2)
}
```

# 使用方法
注册完成之后，通过telnet连接去处理 <br>
敲击help， 就能看到你注册的所有函数，通过敲击命令来调用你注册的函数 enjoy it !!! <br>
```
root@ubuntu:/mnt/hgfs/E/service_platform/src/nemoServers/consoleDemo# telnet 0 8888
Trying 0.0.0.0...
Connected to 0.
Escape character is '^]'.
help
0. help
1. test1 param1 param2
2. test2 param1 param2
3. testStructFunc param1 param2 param3

```
