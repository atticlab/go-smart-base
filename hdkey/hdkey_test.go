package hdkey

import "fmt"
//import (
//	//"encoding/base64"
//	//"log"
//
//	//"strings"
//	"encoding/hex"
//	"fmt"
//	"github.com/stellar/go-stellar-base/keypair"
//	"github.com/stellar/go-stellar-base/strkey"
//	"github.com/agl/ed25519/edwards25519"
//	"strconv"
//	"bytes"
//	"github.com/agl/ed25519"
//	"strings"
//)
//func bigTest() {
//	logSuc := 0
//	logFail := 0
//	verSuc1 := 0
//	verSuc2 := 0
//	m := []byte("It's a test! Hell Dogs")
//	for i := 0; i < 1000; i++ {
//		seed, _ := GenSeed(256)
//		HDKM := MasterKey(seed)
//		for j := 0; j < 10; j++{
//			path11 := "m/" + strconv.Itoa(j)
//			path12 := "M/" + strconv.Itoa(j)
//			HDKC1, _ := HDKM.Derive(path11)
//			HDKC2, _ := HDKM.Derive(path12)
//			if bytes.Compare(HDKC1.PublicKey, HDKC2.PublicKey) == 0 {
//				logSuc = logSuc + 1
//			} else {
//				logFail = logFail + 1
//			}
//			var secret [64]byte
//			var pk1, pk2 [32]byte
//			copy(secret[:32], HDKC1.PrivateKey[:32])
//			copy(secret[32:], HDKC1.PublicKey[:32])
//			copy(pk1[:32], HDKC1.PublicKey[:32])
//			copy(pk2[:32], HDKC2.PublicKey[:32])
//			sig := ed25519.Sign(&secret, m)
//			if ed25519.Verify(&pk1, m, sig) {
//				verSuc1 = verSuc1 + 1
//			}
//
//			if ed25519.Verify(&pk2, m, sig) {
//				verSuc2 = verSuc2 + 1
//			}
//
//		}
//		for j := 0; j < 10; j++ {
//			path21 := "m/3/" + strconv.Itoa(j)
//			path22 := "M/3/" + strconv.Itoa(j)
//			HDKC1, _ := HDKM.Derive(path21)
//			HDKC2, _ := HDKM.Derive(path22)
//			if bytes.Compare(HDKC1.PublicKey, HDKC2.PublicKey) == 0 {
//				logSuc = logSuc + 1
//			} else {
//				logFail = logFail + 1
//			}
//			var secret [64]byte
//			var pk1, pk2 [32]byte
//			copy(secret[:32], HDKC1.PrivateKey[:32])
//			copy(secret[32:], HDKC1.PublicKey[:32])
//			copy(pk1[:32], HDKC1.PublicKey[:32])
//			copy(pk2[:32], HDKC2.PublicKey[:32])
//			sig := ed25519.Sign(&secret, m)
//			if ed25519.Verify(&pk1, m, sig) {
//				verSuc1 = verSuc1 + 1
//			}
//
//			if ed25519.Verify(&pk2, m, sig) {
//				verSuc2 = verSuc2 + 1
//			}
//		}
//		for j := 0; j < 10; j++ {
//			path21 := "m/1/4/" + strconv.Itoa(j)
//			path22 := "M/1/4/" + strconv.Itoa(j)
//			HDKC1, _ := HDKM.Derive(path21)
//			HDKC2, _ := HDKM.Derive(path22)
//			if bytes.Compare(HDKC1.PublicKey, HDKC2.PublicKey) == 0 {
//				logSuc = logSuc + 1
//			} else {
//				logFail = logFail + 1
//			}
//			var secret [64]byte
//			var pk1, pk2 [32]byte
//			copy(secret[:32], HDKC1.PrivateKey[:32])
//			copy(secret[32:], HDKC1.PublicKey[:32])
//			copy(pk1[:32], HDKC1.PublicKey[:32])
//			copy(pk2[:32], HDKC2.PublicKey[:32])
//			sig := ed25519.Sign(&secret, m)
//			if ed25519.Verify(&pk1, m, sig) {
//				verSuc1 = verSuc1 + 1
//			}
//
//			if ed25519.Verify(&pk2, m, sig) {
//				verSuc2 = verSuc2 + 1
//			}
//		}
//		for j := 0; j < 10; j++ {
//			path21 := "m/13/4/3/71/94/13/4/3/71/94/13/4/3/71/94/13/4/3/71/94/13/4/3/71/94/" + strconv.Itoa(j)
//			path22 := "M/13/4/3/71/94/13/4/3/71/94/13/4/3/71/94/13/4/3/71/94/13/4/3/71/94/" + strconv.Itoa(j)
//			HDKC1, _ := HDKM.Derive(path21)
//			HDKC2, _ := HDKM.Derive(path22)
//			if bytes.Compare(HDKC1.PublicKey, HDKC2.PublicKey) == 0 {
//				logSuc = logSuc + 1
//			} else {
//				logFail = logFail + 1
//			}
//			var secret [64]byte
//			var pk1, pk2 [32]byte
//			copy(secret[:32], HDKC1.PrivateKey[:32])
//			copy(secret[32:], HDKC1.PublicKey[:32])
//			copy(pk1[:32], HDKC1.PublicKey[:32])
//			copy(pk2[:32], HDKC2.PublicKey[:32])
//			sig := ed25519.Sign(&secret, m)
//			if ed25519.Verify(&pk1, m, sig) {
//				verSuc1 = verSuc1 + 1
//			}
//
//			if ed25519.Verify(&pk2, m, sig) {
//				verSuc2 = verSuc2 + 1
//			}
//		}
//	}
//	fmt.Println("Success = ", logSuc)
//	fmt.Println("Fail = ", logFail)
//	fmt.Println("True1 = ", verSuc1)
//	fmt.Println("True2 = ", verSuc2)
//}
//func testKeyPair() {
//	var raw [32]byte
//	seed := "SADDF3F6LSTEJ5PSQONOQ76G2AQB3LN3YQ73QZB3ZCO6MHUMBIMB3F6U"
//
//	rawS, _ := strkey.Decode(strkey.VersionByteSeed, seed)
//	s := hex.EncodeToString(rawS)
//	enc := encodeMnemonic(s)
//	fmt.Println("enc: ", enc)
//	encS := strings.Join(enc, " ")
//	fmt.Println("enc: ", encS)
//	dec := decodeMnemonic(encS)
//	fmt.Println("dec: ", dec)
//	fmt.Println("ori: ", hex.EncodeToString(rawS))
//	copy(raw[:], rawS[:32])
//	kp, _ := keypair.FromRawSeed(raw)
//	fmt.Println(kp.Seed())
//}
//
//func testHDKey() {
//	seedR := "SADDF3F6LSTEJ5PSQONOQ76G2AQB3LN3YQ73QZB3ZCO6MHUMBIMB3F6U"
//	fmt.Println(seedR)
//	path1 := "m/1"
//	path2 := "M/1"
//	hdkM := masterKeyStr(seedR)
//	m1 := hex.EncodeToString(hdkM.PrivateKey)
//	m2 := hex.EncodeToString(hdkM.PublicKey)
//	fmt.Println("Master=============")
//	fmt.Println(m1)
//	fmt.Println(m2)
//	fmt.Println("===================")
//	hdkC1, _ := derive(path1)
//	hdkC2, _ := derive(path2)
//	pk := hex.EncodeToString(hdkC1.PrivateKey)
//	pk1 := hex.EncodeToString(hdkC1.PublicKey)
//	pk2 := hex.EncodeToString(hdkC2.PublicKey)
//	fmt.Println("Keys================")
//	fmt.Println(pk)
//	fmt.Println(pk1)
//	fmt.Println(pk2)
//	fmt.Println("===================")
//
//}
//func testMath() {
//	seed1 := "SADDF3F6LSTEJ5PSQONOQ76G2AQB3LN3YQ73QZB3ZCO6MHUMBIMB3F6U"
//	seed2 := "SDHOAMBNLGCE2MV5ZKIVZAQD3VCLGP53P3OBSBI6UN5L5XZI5TKHFQL4"
//	seed3 := "SCVQWNPUXGDRW2IOOM6SS5NQB4KK3Z2MH7ZMM4O6CXKX5L3NRZ5E6V2J"
//	rawS1, _ := strkey.Decode(strkey.VersionByteSeed, seed1)
//	rawS2, _ := strkey.Decode(strkey.VersionByteSeed, seed2)
//	rawS3, _ := strkey.Decode(strkey.VersionByteSeed, seed3)
//	sum1 := AddPrivKeys(rawS1, rawS2) //a+b
//	sum2 := AddPrivKeys(rawS2, rawS3) //b+c
//	sum3 := AddPrivKeys(rawS3, rawS1) //c+a
//
//	sum12 := AddPrivKeys(sum2, rawS1) //+a
//	sum13 := AddPrivKeys(sum3, rawS2) //+c
//	sum11 := AddPrivKeys(sum1, rawS3) //+c
//	res1 := PrivToPub(sum1)
//	p1 := PrivToPub(rawS1)
//	p2 := PrivToPub(rawS2)
//	res2 := AddPubKeys(p1, p2)
//	fmt.Println("Result --------------")
//	fmt.Println(hex.EncodeToString(res1))
//	fmt.Println(hex.EncodeToString(res2))
//	fmt.Println("Result --------------")
//	fmt.Println(hex.EncodeToString(sum11))
//	fmt.Println(hex.EncodeToString(sum12))
//	fmt.Println(hex.EncodeToString(sum13))
//}
//func gTest() {
//	var p1, p2 []byte
//	var a, b [32]byte
//	var pa, pb, pc [32]byte
//	a[0] = byte(1)
//	b[0] = byte(2)
//	var A, B, C edwards25519.ExtendedGroupElement
//	edwards25519.GeScalarMultBase(&A, &a)
//	edwards25519.GeScalarMultBase(&B, &a)
//	edwards25519.GeScalarMultBase(&C, &b)
//	A.ToBytes(&pa)
//	B.ToBytes(&pb)
//	C.ToBytes(&pc)
//	p1 = pa[:32]
//	p2 = pb[:32]
//
//	res1 := pc[:32]
//	res2 := AddPubKeys(p1, p2)
//	fmt.Println(hex.EncodeToString(res1))
//	fmt.Println(hex.EncodeToString(res2))
//}
//
func main() {
	fmt.Println("Hello World!")
	//testKeyPair()
	//testHDKey()
	//testMath()
	//gTest()
	//bigTest()

}
