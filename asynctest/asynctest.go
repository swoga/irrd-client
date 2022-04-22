package asynctest

import (
	"sync"
	"testing"
)

// copied from testing.T without private()
type T interface {
	Cleanup(func())
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Setenv(key, value string)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
}

func New(t *testing.T, fns ...func(T)) {
	wg := sync.WaitGroup{}
	wg.Add(len(fns))

	reCh := make(chan func(*testing.T))
	as := asyncTesting{ch: reCh}

	// execute passed functions and decrement wg after they finish
	for _, f := range fns {
		go func(f func(T)) {
			defer wg.Done()
			f(as)
		}(f)
	}

	// close channel after all functions have returned
	go func() {
		wg.Wait()
		close(reCh)
	}()

	for {
		f, ok := <-reCh
		if !ok {
			break
		}
		// execute testing function on local thread
		f(t)
	}
}

type asyncTesting struct {
	ch chan<- func(*testing.T)
}

// wrapper for all public testing.T functions

func (at asyncTesting) Cleanup(f func()) {
	at.ch <- func(t *testing.T) {
		t.Cleanup(f)
	}
}

func (at asyncTesting) Error(args ...any) {
	at.ch <- func(t *testing.T) {
		t.Error(args...)
	}
}

func (at asyncTesting) Errorf(format string, args ...any) {
	at.ch <- func(t *testing.T) {
		t.Errorf(format, args...)
	}
}

func (at asyncTesting) Fail() {
	at.ch <- func(t *testing.T) {
		t.Fail()
	}
}

func (at asyncTesting) FailNow() {
	at.ch <- func(t *testing.T) {
		t.FailNow()
	}
}

func (at asyncTesting) Failed() bool {
	rch := make(chan bool)
	at.ch <- func(t *testing.T) {
		rch <- t.Failed()
	}
	return <-rch
}

func (at asyncTesting) Fatal(args ...any) {
	at.ch <- func(t *testing.T) {
		t.Fatal(args...)
	}
}

func (at asyncTesting) Fatalf(format string, args ...any) {
	at.ch <- func(t *testing.T) {
		t.Fatalf(format, args...)
	}
}

func (at asyncTesting) Helper() {
	at.ch <- func(t *testing.T) {
		t.Helper()
	}
}

func (at asyncTesting) Log(args ...any) {
	at.ch <- func(t *testing.T) {
		t.Log(args...)
	}
}

func (at asyncTesting) Logf(format string, args ...any) {
	at.ch <- func(t *testing.T) {
		t.Logf(format, args...)
	}
}

func (at asyncTesting) Name() string {
	rch := make(chan string)
	at.ch <- func(t *testing.T) {
		rch <- t.Name()
	}
	return <-rch
}

func (at asyncTesting) Setenv(key, value string) {
	at.ch <- func(t *testing.T) {
		t.Setenv(key, value)
	}
}

func (at asyncTesting) Skip(args ...any) {
	at.ch <- func(t *testing.T) {
		t.Skip(args...)
	}
}

func (at asyncTesting) SkipNow() {
	at.ch <- func(t *testing.T) {
		t.SkipNow()
	}
}

func (at asyncTesting) Skipf(format string, args ...any) {
	at.ch <- func(t *testing.T) {
		t.Skipf(format, args...)
	}
}

func (at asyncTesting) Skipped() bool {
	rch := make(chan bool)
	at.ch <- func(t *testing.T) {
		rch <- t.Skipped()
	}
	return <-rch
}

func (at asyncTesting) TempDir() string {
	rch := make(chan string)
	at.ch <- func(t *testing.T) {
		rch <- t.TempDir()
	}
	return <-rch
}
