package hdkey

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	//"github.com/stellar/go-stellar-base/strkey"
	"strconv"
	"strings"
)

var (
	Public  []byte
	Private []byte
)

func init() {
	Public, _ = hex.DecodeString("0488B21E")
	Private, _ = hex.DecodeString("0488ADE4")
}

// HDKey defines the components of a hierarchical deterministic wallet
type HDKey struct {
	Vbytes      []byte //4 bytes
	Depth       uint16 //1 byte
	Fingerprint []byte //4 bytes
	I           []byte //4 bytes
	Chaincode   []byte //32 bytes
	PrivateKey  []byte //32 bytes
	PublicKey   []byte //32 bytes
}

//Derive returns the ith child of wallet w. Values of path = "m/i/j/../x"
//or path = "M/i/j/../x". Func split path, run Child with i = i||j||..||xs
func (w *HDKey) Derive(p string) (*HDKey, error) {
	var path []string
	var keyPath bool
	path = strings.Split(p, "/")
	switch path[0] {
	case "m":
		if bytes.Compare(w.Vbytes, Private) == 0 {
			keyPath = true
		} else {
			return &HDKey{}, errors.New("Invalid Path")
		}
	case "M":
		keyPath = false
	}
	wTemp := w
	for i := 1; i < len(path[:]); i++ {
		ind64, _ := strconv.ParseUint(path[i], 10, 32)
		ind := uint32(ind64)
		wTemp, _ = wTemp.Child(ind, keyPath)
	}

	return &HDKey{wTemp.Vbytes, wTemp.Depth, wTemp.Fingerprint, wTemp.I,
			wTemp.Chaincode, wTemp.PrivateKey, wTemp.PublicKey}, nil
}

// Child returns the ith child of wallet w. Values of i >= 2^31
// signify private key derivation. Attempting private key derivation
// with a public key will throw an error.
func (w *HDKey) Child(i uint32, f bool) (*HDKey, error) {
	var fingerprint, I, newPriv, newPub []byte
	data := append(w.PublicKey, uint32ToByte(i)...)
	switch f {
	case true:
		mac := hmac.New(sha512.New, w.Chaincode)
		mac.Write(data)
		I = mac.Sum(nil)
		iL := modN(I[:32])
		newPriv = AddPrivKeys(iL, w.PrivateKey)
		newPub = PrivToPub(newPriv)
		fingerprint = hash(w.PublicKey)[:4]

	case false:
		mac := hmac.New(sha512.New, w.Chaincode)
		mac.Write(data)
		I = mac.Sum(nil)
		pub := PrivToPub(modN(I[:32]))
		newPriv = nil
		newPub = AddPubKeys(w.PublicKey, pub)
		fingerprint = hash(w.PublicKey)[:4]
	}
	return &HDKey{w.Vbytes, w.Depth + 1, fingerprint, uint32ToByte(i), I[32:], newPriv, newPub}, nil
}

func GenSeed(length int) ([]byte, error) {
	b := make([]byte, length)
	if length < 128 {
		return b, errors.New("length must be at least 128 bits")
	}
	_, err := rand.Read(b)
	return b, err
}

func MasterKeyStr(data string) *HDKey {
	seed, _ := verDecodeCheck(data)
	return MasterKey(seed)
}

// MasterKey returns a new wallet given a random seed.
func MasterKey(seed []byte) *HDKey {
	key := []byte("Stellar seed")
	mac := hmac.New(sha512.New, key)
	mac.Write(seed)
	I := mac.Sum(nil)
	secret := modN(I[:len(I)/2])
	chain_code := I[len(I)/2:]
	depth := 0
	i := make([]byte, 4)
	fingerprint := make([]byte, 4)

	return &HDKey{Private, uint16(depth), fingerprint, i, chain_code, secret, PrivToPub(secret)}
}


// Serialize returns the serialized form of the wallet.
//func (w *HDKey) Serialize() []byte {
//	depth := uint16ToByte(uint16(w.Depth % 256))
//	//bindata = vbytes||depth||fingerprint||i||chaincode||key
//	bindata := append(w.Vbytes, append(depth, append(w.Fingerprint, append(w.I, append(w.Chaincode, w.PublicKey...)...)...)...)...)
//	// chksum := dblSha256(bindata)[:4]
//	// return append(bindata, chksum...)
//	return bindata
//}

// String returns the base32-encoded string form of the wallet.
//func (w *HDKey) String() string {
//	var s string
//	switch {
//	case bytes.Compare(w.Vbytes, Private) == 0:
//		s, _ = strkey.Encode(strkey.VersionByteSeed, w.Serialize())
//	case bytes.Compare(w.Vbytes, Public) == 0:
//		s, _ = strkey.Encode(strkey.VersionByteAccountID, w.Serialize())
//	}
//	return s
//}

// StringWallet returns a wallet given a base32-encoded extended key
//func StringWallet(data string) (*HDKey, error) {
//	dbin, _ := verDecodeCheck(data)
//	// if err := ByteCheck(dbin); err != nil {
//	// 	return &HDKey{}, err
//	// }
//	// if bytes.Compare(dblSha256(dbin[:(len(dbin) - 4)])[:4], dbin[(len(dbin)-4):]) != 0 {
//	// 	return &HDKey{}, errors.New("Invalid checksum")
//	// }
//	vbytes := dbin[0:4]
//	depth := byteToUint16(dbin[4:5])
//	fingerprint := dbin[5:9]
//	i := dbin[9:13]
//	chaincode := dbin[13:45]
//	key := dbin[45:77]
//	return &HDKey{vbytes, depth, fingerprint, i, chaincode, key}, nil
//}

// StringChild returns the ith base32-encoded extended key of a base32-encoded extended key.
// func StringChild(data string, path string) (string, error) {
// 	w, err := StringWallet(data)
// 	if err != nil {
// 		return "", err
// 	} else {
// 		wRes := w.Derive(path)
// 		return wRes.String(), nil
// 	}
// }

//StringToAddress returns the Bitcoin address of a base32-encoded extended key.
// func StringAddress(data string) (string, error) {
// 	w, err := StringWallet(data)
// 	if err != nil {
// 		return "", err
// 	} else {
// 		return w.Address(), nil
// 	}
// }

// Address returns address represented by wallet w.
// func (w *HDKey) Address() string {
// 	//!!!!!!!!
// 	four, _ := hex.DecodeString("04")
// 	padded_key := append(four, append(x.Bytes(), y.Bytes()...)...)
// 	var prefix []byte
// 	if bytes.Compare(w.Vbytes, TestPublic) == 0 || bytes.Compare(w.Vbytes, TestPrivate) == 0 {
// 		prefix, _ = hex.DecodeString("6F")
// 	} else {
// 		prefix, _ = hex.DecodeString("00")
// 	}
// 	addr_1 := append(prefix, hash(padded_key)...)
// 	chksum := dblSha256(addr_1)
// 	return base58.Encode(append(addr_1, chksum[:4]...))
// }

// GenSeed returns a random seed with a length measured in bytes.
// The length must be at least 128.

// StringCheck is a validation check of a base58-encoded extended key.
// func StringCheck(key string) error {
// 	return ByteCheck(base58.Decode(key))
// }
//
//func ByteCheck(dbin []byte) error {
//	// check proper length
//	if len(dbin) != 81 {
//		return errors.New("invalid string")
//	}
//	// check for correct Public or Private vbytes
//	if bytes.Compare(dbin[:4], Public) != 0 && bytes.Compare(dbin[:4], Private) != 0 {
//		return errors.New("invalid string")
//	}
//	// if Public, check x coord is on curve
//	return nil
//}
