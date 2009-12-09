package conc

import (
	"sync";
)

type safeCounter struct {
	count int;
	m sync.Mutex;
}
func (s *safeCounter) incr() {
	s.m.Lock();
	s.count++;
	s.m.Unlock();
}
func (s *safeCounter) decr() {
	s.m.Lock();
	s.count--;
	s.m.Unlock();
}
func (s *safeCounter) val() int {
	return s.count;
}

//try to read a value off a channel until foo returns false
func grabUntil(in chan Box, foo func() bool) (v Box, ok bool) {
	for {
		if v, ok = <- in; ok {
			return;
		}
		if !foo() {
			break;
		}
	}
	return;
}

/*
	concurrent folding - foo is applied to pairs of items and the results are put
	back in the list, until there is only one item in the list.
*/

func Fold(foo func(Box, Box) Box, in chan Box) Box {
	ready := make(chan Box);
	
	//dump all values from in into ready
	go func() {
		for i := range in {
			ready <- i;
		}
	}();
	
	//keep track of how many values are not in the ready queue, but will be eventually
	inProcess := new(safeCounter);
	
	//if this function returns false, all relevant values are currently in the ready queue
	moreComing := func() bool {
		return !closed(in) || inProcess.val() > 0;
	};
	
	for {
		first := <- ready;
		
		second, ok := grabUntil(ready, moreComing);
		
		//if there aren't any more values coming, then we have completed the folding process
		if !ok {
			return first;
		}
		
		//indicate that there is a value not yet in the ready queue that needs to be processed
		inProcess.incr();
		go func() {
			v := foo(first, second);
			ready <- v;
			inProcess.decr();
		}();
	}
	
	return <-ready;
}
