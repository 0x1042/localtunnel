package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Client struct {
	addr       *net.TCPAddr
	tunnel     *net.TCPConn
	author     *Authenticator
	localPort  int
	targetPort uint16
}

func NewClient(localPort, mappingPort int, secret, tunnelAddr string) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", tunnelAddr)
	if err != nil {
		log.Error().Err(err).Msg("invalid tunnel address.")
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Error().Err(err).Msg("connect to tunnel fail")
		return nil, err
	}

	client := &Client{
		tunnel:    conn,
		localPort: localPort,
		addr:      addr,
	}

	if len(secret) > 0 {
		client.author = NewAuthenticator(secret)
		if err := client.handShake(client.tunnel); err != nil {
			log.Error().Err(err).Msg("handShake fail.")
			_ = client.tunnel.Close()
			return nil, err
		}
	}

	log.Info().Msg("client handShake success.")

	bs := helloPacket(uint16(mappingPort))

	if _, err := client.tunnel.Write(bs); err != nil {
		log.Error().Err(err).Msg("send hello fail.")
		return nil, err
	}

	var ttype ttype

	if err := binary.Read(client.tunnel, binary.BigEndian, &ttype); err != nil {
		log.Error().Err(err).Msg("read hello fail.")
		return nil, err
	}

	log.Info().Uint8("type", uint8(ttype)).Msg("read message type.")

	switch ttype {
	case fail:
		msg := parseFailPacket(client.tunnel)
		log.Error().Str("err", msg).Msg("receive error.")
		return nil, errors.New("receive error")
	case hello:
		rport, err := parseHelloPacket(client.tunnel)
		if err != nil {
			log.Error().Err(err).Msg("read hello fail.")
			return nil, err
		}
		client.targetPort = rport

		mappingAddr := addr.IP.String() + ":" + strconv.FormatUint(uint64(rport), 10)
		log.Info().Str("to", mappingAddr).Msg("forward info")
	default:
		panic("unhandled default case")
	}

	return client, nil
}

func (c *Client) Start() error {
	for {
		var ttype ttype
		if err := binary.Read(c.tunnel, binary.BigEndian, &ttype); err != nil {
			if errors.Is(err, io.EOF) {
				log.Warn().Msg("server shutdown.")
			} else {
				log.Error().Err(err).Msg("read fail.")
			}
			os.Exit(1)
		}

		log.Info().Uint8("type", uint8(ttype)).Msg("receive request")

		switch ttype {
		case connect:
			if err := c.handleConnect(); err != nil {
				log.Error().Err(err).Msg("connect fail.")
				os.Exit(1)
			}
		case fail:
			msg := parseFailPacket(c.tunnel)
			log.Error().Str("err", msg).Msg("receive error.")
			os.Exit(1)
		default:
			panic("unhandled default case")
		}

	}
}

func (c *Client) handShake(remote net.Conn) error {
	uid, err := parseChallengePacket(remote)
	if err != nil {
		return err
	}

	log.Info().Str("id", uid.String()).Msg("client handShake.")

	sign := c.author.Sign(uid)

	log.Debug().Str("sign", sign).Msg("client handShake. sign ")

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
	remote, err := net.DialTCP("tcp", nil, c.addr)
	if err != nil {
		return err
	}

	if c.author != nil {
		log.Debug().Msg("client handShake")
		if err := c.handShake(remote); err != nil {
			log.Error().Err(err).Msg("handShake fail.")
			os.Exit(1)
		}
	}
	log.Debug().Msg("client handShake success")

	packet := transferOrConnectPacket(uid, transfer)

	if _, err := remote.Write(packet); err != nil {
		log.Error().Err(err).Msg("Write fail.")
		os.Exit(1)
	}

	addr := "localhost:" + strconv.FormatInt(int64(c.localPort), 10)
	local, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msg("dial localfail.")
		os.Exit(1)
	}

	inCount, outCount, _ := relay(local, remote)

	log.Info().Int64("in_size", inCount).Int64("out_size", outCount).Msg("forward success")
	return nil
}

func StartClient(localPort, targetPort int, secret, addr string) error {
	client, err := NewClient(localPort, targetPort, secret, addr)
	if err != nil {
		return err
	}
	log.Info().Msg("new client success.")
	return client.Start()
}
