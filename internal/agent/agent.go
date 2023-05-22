package agent

import (
	"fmt"
	"sync"
	"time"

	lichee "github.com/issueye/lichee-js"
	"github.com/issueye/log-agent/internal/config"
	"github.com/issueye/log-agent/pkg/ws"
	"github.com/issueye/tail"
)

type Message struct {
	Id   string
	Data string
}

type getMessageFunc func(string) string

type Agent struct {
	ID         string        // 监听名称
	Path       string        // 日志路径
	ScriptPath string        // 脚本路径
	Level      int           // 日志等级 注：需要监听的日志等级，如果没有设置则是所有等级的日志内容都被监听
	close      chan struct{} // 结束
	Message    chan *Message // 消息
	JsCore     *lichee.Core
	getMessage getMessageFunc
}

func New(id, path, script string, level int) (*Agent, error) {
	a := &Agent{
		ID:         id,
		Path:       path,
		ScriptPath: script,
		Level:      level,
		close:      make(chan struct{}),
		Message:    make(chan *Message, 100),
	}

	MapAgent.Store(id, a)

	r := config.GetParam("SERVER-MODE", "release")
	if r.Value == "release" && a.ScriptPath != "" {
		a.JsCore = lichee.NewCore()
		a.getMessage = a.jsRt(a.JsCore)
	}

	return a, nil
}

func (a *Agent) GetLogLevel() string {
	switch a.Level {
	case -1:
		return "debug"
	case 0:
		return "info"
	case 1:
		return "warn"
	case 2:
		return "error"
	default:
		return "debug"
	}
}

// Del
// 删除代理
func Del(id string) {
	value, ok := MapAgent.Load(id)
	if ok {
		a := value.(*Agent)
		a.Close()
	}
}

// CloseAgent
// 关闭
func (a *Agent) Close() {
	a.close <- struct{}{}
	close(a.Message)
	close(a.close)
	MapAgent.Delete(a.ID)
}

// 保存对象
var MapAgent = new(sync.Map)

// 单独监听代理
func (a *Agent) Listen() {
	go a.listen()
}

func (a *Agent) listen() {
	// 设置跟踪
	cfg := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	tails, err := tail.TailFile(a.Path, cfg)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}

	fmt.Println("开启监听：", a.Path)

	// 开始监听文件
	for {
		select {
		case line, ok := <-tails.Lines: // 读内容
			{
				//遍历chan，读取日志内容
				if !ok {
					fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
					time.Sleep(time.Second)
					continue
				}

				a.Message <- &Message{
					Id:   a.ID,
					Data: line.Text,
				}
			}
		case m := <-a.Message:
			{
				// 是否使用异步的方式
				r := config.GetParam("LOG-AGENT-ASYNC", "false")
				if r.Bool() {
					go a.NewLine(m.Data)
				} else {
					a.NewLine(m.Data)
				}

			}
		case <-a.close: // 关闭监听
			{
				// 移除掉map中的内容
				MapAgent.Delete(a.ID)
				goto title
			}
		}
	}

title:
	fmt.Println("监听退出")
}

func (a *Agent) NewLine(data string) {
	var callbackData string

	// 如果没有设置脚本路径，则不使用脚本处理
	if a.ScriptPath != "" {
		// 如果当前运行模式是 debug 模式则每次都会重新加载JS文件，适合在调试时
		r := config.GetParam("SERVER-MODE", "release")
		if r.Value == "debug" {
			c := lichee.NewCore()
			cbFunc := a.jsRt(c)
			callbackData = cbFunc(data)
		} else {
			callbackData = a.getMessage(data)
		}

		if callbackData == "" {
			return
		}
	} else {
		callbackData = data
	}

	fmt.Println("line:", callbackData)
	c, ok := ws.SMap.Load(a.ID)
	if ok {
		err := c.(*ws.WsConn).OutChanWrite([]byte(callbackData))
		if err != nil {
			fmt.Printf("向 websocket [%s]发送消息失败，失败原因：%s\n", a.ID, err.Error())
			return
		}
	}
}

func (a *Agent) jsRt(rt *lichee.Core) getMessageFunc {
	// 设置参数
	rt.SetLogOutMode(lichee.LOM_DEBUG)
	rt.SetGlobalProperty("getLogLevel", a.GetLogLevel)
	err := rt.Run(a.ID, a.ScriptPath)
	if err != nil {
		fmt.Println("运行脚本失败，失败原因：", err.Error())
		return nil
	}

	// 导出方法
	var cbFunc getMessageFunc
	err = rt.ExportFunc("getMessage", &cbFunc)
	if err != nil {
		fmt.Println("导出方法失败【getMessage】，失败原因：", err.Error())
		return nil
	}

	return cbFunc
}
