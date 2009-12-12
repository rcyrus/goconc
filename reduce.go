package conc

/*
	concurrent reduce - foo is applied to pairs of items and the results are put
	back in the list, until there is only one item in the list. The operator provided
	must be both associative and commutative.
*/
func Reduce(foo func(Box, Box) Box, in chan Box, init Box) (result Box) {
	ready := make(chan Box);

	go func() { ready <- init }();

	countSignal := Chain(in, ready);

	count := -1;	//this value means "unknown"
	folds := 0;

	processChan := make(chan Box);

	reduceTwo := func(thunk Box) {
		first, second := thunk.(func() (Box, Box))();
		ready <- foo(first, second);
	};

	For(processChan, reduceTwo);

	result = <-ready;
	for count < 0 || folds != count {
		select {

		//when this fires, all inputs are in ready and we can see how many there were
		case count = <-countSignal:

		//when this fires, we can fold two values (in another goroutine)
		case second := <-ready:
			folds++;
			first := result;
			processChan <- func() (Box, Box) { return first, second };
			result = <-ready;
		}
	}

	close(processChan);

	return;
}

func ReduceChunk(foo func(Box, Box) Box, in chan Box, init Box, numWorkers int) (result Box) {
	ready := make(chan Box);

	go func() { ready <- init }();

	countSignal := Chain(in, ready);

	count := -1;	//this value means "unknown"
	folds := 0;

	processChan := make(chan Box);

	reduceTwo := func(thunk Box) {
		first, second := thunk.(func() (Box, Box))();
		ready <- foo(first, second);
	};

	ForChunk(processChan, reduceTwo, numWorkers);

	result = <-ready;
	for count < 0 || folds != count {
		select {

		//when this fires, all inputs are in ready and we can see how many there were
		case count = <-countSignal:

		//when this fires, we can fold two values (in another goroutine)
		case second := <-ready:
			folds++;
			first := result;
			processChan <- func() (Box, Box) { return first, second };
			result = <-ready;
		}
	}

	close(processChan);

	return;
}
