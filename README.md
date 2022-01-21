# AMTRAK
An unoffical train status wrapper for use with Amtrak!  
This package uses the data available from [amtrak](https://maps.amtrak.com/services/MapDataService/trains/getTrainsData) and converts it into usable train information.

## Usage
```bash
go get -u "github.com/ATTron/amtrak"
```

```go
import (
    "fmt"
    "github.com/ATTron/amtrak"
)

func main() {
    fmt.Println("SHOWING TRAIN 95: ", amtrak.GetTrain(95))
}
```

#### Functions
Currently this only has 2 functions. Return a specific train or return every train.

```go
import (
    "github.com/ATTron/amtrak"
)
func main() {
    allTrains := amtrak.GetAllTrains()
    myTrain := amtrak.GetTrain(95)
}
```

#### Additional Docs
Full documentation can be found on [pkg.go.dev](https://pkg.go.dev/github.com/ATTron/amtrak)