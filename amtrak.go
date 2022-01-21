package amtrak

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gitlab.com/ATTron/amtrak/util"
	"golang.org/x/crypto/pbkdf2"
)

const EndPoint = "https://maps.amtrak.com/services/MapDataService/trains/getTrainsData"

const Salt = "\x9a\x36\x86\xac"
const IValue = "c6eb2f7f5c4740c1a2f708fefd947d39"
const PublicKey = "69af143c-e8cf-47f8-bf09-fc1f61e5cc33"
const MasterSegment = 88

// FetchData - go and get the latest train information
func FetchData() string {
	resp, err := http.Get(EndPoint)
	util.Check(err)

	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	util.Check(err)
	encryptedContent := responseData[:len(responseData)-MasterSegment]
	encryptedPrivateKey := responseData[len(responseData)-MasterSegment:]
	privateKey := DecryptData(encryptedPrivateKey, PublicKey)
	content := DecryptData(encryptedContent, strings.Split(privateKey, "|")[0])

	return content
}

// WriteJson - write out to file called 'trains.json'
func WriteJson(content string) {
	f, err := os.Create("trains.json")
	util.Check(err)

	defer f.Close()

	f.WriteString(content)

}

// DecryptData - decrypt the large JSON provided by amtrak
func DecryptData(encrypted []byte, key string) string {
	cipherTextDecoded, err := base64.StdEncoding.DecodeString(string(encrypted))
	util.Check(err)

	encryptedIV, err := hex.DecodeString(IValue)
	util.Check(err)

	theKey := pbkdf2.Key([]byte(key), []byte(Salt), 1e3, 32, sha1.New)
	keyEncoded := hex.EncodeToString(theKey)[:32]

	keyDecoded, err := hex.DecodeString(keyEncoded)
	util.Check(err)

	block, err := aes.NewCipher(keyDecoded)
	util.Check(err)

	mode := cipher.NewCBCDecrypter(block, encryptedIV)
	decrypted := make([]byte, len(cipherTextDecoded))
	mode.CryptBlocks(decrypted, []byte(cipherTextDecoded))

	return string(decrypted)
}
