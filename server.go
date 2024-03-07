package main

import (
	"encoding/binary"
	"errors"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("server listen error")
		os.Exit(1)
	}

	log.Info().Str("addr", ln.Addr().String()).Msg("server listen success.")

	for {
		stream, err := ln.AcceptTCP()
		if err != nil {
			log.Error().Err(err).Msg("accept fail.")
			break
		}

		remote := stream.RemoteAddr().String()
		log.Info().Str("remote", remote).Msg("incoming connection.")
		Go("handle", func() {
			defer func() {
				if stream != nil {
					_ = stream.Close()
				}
				log.Info().Str("remote", remote).Msg("client exit.")
			}()
			svr.handle(stream)
		})
	}

	return nil
}

func (svr *Server) handle(stream *net.TCPConn) {
	// 鉴权
	if svr.author != nil {
		log.Debug().Msg("handshake start")
		if err := svr.handshake(stream); err != nil {
			log.Error().Err(err).Msg("handshake fail")
			packet := failPacket(err)
			_, _ = stream.Write(packet)
			return
		}
	}

	log.Debug().Msg("handshake success")

	var newType ttype

	if err := binary.Read(stream, binary.BigEndian, &newType); err != nil {
		log.Error().Err(err).Msg("read fail")
		return
	}

	log.Debug().Uint8("ttype", uint8(newType)).Msg("message info")

	switch newType {

	case hello:
		if err := svr.handleHello(stream); err != nil {
			log.Error().Err(err).Msg("client handShake error")
		}

	case transfer:
		if err := svr.handleTransfer(stream); err != nil {
			log.Error().Err(err).Msg("client transfer error")
		}
	case auth:
		log.Error().Msg("unexpect auth message")
		packet := failPacket(ErrNoAuthRequred)
		_, _ = stream.Write(packet)
		return
	default:
		panic("unhandled default case")
	}
}

func (svr *Server) handshake(stream *net.TCPConn) error {
	id := uuid.New()

	log.Debug().Str("id", id.String()).Msg("handshake start,write to client")

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
		log.Error().Err(err).Msg("read hello message error")
		return err
	}

	addr := &net.TCPAddr{
		IP:   nil,
		Port: int(port),
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error().Err(err).Str("addr", addr.String()).Msg("listen error")
		return err
	}

	realPort := ln.Addr().(*net.TCPAddr).Port
	packet := helloPacket(uint16(realPort))

	log.Info().Int("port", realPort).Str("remote", stream.RemoteAddr().String()).Msg("mapping port.")

	if _, err := stream.Write(packet); err != nil {
		log.Error().Err(err).Msg("write hello message error")
		return err
	}

	for {
		incoming, err := ln.AcceptTCP()
		if err != nil {
			log.Error().Err(err).Msg("accept error")
			break
		}

		log.Debug().Msg("incoming request")

		id := uuid.New()
		svr.connStore.Store(id.String(), incoming)

		packet := transferOrConnectPacket(id, connect)

		if _, err := stream.Write(packet); err != nil {
			log.Error().Err(err).Msg("write error")
		}
	}
	return nil
}

func (svr *Server) handleTransfer(stream *net.TCPConn) error {
	log.Debug().Msg("handle forward")
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
		log.Error().Err(err).Msg("relay fail")
		return nil
	}

	log.Info().Int64("in_size", inCount).Int64("out_size", outCount).Msg("forward success.")

	return nil
}

func StartServer(port int, secret string) error {
	server := NewServer(port, secret)
	return server.Start()
}
