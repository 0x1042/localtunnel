package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"

	"github.com/google/uuid"
)

type Authenticator struct {
	hasher hash.Hash
}

func NewAuthenticator(secret string) *Authenticator {
	digest := sha256.New().Sum([]byte(secret))
	return &Authenticator{
		hasher: hmac.New(sha256.New, digest[:]),
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
