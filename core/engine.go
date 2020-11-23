package core

import (
	"CrawlerX/duck"
	"log"
	"os"
	"sync"
)

var (
	wkSize     = 4
	workers    []*worker
	workerRate int
	wqSize     = 100000
	taskQueue  chan string
	stop       = make(chan struct{})
	repeat     sync.Map
	output     = make(chan duck.Result)
	parsers    map[string]duck.Parser
)

// 推送种子页面
func PushSeed(task string) {
	_, ok := repeat.Load(task)
	if !ok {
		//timeout := time.Tick(time.Millisecond * 100)
		select {
		case taskQueue <- task:
			//log.Printf("New Task => %s", task)
			repeat.Store(task, struct{}{})
		default:
			// cancel
			log.Printf("========= Engine is busy, wait queue size %d =========", len(taskQueue))
		}
	}
}

// 订阅输出结果通道
func SubResult() <-chan duck.Result {
	return output
}

// 初始化引擎
func start() {
	createWks()
	go func() {
		tid := 0
		for {
			select {
			case url := <-taskQueue:
				// 平均分配
				index := tid % len(workers)
				w := workers[index]
				tid++
				// task
				log.Printf("select work %d, task => %s", w.id, url)
				w.in <- url
			case <-stop:
				os.Exit(0)
			}
		}
	}()
	// collect result
	go func() {
		for _, w := range workers {
			select {
			case r := <-w.out:
				select {
				case output <- r:
				default:
					log.Printf("//******// Unhandle Result From Worker %d :: %+v", w.id, r)
				}
			default:
				continue
			}
		}
	}()
}

// 创建工人
func createWks() {
	for i := 0; i < wkSize; i++ {
		w := newWorker(i+1, workerRate)
		workers = append(workers, w)
		w.Run()
	}
}

// 设置解析器
func InitEngine(queueSize, workerSize, wkRate int, ps map[string]duck.Parser) func() {
	parsers = ps
	wkSize = workerSize
	wqSize = queueSize
	workerRate = wkRate
	taskQueue = make(chan string, wqSize) //200m
	return start
}
