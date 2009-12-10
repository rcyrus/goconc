package conc

/*
	Take a channel of thunks and put them into a channel of values.
*/
func Realize(thunks chan Thunk) chan Box {
	values := make(chan Box);
	
	go func() {
		for thunk := range thunks {
			values <- thunk();
		}
		close(values);
	}();
	
	return values;
}