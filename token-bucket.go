package gtask

import "errors"

const LimitSize = 1000000

type TokenBucket struct {
	Ch chan struct{}
}

func NewTokenBucket(size int) *TokenBucket {
	if size > LimitSize {
		panic("[TokenBucket] size too large")
	}
	ch := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		ch <- struct{}{}
	}
	return &TokenBucket{ch}
}

var ErrNoBucket = errors.New("no avaliable bucket")

func (t *TokenBucket) Get(fastFail bool) error {
	if !fastFail {
		<-t.Ch
		return nil
	}
	select {
	case <-t.Ch:
		return nil
	default:
		return ErrNoBucket
	}
}

func (t *TokenBucket) Put() {
	t.Ch <- struct{}{}
}
