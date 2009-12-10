package conc

/*
	Take a channel of thunks and put them into a channel of values.
*/
func Realize(thunks chan Thunk) chan Box {
	out := make(chan Box);
	
	go func() {
		for thunk := range thunks {
			out <- thunk();
		}
		close(out);
	}();
	
	return out;
}