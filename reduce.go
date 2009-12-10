package conc

/*
	concurrent reduce - foo is applied to pairs of items and the results are put
	back in the list, until there is only one item in the list. The operator provided
	must be both associative and commutative.
*/
func Reduce(foo func(Box, Box) Box, in chan Box, init Box) (result Box) {
	ready := make(chan Box);
	
	go func() {
		ready <- init;
	}();
	
	inserted := make(chan bool);
	count := 0;
	
	doneSignal := make(chan bool);
	
	//dump all values from in into ready
	go func() {
		for x := range in {
			count++;
			go func(x Box) {
				ready <- x;
				inserted <- true;
			}(x);
		}
		//wait for all of the values read to get onto ready
		for i:=0; i<count; i++ {
			<-inserted;
		}
		//send a signal indicating that the input is exhausted
		doneSignal <- true;
	}();
	
	doneReading := false;
	folds := 0;
	
	result = <- ready;
	for !doneReading || folds != count {
		select {
		case <-doneSignal:
			doneReading = true;
		case second := <- ready:
			folds++;
			go func(first, second Box) {
				ready <- foo(first, second);
			}(result, second);
			result = <- ready;
		}
	}
	
	return;
}
