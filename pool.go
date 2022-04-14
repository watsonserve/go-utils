package goutils

type Any_t interface{}

type Pool_t struct {
	taskPipe   chan Any_t
	notifyPipe chan Any_t
	count      int
}

type WorkerGenerator func() (Worker, error)

type Worker interface {
	Work(params Any_t) Any_t
	Destroy()
}

func NewPool(workerInit WorkerGenerator, size int) *Pool_t {
	taskPipe := make(chan Any_t, size)
	notifyPipe := make(chan Any_t, size<<1)

	pool := &Pool_t{
		taskPipe:   taskPipe,
		notifyPipe: notifyPipe,
		count:      0,
	}

	for pool.count = 0; pool.count < size; pool.count++ {
		go pool.thread(taskPipe, notifyPipe, workerInit)
	}

	return pool
}

func (pool *Pool_t) thread(taskPipe chan Any_t, notifyPipe chan Any_t, workerInit WorkerGenerator) {
	worker, err := workerInit()
	if nil == err {
		for params := <-taskPipe; nil != params; params = <-taskPipe {
			notifyPipe <- worker.Work(params)
		}
	}
	worker.Destroy()
	// 得到的任务为nil则传出nil
	notifyPipe <- nil
}

func (pool *Pool_t) Count() int {
	return pool.count
}

func (pool *Pool_t) Push(foo Any_t) {
	pool.taskPipe <- foo
}

func (pool *Pool_t) Wait() Any_t {
	ret := <-pool.notifyPipe
	if nil == ret {
		pool.count--
	}
	return ret
}
