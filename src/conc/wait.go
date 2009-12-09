package conc

import "runtime"

//A not-so-busy wait
func Wait(foo func() bool) {
	for !foo() {
		runtime.Gosched();
	}
}