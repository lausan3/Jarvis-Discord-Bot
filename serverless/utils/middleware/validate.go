package middleware

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/sirupsen/logrus"
)

func ValidateDiscordSecurityHeaders(body string, signature string, timestamp string, pubkey string) bool {
	pubkeyBytes, err := hex.DecodeString(pubkey)
	if err != nil || len(pubkeyBytes) != ed25519.PublicKeySize {
		logrus.Error("Invalid public key format")
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil || len(sig) != ed25519.SignatureSize {
		logrus.Error("Invalid signature format")
		return false
	}

	message := append([]byte(timestamp), []byte(body)...)

	return ed25519.Verify(pubkeyBytes, message, sig)
}
