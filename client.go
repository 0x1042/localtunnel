package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"golang.org/x/exp/slog"
)

type Client struct {
	tunnelAddr  *net.TCPAddr
	tunnel      *net.TCPConn
	author      *Authenticator
	localPort   int
	mappingPort uint16
}

func NewClient(localPort, mappingPort int, secret, tunnelAddr string) *Client {
	addr, err := net.ResolveTCPAddr("tcp", tunnelAddr)
	if err != nil {
		slog.Error("invalid tunnel address.", slog.Any("err", err))
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		slog.Error("connect to tunnel fail", slog.Any("err", err))
		os.Exit(1)
	}

	client := &Client{
		tunnel:     conn,
		localPort:  localPort,
		tunnelAddr: addr,
	}

	if len(secret) > 0 {
		client.author = NewAuthenticator(secret)
		if err := client.handShake(client.tunnel); err != nil {
			slog.Error("handShake fail.", slog.Any("err", err))
			_ = client.tunnel.Close()
			os.Exit(1)
		}
	}

	slog.Info("client handShake success.")

	bs := helloPacket(uint16(mappingPort))

	if _, err := client.tunnel.Write(bs); err != nil {
		slog.Error("send hello fail.", slog.Any("err", err))
		os.Exit(1)
	}

	var ttype ttype

	if err := binary.Read(client.tunnel, binary.BigEndian, &ttype); err != nil {
		slog.Error("read hello fail.", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("read message type.", slog.Any("ttype", ttype))

	switch ttype {
	case fail:
		msg := parseFailPacket(client.tunnel)
		slog.Error("receive error.", slog.String("err", msg))
		os.Exit(1)
	case hello:
		rport, err := parseHelloPacket(client.tunnel)
		if err != nil {
			slog.Error("read hello fail.", slog.Any("err", err))
			os.Exit(1)
		}
		client.mappingPort = rport

		mappingAddr := fmt.Sprintf("%s:%d", addr.IP.String(), rport)
		slog.Info("forward info", slog.String("to", mappingAddr))
	}

	return client

}

func (c *Client) Start() error {
	for {
		var ttype ttype
		if err := binary.Read(c.tunnel, binary.BigEndian, &ttype); err != nil {
			slog.Error("read fail.", slog.Any("err", err))
			os.Exit(1)
		}

		slog.Info("ttype ", slog.Any("ttype", ttype))
		switch ttype {
		case connect:
			if err := c.handleConnect(); err != nil {
				slog.Error("connect fail.", slog.Any("err", err))
				os.Exit(1)
			}
		case fail:
			msg := parseFailPacket(c.tunnel)
			slog.Error("receive error.", slog.String("err", msg))
			os.Exit(1)
		}

	}
}

func (c *Client) handShake(remote net.Conn) error {
	uid, err := parseChallengePacket(remote)
	if err != nil {
		return err
	}

	slog.Info("client handShake.", slog.Any("id", uid.String()))

	sign := c.author.Sign(uid)

	slog.Debug("client handShake. sign ", slog.Any("sign", sign))

	bs := authPacket(sign)

	if _, err := remote.Write(bs); err != nil {
		return err
	}
	return nil
}

func (c *Client) handleConnect() error {
	uid, err := parseTransferOrConnectPacket(c.tunnel)
	if err != nil {
		return err
	}
	remote, err := net.DialTCP("tcp", nil, c.tunnelAddr)
	if err != nil {
		return err
	}

	if c.author != nil {
		slog.Debug("client handShake")
		if err := c.handShake(remote); err != nil {
			slog.Error("handShake fail.", slog.Any("err", err))
			os.Exit(1)
		}
	}
	slog.Debug("client handShake success")

	packet := transferOrConnectPacket(uid, transfer)

	if _, err := remote.Write(packet); err != nil {
		slog.Error("Write fail.", slog.Any("err", err))
		os.Exit(1)
	}

	local, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", c.localPort))
	if err != nil {
		slog.Error("dial localfail.", slog.Any("err", err))
		os.Exit(1)
	}

	inCount, outCount, _ := relay(local, remote)

	slog.Info("forward success.", slog.Int64("read_size", inCount),
		slog.Int64("write_size", outCount))

	return nil
}

func StartClient(localPort, mappingPort int, secret, tunnelAddr string) error {
	client := NewClient(localPort, mappingPort, secret, tunnelAddr)
	client.Start()
	return nil
}
