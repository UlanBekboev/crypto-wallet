package utils

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

// VerifySignature проверяет, что подпись соответствует адресу
func VerifySignature(address, signature, message string) error {
	// Добавим префикс как делает это MetaMask (EIP-191)
	prefixedMsg := "\x19Ethereum Signed Message:\n" + string(len(message)) + message
	msgHash := crypto.Keccak256Hash([]byte(prefixedMsg))

	// Декодируем подпись
	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return errors.New("невалидная подпись")
	}
	if sig[64] != 27 && sig[64] != 28 {
		return errors.New("v должен быть 27 или 28")
	}
	sig[64] -= 27 // Go-ethereum требует 0 или 1

	pubKey, err := crypto.SigToPub(msgHash.Bytes(), sig)
	if err != nil {
		return errors.New("не удалось извлечь публичный ключ")
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()
	if strings.ToLower(recoveredAddr) != strings.ToLower(address) {
		return errors.New("адрес не соответствует подписи")
	}

	return nil
}
