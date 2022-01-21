package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gitlab.com/ATTron/amtrak/util"
	"golang.org/x/crypto/pbkdf2"
)

const dataUrl = "https://maps.amtrak.com/services/MapDataService/trains/getTrainsData"

const sValue = "\x9a\x36\x86\xac"
const iValue = "c6eb2f7f5c4740c1a2f708fefd947d39"
const publicKey = "69af143c-e8cf-47f8-bf09-fc1f61e5cc33"
const masterSegment = 88

func main() {
	resp, err := http.Get(dataUrl)
	util.Check(err)

	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	util.Check(err)
	encryptedContent := responseData[:len(responseData)-masterSegment]
	encryptedPrivateKey := responseData[len(responseData)-masterSegment:]
	privateKey := decryptData(encryptedPrivateKey, publicKey)
	content := decryptData(encryptedContent, strings.Split(privateKey, "|")[0])

	fmt.Println("got private key:", privateKey)
	fmt.Println("got content:", content)

	writeJson(content)
}

func writeJson(content string) {
	f, err := os.Create("trains.json") // creates a file at current directory
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	// outcontent := []byte(strings.ReplaceAll(content, "\\", ""))

	f.WriteString(content)

}

// this is an absolute mess and its a total unnecessary PITA
func decryptData(encrypted []byte, key string) string {
	fmt.Println("trying to decrypt:", string(encrypted))
	cipherTextDecoded, err := base64.StdEncoding.DecodeString(string(encrypted))
	util.Check(err)

	encryptedIV, err := hex.DecodeString(iValue)
	util.Check(err)

	fmt.Println("GOT KEY:", string(key))
	theKey := pbkdf2.Key([]byte(key), []byte(sValue), 1e3, 32, sha1.New)
	keyEncoded := hex.EncodeToString(theKey)[:32]

	keyDecoded, err := hex.DecodeString(keyEncoded)
	util.Check(err)

	fmt.Println("GOT KEY:", keyEncoded)
	block, err := aes.NewCipher(keyDecoded)
	util.Check(err)

	mode := cipher.NewCBCDecrypter(block, encryptedIV)
	decrypted := make([]byte, len(cipherTextDecoded))
	mode.CryptBlocks(decrypted, []byte(cipherTextDecoded))

	return string(decrypted)
}

func processRespJSON(resp *http.Response) (string, error) {
	switch resp.StatusCode {
	case 200:
		b, _ := ioutil.ReadAll(resp.Body)
		var jsonResp map[string]interface{}
		err := json.Unmarshal(b, &jsonResp)
		util.Check(err)
		returnResp, err := json.MarshalIndent(&jsonResp, "", "  ")
		util.Check(err)
		return string(returnResp), nil
	case 404:
		log.Println("Unable to find the request you are looking for")
		return "", util.ErrNotFound
	default:
		log.Fatal("Could not return valid response from server . . .")
		return "", util.ErrBadType
	}
}

func processRespArray(resp *http.Response) (string, error) {
	switch resp.StatusCode {
	case 200:
		b, _ := ioutil.ReadAll(resp.Body)
		var jsonResp []interface{}
		err := json.Unmarshal(b, &jsonResp)
		util.Check(err)
		returnResp, err := json.MarshalIndent(&jsonResp, "", "  ")
		util.Check(err)
		return string(returnResp), nil
	case 404:
		log.Println("Unable to find the request you are looking for")
		return "", util.ErrNotFound
	default:
		log.Fatal("Could not return valid response from server . . .")
		return "", util.ErrBadType
	}
}
