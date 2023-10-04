package main

import (
	"bytes"
	"encoding/binary"
	"log/slog"
	"net"

	"github.com/google/uuid"
)

type (
	ttype uint8
)

func (t ttype) String() string {
	switch t {
	case challenge:
		return "challenge"
	case auth:
		return "auth"
	case hello:
		return "hello"
	case connect:
		return "connect"
	case transfer:
		return "transfer"
	case fail:
		return "fail"
	default:
		return "unknown"
	}
}

var emptyUUID = [16]byte{}

const (
	challenge ttype = iota + 1
	auth
	hello
	connect
	transfer
	fail
)

func challengePacket(uuid uuid.UUID) []byte {
	buf := bytes.Buffer{}
	_ = binary.Write(&buf, binary.BigEndian, challenge)
	_ = binary.Write(&buf, binary.BigEndian, uuid[:])
	return buf.Bytes()
}

func parseChallengePacket(conn net.Conn) (uuid.UUID, error) {
	var ttype ttype
	if err := binary.Read(conn, binary.BigEndian, &ttype); err != nil {
		return emptyUUID, err
	}

	id := make([]byte, 16)
	if err := binary.Read(conn, binary.BigEndian, &id); err != nil {
		return emptyUUID, err
	}
	uid, err := uuid.FromBytes(id)
	if err != nil {
		return emptyUUID, err
	}
	return uid, nil
}

func authPacket(tag string) []byte {
	buf := bytes.Buffer{}
	_ = binary.Write(&buf, binary.BigEndian, auth)
	_ = binary.Write(&buf, binary.BigEndian, uint8(len(tag)))
	_ = binary.Write(&buf, binary.BigEndian, []byte(tag))
	return buf.Bytes()
}

func parseFailPacket(conn net.Conn) string {
	var size uint32
	if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
		slog.Error("parseFailPacket err.", slog.Any("err", err))
		return ""
	}
	data := make([]byte, size)
	if err := binary.Read(conn, binary.BigEndian, &data); err != nil {
		slog.Error("parseFailPacket err.", slog.Any("err", err))
		return ""
	}
	return string(data)
}

func failPacket(err error) []byte {
	buf := bytes.Buffer{}
	_ = binary.Write(&buf, binary.BigEndian, fail)
	_ = binary.Write(&buf, binary.BigEndian, uint32(len(err.Error())))
	_ = binary.Write(&buf, binary.BigEndian, []byte(err.Error()))
	return buf.Bytes()
}

func parseAuthPacket(conn net.Conn) (string, error) {
	var ttype byte
	if err := binary.Read(conn, binary.BigEndian, &ttype); err != nil {
		return "", err
	}

	var size uint8
	if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
		return "", err
	}

	buf := make([]byte, size)
	if err := binary.Read(conn, binary.BigEndian, &buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

func parseHelloPacket(conn net.Conn) (uint16, error) {
	var port uint16
	if err := binary.Read(conn, binary.BigEndian, &port); err != nil {
		return 0, err
	}
	return port, nil
}

func helloPacket(port uint16) []byte {
	buf := bytes.Buffer{}
	_ = binary.Write(&buf, binary.BigEndian, hello)
	_ = binary.Write(&buf, binary.BigEndian, port)
	return buf.Bytes()
}

func transferOrConnectPacket(id uuid.UUID, mtype ttype) []byte {
	buf := bytes.Buffer{}
	_ = binary.Write(&buf, binary.BigEndian, mtype)
	_ = binary.Write(&buf, binary.BigEndian, id[:])
	return buf.Bytes()
}

func parseTransferOrConnectPacket(conn net.Conn) (uuid.UUID, error) {
	id := make([]byte, 16)
	if err := binary.Read(conn, binary.BigEndian, &id); err != nil {
		return emptyUUID, err
	}
	uid, err := uuid.FromBytes(id)
	if err != nil {
		return emptyUUID, err
	}
	return uid, nil
}
