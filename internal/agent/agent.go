package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/issueye/tail"
)

type Agent struct {
	ID         string        // 监听名称
	Path       string        // 日志路径
	ScriptPath string        // 脚本路径
	Level      int           // 日志等级 注：需要监听的日志等级，如果没有设置则是所有等级的日志内容都被监听
	Close      chan struct{} // 结束
}

// 中间对象管道
var ChanAgent = make(chan *Agent, 2)

// 保存对象
var MapAgent = new(sync.Map)

// 初始化
func Init() {
	fmt.Println("开启监听代理")
	go func() {
		for {
			select {
			case a := <-ChanAgent:
				{
					Listen(a)
				}
			}
		}
	}()
}

// 单独监听代理
func Listen(a *Agent) {
	_, ok := MapAgent.Load(a.ID)
	if ok {
		fmt.Printf("监听【%s】已经添加，请勿重新添加\n", a.ID)
		return
	} else {
		MapAgent.Store(a.ID, a)
		go listen(a)
	}
}

func listen(a *Agent) {
	// 设置跟踪
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	tails, err := tail.TailFile(a.Path, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}

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

				fmt.Println("line:", line.Text)

				// TODO 读出的内容，通过脚本进行处理
			}
		case <-a.Close: // 关闭监听
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
