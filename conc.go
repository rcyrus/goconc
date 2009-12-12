package conc

type Box interface{}
type Thunk func() Box
type ThunkChan chan Box
