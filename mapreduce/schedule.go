package mapreduce

import (
	"fmt"
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//

	//Remember that workers may fail, and that any given worker may finish
	//	// multiple tasks: the server cannot reach workers only by rpc call
	var wg sync.WaitGroup
	wg.Add(ntasks)
	task := 0
	for task < ntasks {
		go func(i int) {
			arguments := DoTaskArgs{mr.jobName, mr.files[i], phase, i, nios}
			for {
				worker := <-mr.registerChannel
				ok := call(worker, "Worker.DoTask", &arguments, nil)
				if ok {
					go func() { mr.registerChannel <- worker }()
					wg.Done()
					break
				}
			}
		}(task)
		task = task + 1
	}
	wg.Wait()
	fmt.Printf("Schedule: %v phase done\n", phase)
}
