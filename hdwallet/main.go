package hdwallet
import (
	"bitbucket.com/atticlab/go-stellar-base/strkey"
	"bitbucket.com/atticlab/go-stellar-base/hdkey"
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
	Key string
	Sum uint32
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

	if sliceCompare(ver, versionBytes["mpub"]) {
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
//TODO: Check this cycle
	for i := uint32(0); i < currentLookAhead; i++ {
		derivedKey, _ := hdw.HDK.Derive(path + strconv.FormatUint(uint64(i), 10))
		accountId, _ := strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
		accountStatus := checkAccount(accountId)
		if !accountStatus.Valid {
			currentLookAhead++
		} else if accountStatus.Balance > uint32(0) {
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
			if accountStatus.Valid == false {
				currentLookAhead++
			} else if accountStatus.Balance  > 0 {
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



func (hdw *HDWallet) MakeWithdrawalList(sum uint32) *[]accountList {
	list := []accountList{{"0", 0}}
	currentBranchAhead := branchAhead
	path := []string{"m/1/", "m/2/",}
	currentPath := path[0]
	resList, suc := hdw.findMoneyInBranch(list, currentPath, sum)
	list = resList
	if suc {
		return &list
	}
	for d := uint32(0); d < currentBranchAhead; d++{

		currentPath = path[1] + strconv.FormatUint(uint64(d), 10) + "/"

		resList, suc := hdw.findMoneyInBranch(list, currentPath, sum)
		list = resList
		if suc {
			break
		} else {
			currentBranchAhead++
		}
	}

	return &list
}

func (hdw *HDWallet) findMoneyInBranch(list []accountList, path string, sum uint32) (resultList []accountList, success bool) {
	zeroPair := accountList{"0", 0}
	currentSum := uint32(0)
	j := uint32((len(list) - 1))
	if j > 0 {
		for i := 0; i < len(list); i ++ {
			currentSum += uint32(list[i].Sum)
		}
	}
	currentLookAhead := lookAhead
	for i := uint32(0); i < currentLookAhead; i++ {
		derivedKey, _ := hdw.HDK.Derive(path + strconv.FormatUint(uint64(i), 10))
		accountId, _ := strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
		accountStatus := checkAccount(accountId)

		if accountStatus.HasBalance {
			if currentSum + accountStatus.Balance < sum {
				currentSum += accountStatus.Balance
				list[j].Key, _ = strkey.Encode(strkey.VersionByteSeed, derivedKey.PrivateKey)
				list[j].Sum = accountStatus.Balance
				list = append(list, zeroPair)
				currentLookAhead++
				j += 1
			} else if currentSum + accountStatus.Balance >= sum {
				delta := sum - currentSum
				list[j].Key, _ = strkey.Encode(strkey.VersionByteSeed, derivedKey.PrivateKey)
				list[j].Sum = delta
				resultList = list
				return resultList, true
			}}
	}
	resultList = list
	return resultList, false

}


func (hdw *HDWallet) MakeInvoiceList(sum uint32) *[]accountList {
	var path string
	list := []accountList{{"0", 0}}
	zeroPair := accountList{"0", 0}
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
				list[j].Key, _ = strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
				list[j].Sum = accountBalanceLimit
				list = append(list, zeroPair)
				currentLookAhead++
				j++
			} else if currentSum + accountBalanceLimit >= sum {
				delta := sum - currentSum
				list[j].Key, _ = strkey.Encode(strkey.VersionByteAccountID, derivedKey.PublicKey)
				list[j].Sum = delta
				break
			}
		}
	}
	return &list
}

func checkAccount(accountId string)  *accountInfo {
	id, _ := strkey.Decode(strkey.VersionByteAccountID, accountId)
	a := (id[0] & 1) > 0
	b := (id[0] & 2) > 0
	var c uint32
	if a&&b {
		c = uint32(id[0] ^ 5)
	} else {
		c = uint32(0)
	}
	return &accountInfo{a, a&&b, c}
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

