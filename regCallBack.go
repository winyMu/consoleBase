package consoleBase

import (
	"nemoCommon/alog"
	"net"
	"strconv"
	"strings"
	"sync"
)

var regCbInstance *RegCallBack
var once sync.Once

func GetRegCbInstance() *RegCallBack {
	return regCbInstance
}

type ConsoleCallBack func(args string) string

type ConsoleEntry struct {
	CmdStr  string
	HelpStr string
	PFunc   ConsoleCallBack
}

type ConsoleCmd2Entry_t map[string]ConsoleEntry

type RegCallBack struct {
	gs_consolecmd2Entry ConsoleCmd2Entry_t
	gs_consoleCmdHelper []string
}

const (
	LISTEN_PORT_RETRY = 20
)

func (rgcb *RegCallBack) RegConsoleCallBack(fp ConsoleCallBack, cmdStr string, helpStr string) bool {

	alog.Infof("RegConsoleCallBack cmdStr %s helpStr %s", cmdStr, helpStr)
	if _, ok := rgcb.gs_consolecmd2Entry[cmdStr]; ok {
		alog.Infof("cmdStr %s has exist ", cmdStr)
		return false
	}

	rEntry := ConsoleEntry{}
	rEntry.CmdStr = cmdStr
	rEntry.HelpStr = helpStr
	rEntry.PFunc = fp

	rgcb.gs_consolecmd2Entry[cmdStr] = rEntry

	rgcb.gs_consoleCmdHelper = append(rgcb.gs_consoleCmdHelper, helpStr)
	return true
}

func (rgcb *RegCallBack) onHandlerConsoleCmd(conn net.Conn, params string, cmdStr string) bool {
	alog.Infof("come in params %s cmdStr %s", params, cmdStr)
	if rEntry, ok := rgcb.gs_consolecmd2Entry[cmdStr]; ok {
		alog.Infof("found handler function params %s cmdStr %s ", params, cmdStr)
		resultStr := rEntry.PFunc(params)
		resultStr += "\r\n"
		alog.Infof("Begin to write right cmd ack")
		n, err := conn.Write([]byte(resultStr))
		alog.Infof("write %v bytes err:%v", n, err)
		alog.Infof("After to write right cmd ack")
		return true
	} else {
		return rgcb.onReservedCmd(params, cmdStr)
	}
}

func (rgcb *RegCallBack) onReservedCmd(params string, cmdStr string) bool {
	alog.Errorf("not found handler function params %s cmdStr %s", params, cmdStr)
	return false
}

func (rgcb *RegCallBack) onConsoleHelp(params string) string {
	var helpStr string = ""
	for k, cmdStr := range rgcb.gs_consoleCmdHelper {
		helpStr += strconv.Itoa(k) + ". " + cmdStr + "\r\n"
	}
	return helpStr
}

func (rgcb *RegCallBack) processCommand(cmd string, conn net.Conn) {
	firstIndex := strings.Index(cmd, " ")

	cmdStr := ""
	params := ""

	if firstIndex == -1 {
		cmdStr = cmd
	} else {
		cmdStr = cmd[0:firstIndex]
		params = cmd[firstIndex+1:]
	}

	if false == rgcb.onHandlerConsoleCmd(conn, params, cmdStr) {
		response := "Invalid Command\r\n>"
		alog.Infof("Begin to write wrong cmd ack")
		n, err := conn.Write([]byte(response))
		alog.Infof("write %v bytes err:%v", n, err)
		alog.Infof("After to write wrong cmd ack")
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		// read from the connection
		var buf = make([]byte, 65535)
		alog.Info("start to read from conn")
		n, err := conn.Read(buf)
		if err != nil {
			alog.Infof("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			break
		}
		if n < 3 || buf[0] == 255 && buf[1] == 244 && buf[2] == 255 { // 针对ctrl+c 断开连接
			break
		} else {
			cmdStrs := string(buf[:n-1]) // 去除"\n"两个字符
			if buf[n-2] == byte('\r') {
				cmdStrs = string(buf[:n-2])
			}
			alog.Infof("read %d bytes, content is %s bytes:%v\n", n, cmdStrs, []byte(cmdStrs))

			GetRegCbInstance().processCommand(cmdStrs, conn)
		}
	}
}

func InitConsolePort(consolePort int) {
	once.Do(func() {
		regCbInstance = &RegCallBack{
			gs_consolecmd2Entry: make(ConsoleCmd2Entry_t),
		}
	})

	regCbInstance.RegConsoleCallBack(regCbInstance.onConsoleHelp, "help", "help")

	var consolePortSuccess int = consolePort
	var listener net.Listener
	var err error
	index := 0
	for index < LISTEN_PORT_RETRY {
		index += 1
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(consolePortSuccess))
		if err != nil {
			alog.Errorf("listen error:", err)
			consolePortSuccess += 1
			continue
		} else {
			break
		}
	}
	if index == LISTEN_PORT_RETRY {
		alog.Fatal("Listen ConsolePort Fail, OVER 20 times !!!")
	}
	alog.Infof("Listen Success Console Port:%v", consolePortSuccess)
	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				alog.Errorf("accept error:", err)
				break
			}
			// start a new goroutine to handle
			// the new connection.
			go handleConn(c)
		}
	}()
}
