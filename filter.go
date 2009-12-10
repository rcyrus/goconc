package conc

/*
	Things in in get to out if they pass foo(). Order on out is arbitrary.
*/
func Filter(foo func(i Box) bool, in chan Box) chan Box {
	out := make(chan Box);
	
	done := make(chan bool);
	count := 0;
	
	go func() {
		for i := range in {
			count++;
			i := i;
			go func() {
				if foo(i) {
					out <- i;
				}
				done <- true;
			}();
		}
	}();
	
	go func() {
		WaitUntil(func() bool {return closed(in)});
		for i := 0; i<count; i++ {
			<-done;
		}
		close(out);
	}();
	
	return out;
}