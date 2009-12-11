package conc

import "sync"

type MultiReader struct {
	ch <-chan Box;
	m sync.Mutex;
}

func (cw *MultiReader) read() (v Box, ok bool) {
	cw.m.Lock();
	v, ok = <-cw.ch, closed(cw.ch);
	cw.m.Unlock();
	return;
}
