package monitor

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/nichuanfang/gymdl/utils"
)

//目录监控

// WatchManager监听管理器
type WatchManager struct {
	//目录-目录监听器映射表
	watchers map[string]*fsnotify.Watcher
	//互斥锁
	mu sync.Mutex
	//事件缓冲池(相当于消息队列)
	eventCh chan fsnotify.Event
}

// 创建监听管理器
func NewWatchManager() *WatchManager {
	return &WatchManager{
		watchers: make(map[string]*fsnotify.Watcher),
		eventCh:  make(chan fsnotify.Event, 1024),
	}
}

// 添加目录
func (wm *WatchManager) AddDir(dir string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if _, ok := wm.watchers[dir]; ok {
		//当前目录已注册
		return nil
	}

	//添加监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	//监听器注册目录
	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	//将目录和目录监听器绑定 建立一对一关系
	wm.watchers[dir] = watcher

	//启动一个协程 执行监听
	go wm.watchLoop(watcher)
	return nil
}

// 监听
func (wm *WatchManager) watchLoop(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// 可加入去抖动逻辑
			wm.eventCh <- event
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			utils.ErrorWithFormat("watcher error:", err)
		}
	}
}

func (wm *WatchManager) Stop() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 关闭所有 watcher
	for _, watcher := range wm.watchers {
		watcher.Close()
	}

	// 关闭事件通道，通知所有 worker 退出
	close(wm.eventCh)

	// 清空 map
	wm.watchers = make(map[string]*fsnotify.Watcher)
}

// 启动workers
func (wm *WatchManager) StartWorkerPool(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for event := range wm.eventCh {
				utils.DebugWithFormat("[Monitor] Worker %d handling event: %s\n", id, event)
				// todo 处理业务逻辑
				// todo 1如果event对应的是目录 需要解锁该目录下的文件 读取该目录下所有加密文件 用um cli 批量解密后整理
				// todo 2如果event对应的是文件 则直接整理无需解锁
				// todo 3 第1步比第二步多一个解密的环节 整理可以一起做
			}
		}(i)
	}
}
