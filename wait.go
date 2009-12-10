package conc

import "runtime"

//A not-so-busy wait
func WaitUntil(foo func() bool) {
	for !foo() {
		runtime.Gosched();
	}
}
