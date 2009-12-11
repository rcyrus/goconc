package conc

import (
	"sync";
	"runtime";
)

var lookup = make(map[<-chan Box] chan chan Box)
var safeChanLock sync.Mutex

/*
	Returns a channel that will grab values from the input, with
	the idea that you can use range and <-ch,closed(ch) in a threadsafe fashion.
	Multiple calls to SafeChan() with the same input channel will result in
	different channels that can all be read in different goroutines.
*/
func SafeChan(inputs <-chan Box) chan Box {
	safeChanLock.Lock();
	defer safeChanLock.Unlock();
	if outgoing, ok := lookup[inputs]; ok {
		return <-outgoing;
	}

	loop := make(chan chan Box);
	outgoing := make(chan chan Box);
	
	lookup[inputs] = outgoing;
	
	go func() {
		for {
			nextChan := make(chan Box);
			outgoing <- nextChan;
			//once the new channel is grabbed by someone, add it to the loop
			loop <- nextChan;
		}
	}();
	
	go func() {
		for v := range inputs {
			for receiver := range loop {
				go func(receiver chan Box) {
					loop <- receiver;
				}(receiver);
				if receiver <- v {
					break;
				}
				runtime.Gosched();
			}
		}
		
		for receiver := range loop {
			close(receiver);
		}
		
		lookup[inputs] = nil, false;
	}();
	
	return <-outgoing;
}
