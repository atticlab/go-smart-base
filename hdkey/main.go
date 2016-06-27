package hdkey

import (
	"encoding/hex"
	"github.com/stellar/go-stellar-base/strkey"

)

func GetMnemonic() string {
	seed, _ := genSeed(32)
	str := hex.EncodeToString(seed)
	return encodeMnemonic(str)
}

func GetMasterPriv(str string) string {
	decoded := decodeMnemonic(str)
	mkp := masterKey(decoded, true)
	key, _ := strkey.Encode(strkey.VersionByteMPriv, mkp.PrivateKey)
	return key
}

func GetMasterPub(str string) string {
	decoded := decodeMnemonic(str)
	mkp := masterKey(decoded, false)
	key, _ := strkey.Encode(strkey.VersionByteMPub, mkp.PublicKey)
	return key
}
