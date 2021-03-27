package gtask

import (
	"errors"
	"testing"
)

func TestOnceReset(t *testing.T) {
	var calls int
	var c Once
	c.Do(func() error {
		calls++
		return nil
	})
	c.Do(func() error {
		calls++
		return nil
	})
	c.Do(func() error {
		calls++
		return nil
	})
	if expect := 1; calls != expect {
		t.Fatalf("want %v, got %v", expect, calls)
	}
	c.Reset()
	c.Do(func() error {
		calls++
		return nil
	})
	c.Do(func() error {
		calls++
		return nil
	})
	c.Do(func() error {
		calls++
		return nil
	})
	if expect := 2; calls != expect {
		t.Fatalf("want %v, got %v", expect, calls)
	}
}

func TestOnceError(t *testing.T) {
	var calls int
	var c Once
	c.Do(func() error {
		calls++
		return errors.New("try again")
	})
	c.Do(func() error {
		calls++
		return errors.New("try again")
	})
	c.Do(func() error {
		calls++
		return nil
	})
	c.Do(func() error {
		calls++
		return nil
	})
	c.Do(func() error {
		calls++
		return nil
	})
	if expect := 3; calls != expect {
		t.Fatalf("want %v, got %v", expect, calls)
	}
}
