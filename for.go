package conc

/*
	concurrent for loop - numWorkers iterations execute in parallel
*/
func ForChunk(inputs <-chan Box, foo func(i Box), numWorkers int) (wait func()) {
	///	workerInputs := SafeChan(inputs);

	var multiInputs MultiReader;
	multiInputs.ch = inputs;

	block := make(chan bool, numWorkers);
	for j := 0; j < numWorkers; j++ {
		go func() {
			/*
				//SafeChan involves a busy(Goshed())wait, which makes people angry
				myInput := SafeChan(inputs);
				for i := range myInput {
					foo(i);
				}
			*/

			for i, done := multiInputs.read(); !done; i, done = multiInputs.read() {
				foo(i)
			}
			block <- true;
		}()
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
	exhausted := make(chan bool);
	go func() {
		for i := range inputs {
			count++;
			go func(i Box) {
				foo(i);
				block <- true;
			}(i);
		}
		exhausted <- true;
	}();
	wait = func() {
		<-exhausted;
		for i := 0; i < count; i++ {
			<-block
		}
	};
	return wait;
}
