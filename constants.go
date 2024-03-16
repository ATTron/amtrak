package amtrak

const endPoint = "https://maps.amtrak.com/services/MapDataService/trains/getTrainsData"

const salt = "\x9a\x36\x86\xac"
const iValue = "c6eb2f7f5c4740c1a2f708fefd947d39"
const publicKey = "69af143c-e8cf-47f8-bf09-fc1f61e5cc33"
const masterSegment = 88

var allTrains = make(map[string]Train)
var attempts = 0
