package main

import (
	"errors"
	"io"
	"net"
	"runtime"
	"sync"
	"unsafe"

	"github.com/rs/zerolog/log"
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
				log.Error().Any("err", err).Msg("pipe panic.")
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
		name := "routine-[" + name + "]"
		defer func() {
			if err := recover(); err != nil {
				var buf [stackSize]byte
				n := runtime.Stack(buf[:], false)
				log.Error().Str("job", name).Any("err", err).Str("stack", string(buf[:n])).Msg("routine panic.")
			}
		}()
		log.Info().Str("job", name).Msg("routine start.")
		f()
	}()
}

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
