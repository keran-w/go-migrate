package phaul

// Remote interface
// Rpc between PhaulClient and PhaulServer. When client
// calls anything on this one, the corresponding method
// should be called on PhaulServer object.

type Remote struct {
}

func (ri *Remote) StartIter() error {
	return nil
}

func (ri *Remote) StopIter() error {
	return nil
}
