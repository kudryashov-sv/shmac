package mocktrans

import "fmt"

type ioError struct {
	msg       string
	isTemp    bool
	isTimeout bool
}

func NewError(msg string, temp, timeout bool) *ioError {
	return &ioError{
		msg:       msg,
		isTemp:    temp,
		isTimeout: timeout,
	}
}

func (ie *ioError) Error() string {
	return fmt.Sprintf("%s temp=%v timeout=%v", ie.msg, ie.isTemp, ie.isTimeout)
}

func (ie *ioError) Timeout() bool {
	return ie.isTimeout
}

func (ie *ioError) Temporary() bool {
	return ie.isTemp
}
