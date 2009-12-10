package conc

/*
	Transfers everything in "in" to "out", and when done sends the number of items along
	the returned channel (and closes it).
*/
func Chain(in, out chan interface{}) chan int {
	count := make(chan int);
	go func() {
		c := 0;
		for i := range in {
			c++;
			out <- i;
		}
		count <- c;
		close(count);
	}();
	return count;
}