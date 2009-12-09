package conc

import (
	"testing";
//	"fmt";
)

func TestMap(t *testing.T) {
	incr := func(a Box) Box {
		return a.(int)+1;
	};
	numbers := make(chan Box);
	go func() {
		for i:=0; i<10; i++ {
			numbers <- i;
		}
		close(numbers);
	}();
	incrNumbers := Map(incr, numbers);
	for i:=0; i<10; i++ {
		j := <- incrNumbers;
		if i+1 != j.(int) {
			t.Fail();
		}
	}
}

func TestFold(t *testing.T) {
	sum := func(a Box, b Box) Box {
		return a.(int)+b.(int);
	};
	numbers := make(chan Box);
	go func() {
		for i:=0; i<10; i++ {
			numbers <- i;
		}
		close(numbers);
	}();
	totalSum := Fold(sum, numbers);

	if totalSum.(int) != 45 {
		t.Fail();
	}
}