package conc

//import "fmt"

/*
	Apply foo to each value in the channel, put on a new channel. Processing
	is done in parallel, though retreival is in sequence.
*/
func Map(foo func(i Box) Box, in chan Box, length int) chan Box {
	futures := make(chan Thunk, length);
	go func() {
		for i := range in {
			i := i;
			futures <- Future(func() Box {return foo(i)});
		}
		close(futures);
	}();
	return Realize(futures);
}
