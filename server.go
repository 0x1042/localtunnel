package main

import (
	"encoding/binary"
	"errors"
	"log/slog"
	"net"
	"os"

	"github.com/google/uuid"
)

type (
	db = Map[string, net.Conn]
)

type Server struct {
	addr      *net.TCPAddr
	connStore *db
	author    *Authenticator
}

func NewServer(port int, secret string) *Server {
	addr := &net.TCPAddr{
		IP:   nil,
		Port: port,
	}
	srv := &Server{
		addr:      addr,
		connStore: new(db),
	}
	if len(secret) > 0 {
		srv.author = NewAuthenticator(secret)
	}
	return srv
}

func (svr *Server) Start() error {
	ln, err := net.ListenTCP("tcp", svr.addr)
	if err != nil {
		slog.Error("server listen error", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("server listen success.", slog.Any("addr", ln.Addr().String()))

	for {
		stream, err := ln.AcceptTCP()
		if err != nil {
			slog.Error("accept fail.", slog.Any("err", err))
			break
		}

		remote := stream.RemoteAddr().String()
		slog.Info("incoming connection.", slog.String("remote", remote))
		Go("handle", func() {
			defer func() {
				if stream != nil {
					_ = stream.Close()
				}
				slog.Info("client exit", slog.String("remote", remote))
			}()
			svr.handle(stream)
		})
	}

	return nil
}

func (svr *Server) handle(stream *net.TCPConn) {
	// 鉴权
	if svr.author != nil {
		slog.Debug("handshake start")
		if err := svr.handshake(stream); err != nil {
			slog.Error("handshake fail", slog.Any("err", err))
			packet := failPacket(err)
			_, _ = stream.Write(packet)
			return
		}
	}

	slog.Debug("handshake success")

	var newType ttype

	if err := binary.Read(stream, binary.BigEndian, &newType); err != nil {
		slog.Error("read fail", slog.Any("err", err))
		return
	}

	slog.Debug("message type", slog.Any("ttype", newType))

	switch newType {

	case hello:
		if err := svr.handleHello(stream); err != nil {
			slog.Error("client handShake error", slog.Any("err", err))
		}

	case transfer:
		if err := svr.handleTransfer(stream); err != nil {
			slog.Error("client transfer error", slog.Any("err", err))
		}
	case auth:
		slog.Error("unexpect auth message")
		packet := failPacket(ErrNoAuthRequred)
		_, _ = stream.Write(packet)
		return
	}
}

func (svr *Server) handshake(stream *net.TCPConn) error {
	id := uuid.New()

	slog.Debug("handshake start,write to client", slog.String("id", id.String()))

	// 1. send
	if _, err := stream.Write(challengePacket(id)); err != nil {
		return err
	}

	// 2. wait client
	tag, err := parseAuthPacket(stream)
	if err != nil {
		return err
	}

	// 3. auth
	if verify := svr.author.Verify(id, tag); !verify {
		return ErrInvalidSecret
	}
	return nil
}

func (svr *Server) handleHello(stream *net.TCPConn) error {
	port, err := parseHelloPacket(stream)
	if err != nil {
		slog.Error("read hello message error", slog.Any("err", err))
		return err
	}

	addr := &net.TCPAddr{
		IP:   nil,
		Port: int(port),
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		slog.Error("listen error", slog.Any("err", err), slog.String("addr", addr.String()))
		return err
	}

	realPort := ln.Addr().(*net.TCPAddr).Port
	packet := helloPacket(uint16(realPort))
	slog.Info("mapping port.", slog.Any("port", realPort))
	if _, err := stream.Write(packet); err != nil {
		slog.Error("write hello message error", slog.Any("err", err))
		return err
	}

	for {
		incoming, err := ln.AcceptTCP()
		if err != nil {
			slog.Error("accept error", slog.Any("err", err))
			break
		}

		slog.Debug("incoming request")

		id := uuid.New()
		svr.connStore.Store(id.String(), incoming)

		packet := transferOrConnectPacket(id, connect)

		if _, err := stream.Write(packet); err != nil {
			slog.Error("write error", slog.Any("err", err))
		}
	}
	return nil
}

func (svr *Server) handleTransfer(stream *net.TCPConn) error {
	slog.Debug("handle forward")
	uid, err := parseTransferOrConnectPacket(stream)
	if err != nil {
		return err
	}

	session, ok := svr.connStore.Load(uid.String())

	if !ok {
		return errors.New("not exist:" + uid.String())
	}

	inCount, outCount, errs := relay(session, stream)

	if errs != nil {
		slog.Error("relay fail", slog.Any("err", errs))
		return nil
	}
	slog.Info("forward success.", slog.Int64("read_size", inCount), slog.Int64("write_size", outCount))
	return nil
}

func StartServer(port int, secret string) error {
	server := NewServer(port, secret)
	return server.Start()
}
