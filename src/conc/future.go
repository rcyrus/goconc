package conc

/*
	Returns a thunk that, when evaluated, will return the result of
	foo(), waiting if necessary.
*/
func Future(foo func() Box) (thunk func() Box) {
	wormhole := make(chan Box);
	go func() {
		wormhole <- foo();
	}();
	var result Box;
	thunk = func() Box {
		if closed(wormhole) {
			return result;
		}
		result = <- wormhole;
		return result;
	};
	return thunk;
}