package insight

import (
	"sync"
)

const EarthRadius = 6371.0

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type CellRes []LatLng
type JsonRes []CellRes

var (
	mu sync.Mutex
)
