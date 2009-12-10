package conc

/*
	concurrent folding - foo is applied to pairs of items and the results are put
	back in the list, until there is only one item in the list. The operator provided
	must be both associative and commutative.
*/
func Fold(foo func(Box, Box) Box, in chan Box) (result Box) {
	ready := make(chan Box);
	
	inserted := make(chan bool);
	count := 0;
	
	doneSignal := make(chan bool);
	
	//dump all values from in into ready
	go func() {
		for {
			if v, isLast := <-in, closed(in); !isLast {
				count++;
				go func() {
					ready <- v;
					inserted <- true;
				}();
			}
			else {
				break;
			}
		}
		for i:=0; i<count; i++ {
			<-inserted;
		}
		doneSignal <- true;
	}();
	
	doneReading := false;
	folds := 0;
	for {
		first := <- ready;
		if doneReading && folds == count-1 {
			result = first;
			break;
		}
		select {
		case second := <- ready:
			folds++;
			go func() {
				ready <- foo(first, second);
			}();
			break;
		case <-doneSignal:
			doneReading = true;
			go func() {
				ready <- first;
			}();
			break;
		}
	}
	
	return;
}
