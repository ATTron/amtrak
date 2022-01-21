package amtrak

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	re "regexp"
	"strconv"
	"strings"

	"github.com/ATTron/amtrak/util"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/pbkdf2"
)

const endPoint = "https://maps.amtrak.com/services/MapDataService/trains/getTrainsData"

const salt = "\x9a\x36\x86\xac"
const iValue = "c6eb2f7f5c4740c1a2f708fefd947d39"
const publicKey = "69af143c-e8cf-47f8-bf09-fc1f61e5cc33"
const masterSegment = 88

var allTrains = make(map[string]Train)
var attempts = 0

// GetAllTrains - return information on all trains
func GetAllTrains() map[string]Train {
	cleanData()
	return allTrains
}

// GetTrain - return information about 1 specific train
func GetTrain(trainNum int) Train {
	cleanData()
	return allTrains[strconv.Itoa(trainNum)]
}

// fetchData - go and get the latest train information
func fetchData() (string, error) {
	content := ""
attempt:
	for attempts < 3 {
		resp, err := http.Get(endPoint)
		util.Check(err)
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 200:
			responseData, err := ioutil.ReadAll(resp.Body)
			util.Check(err)
			encryptedContent := responseData[:len(responseData)-masterSegment]
			encryptedPrivateKey := responseData[len(responseData)-masterSegment:]
			privateKey := decryptData(encryptedPrivateKey, publicKey)
			content = decryptData(encryptedContent, strings.Split(privateKey, "|")[0])
			break attempt
		default:
			attempts++
		}
	}
	if attempts >= 3 {
		return "", util.ErrNotFound
	}
	return content, nil
}

// cleanData - massage the data out
func cleanData() {
	stationNumRegex := re.MustCompile(`\D*`)
	content, err := fetchData()
	util.Check(err)

	features := gjson.Get(content, "features")

	features.ForEach(func(tkey, tvalue gjson.Result) bool {
		geometry := tvalue.Get("geometry")
		newGeometry := Geometry{}
		rawGeo := json.RawMessage(geometry.String())
		bytes, err := rawGeo.MarshalJSON()
		util.Check(err)
		err = json.Unmarshal(bytes, &newGeometry)
		util.Check(err)

		properties := tvalue.Get("properties")
		newProperties := Properties{}
		rawProps := json.RawMessage(properties.String())
		bytes, err = rawProps.MarshalJSON()
		util.Check(err)
		err = json.Unmarshal(bytes, &newProperties)
		util.Check(err)
		properties.ForEach(func(pkey, pvalue gjson.Result) bool {
			if strings.HasPrefix(pkey.String(), "Station") && pvalue.String() != "" {
				stationNum := stationNumRegex.Split(pkey.String(), -1)
				newStation := Station{Station: strings.Join(stationNum, "")}
				rawStationData := json.RawMessage(pvalue.String())
				bytes, err := rawStationData.MarshalJSON()
				util.Check(err)

				err = json.Unmarshal(bytes, &newStation)
				util.Check(err)
				newStation.TZ = TimeZones[newStation.TZ]
				newProperties.Stations = append(newProperties.Stations, newStation)
			}
			return true
		})
		newTrain := Train{ID: int(tvalue.Get("id").Int()), Geometry: newGeometry, Properties: newProperties}
		if newTrain.Properties.TrainNum != "" {
			allTrains[newTrain.Properties.TrainNum] = newTrain
		}
		return true
	})
}

// WriteJson - write out to file called 'trains.json'
func writeJson(content string) {
	f, err := os.Create("trains.json")
	util.Check(err)

	defer f.Close()

	f.WriteString(content)

}

// DecryptData - decrypt the large JSON provided by amtrak
func decryptData(encrypted []byte, key string) string {
	cipherTextDecoded, err := base64.StdEncoding.DecodeString(string(encrypted))
	util.Check(err)

	encryptedIV, err := hex.DecodeString(iValue)
	util.Check(err)

	theKey := pbkdf2.Key([]byte(key), []byte(salt), 1e3, 32, sha1.New)
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
