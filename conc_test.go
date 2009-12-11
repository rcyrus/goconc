package conc

import (
	"testing";
	"time";
)

func slowNumbers(max int, delay int64) chan Box {
	out := make(chan Box);
	go func() {
		for i:=0; i<max; i++ {
			out <- i;
			time.Sleep(delay);
		}
		close(out);
	}();
	return out;
}

func TestFor(t *testing.T) {
	var vals [20]int;
	wait := ForChunk(CountStream(0, 20), func(i Box) {vals[i.(int)] = 1}, 3);
	wait();
	
	total := 0;
	for i:=0; i<20; i++ {
		total += vals[i];
	}
	
	if total != 20 {
		t.Fail();
	}
	
	wait = For(CountStream(0, 20), func(i Box) {vals[i.(int)] = 1});
	wait();
	
	total = 0;
	for i:=0; i<20; i++ {
		total += vals[i];
	}
}

func TestMap(t *testing.T) {
	incr := func(a Box) Box {
		return a.(int)+1;
	};
	incrNumbers := Map(incr, CountStream(0, 20));
	for i:=0; i<20; i++ {
		j := <- incrNumbers;
		if i+1 != j.(int) {
			t.Fail();
		}
	}
}

func TestReduce(t *testing.T) {
	sum := func(a Box, b Box) Box {
		return a.(int)+b.(int);
	};
	totalSum := Reduce(sum, CountStream(0, 10), 0);

	if totalSum.(int) != 45 {
		t.Fail();
	}
}

func TestFilter(t *testing.T) {	
	results := Filter(func(i Box) bool { return i.(int)%2==0 }, CountStream(0, 10));
	trueRes := make([]bool, 10);
	for i := range results {
		trueRes[i.(int)] = true;
	}
	for i,v := range trueRes {
		if v != (i%2==0) {
			t.Fail();
		}
	}
}

func TestMapReduce(t *testing.T) {
	
	incr := func(a Box) Box {
		return a.(int)+1;
	};
	sum := func(a Box, b Box) Box {
		return a.(int)+b.(int);
	};
	
	result := Reduce(sum, MapUnordered(incr, CountStream(0, 10)), 0);

	if result.(int) != 55 {
		t.Fail();
	}
}

func TestSafeChan(t *testing.T) {
	slowCounter := func () (foo func() Box) {
		i := 0;
		return func() Box {
			defer func() {i++}();
			time.Sleep(1e8);
			return i;
		};
	};
	baseChan := EvalStream(slowCounter());
	
	chan1 := SafeChan(baseChan);
	chan2 := SafeChan(baseChan);
	chan3 := SafeChan(baseChan);
	
	collector := make(chan int);
	
	go func() {
		for v := range chan1 {
			collector <- v.(int);
		}
	}();
	go func() {
		for v := range chan2 {
			collector <- v.(int);
		}
	}();
	go func() {
		for v := range chan3 {
			collector <- v.(int);
		}
	}();
	
	for i:=0; i<10; i++ {
		<-collector;
	}
}
