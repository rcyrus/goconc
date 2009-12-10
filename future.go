package conc

/*
	Returns a thunk that, when evaluated, will return the result of
	foo(), waiting if necessary.
*/
func Future(foo func() Box) (thunk func() Box) {
	wormhole := make(chan Box, 1);
	go func() {
		wormhole <- foo();
	}();
	thunk = func() (result Box) {
		result = <- wormhole;
		wormhole <- result;
		return;
	};
	return thunk;
}