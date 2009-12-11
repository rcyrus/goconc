package conc

func NaturalNumbers() (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		for i:=1;; i++ {
			out <- i;
		}
	}(out);
	return;
}

func WholeNumbers() (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		for i:=0;; i++ {
			out <- i;
		}
	}(out);
	return;
}

func CountStream(start, cap Box) (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		switch start.(type) {
		case int:
			start := start.(int);
			cap := cap.(int);
			for i:=start; i<cap; i++ {
				out <- i;
			}
		case float:
			start := start.(float);
			cap := cap.(float);
			for i:=start; i<cap; i++ {
				out <- i;
			}
		}
		close(out);
	}(out);
	return;
}

func IncrementStream(start, cap, incr Box) (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		switch start.(type) {
		case int:
			start := start.(int);
			cap := cap.(int);
			incr := incr.(int);
			for i:=start; i<cap; i+=incr {
				out <- i;
			}
		case float:
			start := start.(float);
			cap := cap.(float);
			incr := incr.(float);
			for i:=start; i<cap; i+=incr {
				out <- i;
			}
		}
		close(out);
	}(out);
	return;
}

func OnceStream(val Box) (out chan Box) {
	out = make(chan Box, 1);
	out <- val;
	close(out);
	return out;
}

func RepeatStream(val Box, n int) (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		for i:=0; i<n; i++ {
			out <- val;
		}
		close(out);
	}(out);
	return;
}

func RepeatForeverStream(val Box) (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		for {
			out <- val;
		}
	}(out);
	return;
}

func EvalStream(foo func() Box) (out chan Box) {
	out = make(chan Box);
	go func(out chan Box) {
		for {
			out <- foo();
		}
	}(out);
	return;
}
