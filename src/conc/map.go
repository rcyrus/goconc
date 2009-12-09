package conc

func Map(in chan Box, foo func(i Box) Box) chan Box {
	out := make(chan Box);
	futures := make(chan Thunk);
	for i := range in {
		futures <- Future(func() Box {return foo(i)});
	}
	go func() {
		for f := range futures {
			out <- f();
		}
		close(out);
	}();
	return out;
}
