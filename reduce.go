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
	
	countSignal := Chain(in, ready);
	
	doneReading := false;
	count := 0;
	folds := 0;
	
	result = <- ready;
	for !doneReading || folds != count {
		select {
		//when this fires, all inputs are in ready and we can see how many there were
		case count = <- countSignal:
			doneReading = true;
		//when this fires, we can fold two values (in another goroutine)
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
