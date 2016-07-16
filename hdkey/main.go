package hdkey

import (
	"encoding/hex"
	"bitbucket.org/atticlab/go-stellar-base/strkey"
)

func GetMnemonic() string {
	seed, _ := GenSeed(32)
	str := hex.EncodeToString(seed)
	return encodeMnemonic(str)
}

func GetMnemonicFromSeed(seed []byte) string{
	str := hex.EncodeToString(seed)
	return encodeMnemonic(str)
}

func GetSeedFromMnemonic(str string) string {
	decoded := decodeMnemonic(str)
	return decoded
}
//TODO: Fix encoding to Base32 (last 4 char is "====")
func GetMasterPriv(str string) string {
	decoded, _ := hex.DecodeString(decodeMnemonic(str))
	mkp := MasterKey(decoded, Private)
	key, _ := strkey.Encode(strkey.VersionByteMPriv, append(mkp.Chaincode, mkp.PrivateKey...))
	return key//[:len(key)-4]
}

func GetMasterPub(str string) string {
	decoded, _ := hex.DecodeString(decodeMnemonic(str))
	mkp := MasterKey(decoded, Public)
	key, _ := strkey.Encode(strkey.VersionByteMPub, append(mkp.Chaincode, mkp.PublicKey...))
	return key
}
