package amtrak

type Train struct {
	ID         int        `json:"id"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	GeoType     string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	ObjectID                 int      `json:"OBJECTID"`
	Longitude                *float64 `json:"lon"`
	Latitude                 *float64 `json:"lat"`
	Gx_ID                    string   `json:"gx_id"`
	StatusMsg                string   `json:"StatusMsg"`
	Stations                 []Station
	Heading                  string `json:"Heading"`
	LastValueTS              string `json:"LastValTS"`
	EventCode                string `json:"EventCode"`
	DestinationCode          string `json:"DestCode"`
	OriginCode               string `json:"OrigCode"`
	RouteName                string `json:"RouteName"`
	TrainState               string `json:"TrainState"`
	OriginTZ                 string `json:"OriginTZ"`
	OriginScheduledDeparture string `json:"OrigSchDep"`
	TrainNum                 string `json:"TrainNum"`
	Velocity                 string `json:"velocity"`
	CMSID                    string `json:"CMSID"`
	ID                       int    `json:"ID"`
}

type Station struct {
	Station          string
	Code             string `json:"code"`
	TZ               string `json:"tz"`
	Bus              bool   `json:"bus"`
	ScheduledArrival string `json:"scharr"`
	ScheduledComment string `json:"schcmnt"`
	AutoArrive       bool   `json:"autoarr"`
	AutoDepart       bool   `json:"autodep"`
	EstimatedArrival string `json:"estarr"`
	EstimatedComment string `json:"estarrcmnt"`
	ActualArrival    string `json:"postarr"`
	ActualDeparture  string `json:"postdep"`
	ActualComment    string `json:"postcmnt"`
}

var TimeZones = map[string]string{
	"E": "EST",
	"C": "CST",
	"M": "MST",
	"P": "PST",
}
