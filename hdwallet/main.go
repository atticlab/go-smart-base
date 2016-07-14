package hdwallet
import (
	"github.com/stellar/go-stellar-base/strkey"
	"github.com/stellar/go-stellar-base/hdkey"
	"encoding/binary"
	"encoding/hex"
	"strconv"

)

const (
	lookAhead = uint32(20)
	branchAhead = uint32(20)
	accountBalanceLimit = uint32(500)
)
var versionBytes = make(map[string][]byte)

func init() {
	versionBytes["accountId"], _ = hex.DecodeString("30") // "G" in base32
	versionBytes["seed"], _ = hex.DecodeString("90") // "S" in base32
	versionBytes["mpriv"], _ = hex.DecodeString("60") // "M" in base32
	versionBytes["mpub"], _ = hex.DecodeString("78") // "P" in base32
	versionBytes["privW"], _ = hex.DecodeString("b0") // "W" in base32
	versionBytes["pubW"], _ = hex.DecodeString("c8") // "Z" in base32
}
type accountInfo struct {
	Valid bool
	HasBalance bool
	Balance uint32
}

type accountList struct {
	Key []string
	Sum []uint32
}

type HDWallet struct {
	Vbytes      	[]byte //1 bytes
	FirstWithMoney 	uint32 //4 bytes
	FirstUnused     uint32 //4 bytes uint32ToByte()
	Seed   		[]byte //32 bytes
	Maps		[]uint32
	HDK 		hdkey.HDKey
}

func SetByStrKey(str string) *HDWallet {
	switch str[0] {
	case "P"[0]: {
		key, err := strkey.Decode(strkey.VersionByteMPub, str)
		if err != nil {
			panic(err)
		}
		return initKey(versionBytes["mpub"] , key)
	}
	case "S"[0]: {
		key, err := strkey.Decode(strkey.VersionByteSeed, str)
		if err != nil {
			panic(err)
		}
		return SetBySeed(key)
	}
	case "W"[0]: {
		key, err := strkey.Decode(strkey.VersionBytePrivWallet, str)
		if err != nil {
			panic(err)
		}
		return Deserialize(versionBytes["mpriv"], key)
	}
	case "Z"[0]: {
		key, err := strkey.Decode(strkey.VersionBytePubWallet, str)
		if err != nil {
			panic(err)
		}
		return Deserialize(versionBytes["mpub"], key)
	}
	default:
		panic("Invalid version of StrKey")
	}
}

func SetByPhrase(s string)  *HDWallet {
	strSeed := hdkey.GetSeedFromMnemonic(s)
	seed, _ := hex.DecodeString(strSeed)
	return SetBySeed(seed)
}

func SetBySeed(seed []byte) *HDWallet {
	hdw := new(HDWallet)
	hdw.Vbytes = versionBytes["mpriv"]
	hdw.Seed = seed
	hdw.HDK = *hdkey.MasterKey(seed, hdw.Vbytes)
	return setAllIndex(hdw)
}

func Deserialize(ver []byte, serWallet []byte) *HDWallet {
	//var tmpMap []uint32
	hdw := new(HDWallet)
	hdw.Vbytes = ver
	if sliceCompare(ver, versionBytes["mpriv"]) {
		hdw.Seed = serWallet[0:32]
		hdw.HDK = *hdkey.MasterKey(hdw.Seed, ver)
	} else if sliceCompare(ver, versionBytes["mpub"]) {
		hdw.HDK.PublicKey = serWallet[0:32]
		hdw.HDK.Chaincode = serWallet[32:64]
	}
	hdw.FirstWithMoney =  byteToUInt32(serWallet[64:68])
	hdw.FirstUnused = byteToUInt32(serWallet[68:72])
	mapsLen := byteToUInt32(serWallet[72:76])

	tmpMap := []uint32{byteToUInt32(serWallet[76 : 80]),}

	j := uint32(0)
	for i := uint32(1) ; i < mapsLen; i++ {
		tmpMap = append(tmpMap, byteToUInt32(serWallet[80 + j : 84 + j]))
		j += 4
	}
	hdw.Maps = tmpMap
	return hdw
}
func initKey(ver []byte, rawKey []byte)  *HDWallet {
	var masterKP hdkey.HDKey
	masterKP.Vbytes = ver
	masterKP.I, _ = hex.DecodeString("0")
	masterKP.Depth = uint16(0)
	masterKP.Chaincode = rawKey[:32]

	if sliceCompare(ver, versionBytes["mpriv"]) {
		masterKP.PrivateKey = rawKey[32:]
		masterKP.PublicKey = hdkey.PrivToPub(rawKey[32:])
	} else if sliceCompare(ver, versionBytes["mpub"]) {
		masterKP.PrivateKey = nil
		masterKP.PublicKey = rawKey[32:]
	}

	hdw := new(HDWallet)
	hdw.Vbytes = ver
	hdw.HDK = masterKP

	return setAllIndex(hdw)
}
func setAllIndex(hdw *HDWallet) *HDWallet{
	var path string
	var FirstUnused, FirstWithMoney uint32
	FirstUnused = 0
	FirstWithMoney = 0
	currentLookAhead := lookAhead

	if sliceCompare(hdw.Vbytes, versionBytes["mpriv"]) {
		path = "m/1/"
		hdw.Maps = getPublicMap(&hdw.HDK)
	} else if sliceCompare(hdw.Vbytes, versionBytes["mpub"]) {
		path = "M/"
	} else {
		panic("Invalid version")
	}

	for i := uint32(0); i < currentLookAhead; i++ {
		derivedKey, _ := hdw.HDK.Derive(path + strconv.FormatUint(uint64(i), 10))
		accountId, _ := strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
		accountStatus := checkAccount(accountId)
		if !accountStatus.Valid {
			currentLookAhead++
		} else if (accountStatus.Balance > 0) && (FirstWithMoney == 0) {
			FirstWithMoney = i
		}
		FirstUnused = i + 1
	}

	hdw.FirstWithMoney = FirstWithMoney
	hdw.FirstUnused = FirstUnused
	return hdw
}

func getPublicMap(hd *hdkey.HDKey) []uint32 {
	var maps []uint32
	path := "M/2/"
	currentLookAhead := lookAhead
	currentBranchAhead := branchAhead
	j := uint32(0)
	for d := uint32(0); d < currentBranchAhead; d ++ {
		jT := j
		maps = append(maps, uint32(0));
		for i := uint32(0); i < currentLookAhead; i ++ {
			derivedKey, _ := hd.Derive(path + strconv.FormatUint(uint64(d), 10) + "/" + strconv.FormatUint(uint64(i), 10))
			accountId, _ := strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
			accountStatus := checkAccount(accountId)
			if !accountStatus.Valid {
				currentLookAhead++
			} else if (accountStatus.Balance > 0) && maps[j] == 0 {
				maps[j] = i
				j++
				break
			}
		}
		if j == jT {
			currentBranchAhead++
		}}

	return maps
}

func (hdw *HDWallet) Serialize() string {
	var key, arr []byte
	var ver strkey.VersionByte
	mapsLen := uint32(len(hdw.Maps))
	if sliceCompare(hdw.Vbytes, versionBytes["mpriv"]) {
		ver = strkey.VersionBytePrivWallet
		key = hdw.Seed[:]
	} else if sliceCompare(hdw.Vbytes, versionBytes["mpub"]) {
		ver = strkey.VersionBytePubWallet
		key = hdw.HDK.PublicKey[:]
	}
	tmpMap := uint32ToByte(hdw.Maps[0])
	for i := uint32(1); i < mapsLen; i++ {
		tmpMap = append(tmpMap, uint32ToByte(hdw.Maps[i])...)
	}

	arr = append(key, append(hdw.HDK.Chaincode, append(uint32ToByte(hdw.FirstWithMoney), append(uint32ToByte(hdw.FirstUnused), append(uint32ToByte(mapsLen), tmpMap...)...)...)...)...)
	str, _ := strkey.Encode(ver, arr)
	return str
}

func (hdw *HDWallet) MakeWithdrawalList(sum uint32) *accountList {
	list := new(accountList)
	currentSum := uint32(0)
	currentLookAhead := lookAhead
	currentBranchAhead := branchAhead
	path := []string{"m/1/", "m/2/",}
	for p := 0; p < 2; p++ {
		currentPath := path[p]
		jd := uint32(0)
		for d := uint32(0); d < currentBranchAhead; d++{
			jT := jd
			j := uint32(0)
			if p == 1 {
				currentPath = currentPath + strconv.FormatUint(uint64(d), 10) + "/"
			}
			for i := uint32(0); i < currentLookAhead; i++ {
				derivedKey, _ := hdw.HDK.Derive(currentPath + strconv.FormatUint(uint64(i), 10))
				accountId, _ := strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
				accountStatus := checkAccount(accountId)
				if accountStatus.HasBalance {
					if currentSum + accountStatus.Balance < sum {
						currentSum += accountStatus.Balance
						list.Key[j], _ = strkey.Encode(strkey.VersionByteSeed, derivedKey.PrivateKey)
						list.Sum[j] = accountStatus.Balance
						currentLookAhead++
						jd++
						j++
					} else if currentSum + accountStatus.Balance >= sum {
						delta := sum - currentSum
						list.Key[j], _ = strkey.Encode(strkey.VersionByteSeed, derivedKey.PrivateKey)
						list.Sum[j] = delta
						return list
					}}}
			if (jd == jT){
				currentBranchAhead++
		}}}

	return list
}

func (hdw *HDWallet) MakeInvioceList(sum uint32) *accountList {
	var path string
	list := new(accountList)
	currentSum := uint32(0)
	currentLookAhead := lookAhead

	if sliceCompare(hdw.Vbytes, versionBytes["mpriv"]) {
		path = "m/1/"
	} else if sliceCompare(hdw.Vbytes, versionBytes["mpub"]) {
		path = "M/"
	}
	j := uint32(0)
	for i := hdw.FirstUnused; i < i + currentLookAhead; i++ {
		derivedKey, _ := hdw.HDK.Derive(path + strconv.FormatUint(uint64(i), 10))
		accountId, _ := strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
		accountStatus := checkAccount(accountId)
		if !accountStatus.Valid {
			if currentSum + accountBalanceLimit < sum {
				currentSum += accountBalanceLimit
				list.Key[j], _ = strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
				list.Sum[j] = accountBalanceLimit
				currentLookAhead++
				j++
			} else if currentSum + accountBalanceLimit >= sum {
				delta := sum - currentSum
				list.Key[j], _ = strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
				list.Sum[j] = delta
				break
			}
		}
	}
	return list
}

func checkAccount(accountId string)  *accountInfo {
	id, _ := strkey.Decode(strkey.VersionByteAccountID, accountId)
	a := (id[0] & 1) > 0
	b := (id[0] & 2) > 0
	c := (id[0] ^ 5)
	return &accountInfo{a, a&&b, uint32(c)}
}

func uint32ToByte(i uint32) []byte {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, i)
	return a
}

func byteToUInt32 (b []byte)  uint32 {
	return binary.BigEndian.Uint32(b)
}

func sliceCompare(a, b []byte) bool {

	if a == nil && b == nil {
		return true;
	}

	if a == nil || b == nil {
		return false;
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

