package main

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

type Authenticator struct {
	hasher hash.Hash
}

func NewAuthenticator(secret string) *Authenticator {
	digest := sha3.Sum256([]byte(secret))
	return &Authenticator{
		hasher: hmac.New(sha3.New256, digest[:]),
	}
}

func (au *Authenticator) Sign(id uuid.UUID) string {
	au.hasher.Reset()
	au.hasher.Write(id[:])
	return hex.EncodeToString(au.hasher.Sum(nil))
}

func (au *Authenticator) Verify(id uuid.UUID, tag string) bool {
	sign := au.Sign(id)
	return hmac.Equal(s2b(tag), s2b(sign))
}
