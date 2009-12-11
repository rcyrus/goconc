package conc

import "runtime"

/*
	Returns a stream of channels that will grab values from the input, with
	the idea that you can use range and <-ch,closed(ch) in a threadsafe fashion.
*/
func SafeChan(inputs <-chan Box) chan chan Box {
	loop := make(chan chan Box);
	outgoing := make(chan chan Box);
	
	go func() {
		for {
			nextChan := make(chan Box);
			outgoing <- nextChan;
			//once the new channel is grabbed by someone, add it to the loop
			loop <- nextChan;
		}
	}();
	
	go func() {
		for v := range inputs {
			for receiver := range loop {
				go func(receiver chan Box) {
					loop <- receiver;
				}(receiver);
				if receiver <- v {
					break;
				}
				runtime.Gosched();
			}
		}
		
		for receiver := range loop {
			close(receiver);
		}
	}();
	
	return outgoing;
}