package conc

//import "fmt"

/*
	Apply foo to each value in the channel, put on a new channel. Processing
	is done in parallel, though retreival is in sequence. Only length
	foos can be applied at a time.
*/
func MapBuffered(foo func(i Box) Box, in chan Box, length int) chan Box {
	futures := make(chan Thunk, length);
	go func() {
		for i := range in {
			i := i;
			futures <- Future(func() Box { return foo(i) });
		}
		close(futures);
	}();
	return Realize(futures);
}

/*
	Empty in into out, then close out.
*/
func chain(in, out chan Thunk) {
	for b := range in {
		out <- b
	}
	close(out);
}

/*
	Apply foo to each value in the channel, put on a new channel. Processing
	is done in parallel, though retreival is in sequence. All foos can be applied,
	through use of channel chaining. Essentially a linked list with messages being
	passed along.
*/
func Map(foo func(i Box) Box, in chan Box) chan Box {
	futures := make(chan Thunk);
	last := futures;

	go func() {
		for i := range in {
			//keep your own i
			i := i;
			next := make(chan Thunk);
			go chain(next, last);
			last = next;
			last <- Future(func() Box { return foo(i) });
		}
		close(last);
	}();
	return Realize(futures);
}

/*
	Same as above, except no guarantee about the order of things coming out.
	Useful for map-reduce.
*/
func MapUnordered(foo func(i Box) Box, in chan Box) chan Box {
	out := make(chan Box);

	done := make(chan bool);
	count := 0;

	go func() {
		for i := range in {
			count++;
			go func(i Box) {
				out <- foo(i);
				done <- true;
			}(i);
		}
		for i := 0; i < count; i++ {
			<-done
		}
		close(out);
	}();

	return out;
}
