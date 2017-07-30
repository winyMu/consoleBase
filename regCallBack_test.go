package consoleBase

import (
	"testing"
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

func TestExample(t *testing.T) {
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
