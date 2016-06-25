package hdkey

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/agl/ed25519/edwards25519"
	"github.com/stellar/go-stellar-base/strkey"
)

func hash(data []byte) []byte {
	sha := sha256.New()
	sha.Write(data)
	return sha.Sum(nil)
}

func dblSha256(data []byte) []byte {
	sha1 := sha256.New()
	sha2 := sha256.New()
	sha1.Write(data)
	sha2.Write(sha1.Sum(nil))
	return sha2.Sum(nil)
}

func modN(n []byte) []byte {
	var a [64]byte
	var b [32]byte
	for i := 32; i <64; i++{
		a[i] = 0
	}
	copy(a[:32], n[:32])
	edwards25519.ScReduce(&b, &a)
	r := b[:32]
	return r

}

func PrivToPub(privateKey []byte) []byte {
	var A edwards25519.ExtendedGroupElement
	var res []byte
	var hBytes [32]byte
	var publicKey [32]byte
	copy(hBytes[:], privateKey[:32])
	edwards25519.GeScalarMultBase(&A, &hBytes)
	A.ToBytes(&publicKey)
	res = publicKey[:32]
	return res
}

func AddPrivKeys(k1, k2 []byte) []byte {
	var k [32]byte
	flag := 0
	for i := 0; i < 32; i++ {
		a := int(k1[i])
		b := int(k2[i])
		buf := a + b + flag
	 	flag = 0
	 	if buf > 255 {
	 		flag = 1
	 	}
		buf = buf % 256
	 	k[i] = byte(buf)
	}
	kS := k[:32]

	return modN(kS)
}

func AddPubKeys(public, tweak []byte) []byte {
	var Pub, Tw, Sum edwards25519.ExtendedGroupElement
	var TwCach edwards25519.CachedGroupElement
	var SumC edwards25519.CompletedGroupElement
	var pub, tw, s [32]byte
	var sum []byte
	copy(pub[:], public[:32])
	copy(tw[:], tweak[:32])
	Pub.FromBytes(&pub)
	Tw.FromBytes(&tw)
	Tw.ToCached(&TwCach)
	edwards25519.GeAdd(&SumC, &Pub, &TwCach)
	SumC.ToExtended(&Sum)
	Sum.ToBytes(&s)
	sum = s[:32]
	return sum
}

func verDecodeCheck(addressOrSeed string) (res []byte, f bool) {
	res, err := strkey.Decode(strkey.VersionByteAccountID, addressOrSeed)
	if err == nil {
		return res, false
	}

	res, err = strkey.Decode(strkey.VersionByteSeed, addressOrSeed)
	if err == nil {
		return res, true
	}

	return
}

func uint32ToByte(i uint32) []byte {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, i)
	return a
}

func uint16ToByte(i uint16) []byte {
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, i)
	return a[1:]
}

func byteToUint16(b []byte) uint16 {
	if len(b) == 1 {
		zero := make([]byte, 1)
		b = append(zero, b...)
	}
	return binary.BigEndian.Uint16(b)
}
