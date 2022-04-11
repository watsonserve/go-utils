package goutils

type Any_t interface{}

type Pool_t struct {
	taskPipe   chan Any_t
	notifyPipe chan Any_t
	count      int
}

type Poolable interface {
	WorkerInit() error
	Worker(ranger Any_t) Any_t
}

func NewPool(poolable Poolable, size int) *Pool_t {
	taskPipe := make(chan Any_t, size)
	notifyPipe := make(chan Any_t, size<<1)

	pool := &Pool_t{
		taskPipe:   taskPipe,
		notifyPipe: notifyPipe,
		count:      0,
	}

	for pool.count = 0; pool.count < size; pool.count++ {
		go pool.worker(taskPipe, notifyPipe, poolable)
	}

	return pool
}

func (pool *Pool_t) worker(taskPipe chan Any_t, notifyPipe chan Any_t, poolable Poolable) {
	err := poolable.WorkerInit()
	if nil == err {
		for ranger := <-taskPipe; nil != ranger; ranger = <-taskPipe {
			notifyPipe <- poolable.Worker(ranger)
		}
	}
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
	return ret
}
