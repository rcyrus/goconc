package conc

/*
	Take a channel of thunks and put them into a channel of values.
*/
func Realize(in chan Thunk) chan Box {
	out := make(chan Box);
	
	go func() {
		for i := range in {
			out <- i();
		}
		close(out);
	}();
	
	return out;
}