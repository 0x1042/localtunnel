package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/exp/slog"
)

var (
	ErrNoAuthRequred = errors.New("server unsupport auth")
	ErrInvalidSecret = errors.New("invalid secret")
)

func relay(c1, c2 io.ReadWriteCloser) (inCount int64, outCount int64, errs []error) {
	var wait sync.WaitGroup

	recordErrs := make([]error, 2)
	pipe := func(number int, from, to io.ReadWriteCloser, count *int64) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("pipe panic.", slog.Any("err", err))
			}
		}()
		defer func(to io.ReadWriteCloser) {
			_ = to.Close()
		}(to)
		defer func(from io.ReadWriteCloser) {
			_ = from.Close()
		}(from)
		defer wait.Done()
		*count, recordErrs[number] = io.Copy(to, from)
	}

	wait.Add(2)
	go pipe(0, c1, c2, &inCount)
	go pipe(1, c2, c1, &outCount)
	wait.Wait()
	for _, err := range recordErrs {
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			continue
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

const (
	stackSize = 4096
)

func Go(name string, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				var buf [stackSize]byte
				n := runtime.Stack(buf[:], false)
				msg := fmt.Sprintf("goroutine-[%s]", name)
				slog.Error(msg, slog.Any("err", err), slog.String("stack", string(buf[:n])))
			}
		}()
		slog.Debug(fmt.Sprintf("goroutine-[%s] start", name))
		f()
	}()
}

func b2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
