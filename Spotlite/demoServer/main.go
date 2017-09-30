package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/gorilla/mux"
)

type LatLng struct {
	Lat, Lng float64
}

type DriverInfo struct {
	Lat, Lng float64
	Radius   float64
}

type CellRes []LatLng
type JsonRes []CellRes

func main() {
	r := mux.NewRouter()
	//r.HandleFunc("/Cells", cellsHandler)
	r.HandleFunc("/Overlap", overlapHandler)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("html/"))))

	http.Handle("/", r)
	var port = os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}

	//testOverlap()

	fmt.Println("Listening on port ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}

/*
func testOverlap() {
	//ll := LatLng{37.78007695280165, -122.45520114898682}
	p := LatLng{37.77983951999081, -122.45412826538086}
	d1 := LatLng{37.789878870252856, -122.44739055633545}
	d2 := LatLng{37.77536207277254, -122.46322631835938}
	d3 := LatLng{37.77994127700315, -122.47017860412598}

	pc := getPointCells_(p)
	c1 := getCircleCells_(d1, 1.0)
	c2 := getCircleCells_(d2, 1.0)
	c3 := getCircleCells_(d3, 2.0)

	//res := c1.IntersectsCellID(pcells[0])
	fmt.Println("p and c1: ", c1.IntersectsCellID(pc[0]))
	fmt.Println("p and c2: ", c2.IntersectsCellID(pc[0]))
	fmt.Println("p and c3: ", c3.IntersectsCellID(pc[0]))

}


func getCells(ll LatLng) JsonRes {
	l := s2.LatLngFromDegrees(ll.Lat, ll.Lng)
	u := s2.LatLngFromDegrees(ll.Lat+.010, ll.Lng+.010)

	var res JsonRes

	rect := s2.RectFromLatLng(l)
	rect = rect.AddPoint(u)

	rc := &s2.RegionCoverer{MinLevel: 15, MaxLevel: 20, MaxCells: 25}

	r := s2.Region(rect.CapBound())
	covering := rc.Covering(r)

	for _, c := range covering {
		cell := s2.CellFromCellID(c)
		var cres CellRes
		for i := 0; i < 4; i++ {
			p := cell.Vertex(i)
			pll := s2.LatLngFromPoint(p)
			cres = append(cres, LatLng{pll.Lat.Degrees(), pll.Lng.Degrees()})
		}
		res = append(res, cres)
	}
	return res
}

func cellsHandler(rw http.ResponseWriter, req *http.Request) {
	kv := req.URL.Query()
	l1, _ := strconv.ParseFloat(kv["lat"][0], 64)
	l2, _ := strconv.ParseFloat(kv["lng"][0], 64)

	ll := LatLng{l1, l2}

	// THIS is to find all the locations withing a given radius
	res := getCells(ll)
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(res)
}
*/

func overlapHandler(rw http.ResponseWriter, req *http.Request) {
	kv := req.URL.Query()
	l1, _ := strconv.ParseFloat(kv["lat"][0], 64)
	l2, _ := strconv.ParseFloat(kv["lng"][0], 64)

	ll := LatLng{l1, l2}

	// THIS is to find all the locations withing a given radius
	res := isOverlap(ll)
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(res)
}

// kmToAngle converts a distance on the Earth's surface to an angle.
func kmToAngle(km float64) s1.Angle {
	// The Earth's mean radius in kilometers (according to NASA).
	const earthRadiusKm = 6371.01
	res := s1.Angle(km / earthRadiusKm)
	return res
}

func isOverlap(ll LatLng) CellRes {

	var drivers []DriverInfo

	var res CellRes

	drivers = append(drivers, DriverInfo{37.775760, -122.412360, 2.346178})
	drivers = append(drivers, DriverInfo{37.790050, -122.403840, 0.000000})
	drivers = append(drivers, DriverInfo{37.687110, -122.451570, 0.953019})
	drivers = append(drivers, DriverInfo{37.677030, -122.388410, 2.257210})
	drivers = append(drivers, DriverInfo{37.783710, -122.424370, 1.000000})
	drivers = append(drivers, DriverInfo{37.768670, -122.427360, 0.000000})
	drivers = append(drivers, DriverInfo{37.749700, -122.396470, 3.856807})
	drivers = append(drivers, DriverInfo{37.775210, -122.424350, 0.000000})

	fmt.Println(" ------------------------------------ ")
	rll := s2.LatLngFromDegrees(ll.Lat, ll.Lng)

	for _, di := range drivers {

		dll := s2.LatLngFromDegrees(di.Lat, di.Lng)
		l := s2.RegularLoop(s2.PointFromLatLng(dll), kmToAngle(di.Radius), 20)

		cont := l.ContainsPoint(s2.PointFromLatLng(rll))
		if cont == true {
			res = append(res, LatLng{di.Lat, di.Lng})
			fmt.Println("Overlaps with: ", di.Lat, di.Lng, di.Radius)
		}
	}
	fmt.Println(" ------------------------------------ ")
	return res
}

/*
func getPointCells_(ll LatLng) s2.CellUnion {
	l := s2.LatLngFromDegrees(ll.Lat, ll.Lng)
	rect := s2.RectFromLatLng(l)
	rect = rect.AddPoint(l)

	rc := &s2.RegionCoverer{MinLevel: 17, MaxLevel: 17, MaxCells: 5}

	r := s2.Region(rect.CapBound())
	covering := rc.Covering(r)
	return covering
}

func getPointCells(ll LatLng) JsonRes {
	covering := getPointCells_(ll)

	var res JsonRes
	for _, c := range covering {
		cell := s2.CellFromCellID(c)
		var cres CellRes
		for i := 0; i < 4; i++ {
			p := cell.Vertex(i)
			pll := s2.LatLngFromPoint(p)
			fmt.Println("Point found ", pll.String())
			cres = append(cres, LatLng{pll.Lat.Degrees(), pll.Lng.Degrees()})
		}
		res = append(res, cres)
	}
	return res
}

func pointcellsHandler(rw http.ResponseWriter, req *http.Request) {
	kv := req.URL.Query()
	fmt.Println("Inside pointcellshandler: ", req.URL.String())

	l1, _ := strconv.ParseFloat(kv["lat"][0], 64)
	l2, _ := strconv.ParseFloat(kv["lng"][0], 64)

	ll := LatLng{l1, l2}

	// THIS is to find all the locations withing a given radius
	res := getPointCells(ll)
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(res)
}

// kmToAngle converts a distance on the Earth's surface to an angle.
func kmToAngle(km float64) s1.Angle {
	// The Earth's mean radius in kilometers (according to NASA).
	const earthRadiusKm = 6371.01
	res := s1.Angle(km / earthRadiusKm)
	return res
}

func getCircleCells_(ll LatLng, kmRadius float64) s2.CellUnion {
	l := s2.LatLngFromDegrees(ll.Lat, ll.Lng)

	rc := &s2.RegionCoverer{MinLevel: 15, MaxLevel: 20, MaxCells: 25}
	rd := kmToAngle(kmRadius)
	rect := s2.CapFromCenterAngle(s2.PointFromLatLng(l), rd)
	r := s2.Region(rect.CapBound())
	covering := rc.Covering(r)
	return covering
}

func getCircleCells(ll LatLng, radius float64) JsonRes {

	covering := getCircleCells_(ll, radius)
	var res JsonRes

	for _, c := range covering {
		cell := s2.CellFromCellID(c)
		var cres CellRes
		for i := 0; i < 4; i++ {
			p := cell.Vertex(i)
			pll := s2.LatLngFromPoint(p)
			cres = append(cres, LatLng{pll.Lat.Degrees(), pll.Lng.Degrees()})
		}
		res = append(res, cres)
	}
	return res
}

func circlecellsHandler(rw http.ResponseWriter, req *http.Request) {
	kv := req.URL.Query()
	fmt.Println("Inside circlecellshandler: ", req.URL.String())

	l1, _ := strconv.ParseFloat(kv["lat"][0], 64)
	l2, _ := strconv.ParseFloat(kv["lng"][0], 64)
	r, _ := strconv.ParseFloat(kv["radius"][0], 64)
	ll := LatLng{l1, l2}

	// THIS is to find all the locations withing a given radius
	res := getCircleCells(ll, r)
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(res)
}
*/
