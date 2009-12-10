package conc

/*
	Returns a thunk that, when evaluated, will return the result of
	foo(), waiting if necessary.
*/
func Future(foo func() Box) Thunk {
	wormhole := make(chan Box, 1);
	go func() {
		wormhole <- foo();
	}();
	return func() (result Box) {
		result = <- wormhole;
		wormhole <- result;
		return;
	};
}

func FutureChan(foo func() Box) ThunkChan {
	wormhole := make(ThunkChan);
	go func() {
		result := foo();
		for {
			wormhole <- result;
		}
	}();
	return wormhole;
}