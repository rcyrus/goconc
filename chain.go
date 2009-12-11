package conc

/*
	Transfers everything in "in" to "out", and when done sends the number of items along
	the returned channel.
*/
func Chain(in, out chan Box) chan int {
	count := make(chan int);
	go func() {
		c := 0;
		for i := range in {
			c++;
			out <- i;
		}
		for {
			//send it forever in case multiple sites are interested
			count <- c;
		}
	}();
	return count;
}