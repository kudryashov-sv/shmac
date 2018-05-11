package mocktrans

type fakeAddr struct {
}

func newFakeAddr() fakeAddr {
	return fakeAddr{}
}

func (fakeAddr) Network() string {
	return "fake"
}

func (fakeAddr) String() string {
	return "fake:address"
}
