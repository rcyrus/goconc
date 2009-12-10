package conc

import "sync"

/*
	concurrent for loop - numWorkers iterations execute in parallel
*/
func ForChunk(inputs <-chan Box, foo func(i Box), numWorkers int) (wait func()) {
	m := new(sync.Mutex);
	
	block := make(chan bool, numWorkers);
	for j := 0; j < numWorkers; j++ {
		go func() {
			for {
				m.Lock();
				i, done := <- inputs, closed(inputs);
				m.Unlock();
				if done {
					break;
				}
				foo(i);
			}
			block <- true;
		}();
	}
	wait = func() {
		for i := 0; i < numWorkers; i++ {
			<-block
		}
	};
	return wait;
}

func For(inputs <-chan Box, foo func(i Box)) (wait func()) {
	count := 0;
	block := make(chan bool);
	for i := range inputs {
		count++;
		go func(i Box) {
			foo(i);
			block <- true;
		}(i);
	}
	wait = func() {
		for i := 0; i < count; i++ {
			<-block
		}
	};
	return wait;
}