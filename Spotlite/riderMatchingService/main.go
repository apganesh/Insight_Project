package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	insight "github.com/apganesh/Insight_Project/common"
	"github.com/pkg/profile"
)

const (
	port = ":50051"
)

// server is used to implement customer.CustomerServer.
type server struct {
	//drivers []*insight.DriverInfo
}

var driverMap map[int64]*insight.Driver
var mu sync.Mutex
var recordId int

func init() {
	driverMap = make(map[int64]*insight.Driver)
	recordId = 0
	err := RedisNewClient()
	if err != nil {
		log.Fatal("Cannot start redis client: ", err)
	}
	err = createCassandraSession()
	if err != nil {
		log.Fatal("Cannot start cassandra session: ", err)
	}
}

type distMap struct {
	Dist float64
	Id   int64
}

type byDistance []distMap

func (a byDistance) Len() int           { return len(a) }
func (a byDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDistance) Less(i, j int) bool { return a[i].Dist < a[j].Dist }

func sortByDistance(r insight.Rider, ids []int64) []int64 {
	var idDistMap []distMap

	for _, id := range ids {
		obj := driverMap[id]
		hdist := getHaversineDist(r.SLat, r.SLng, obj.Lat, obj.Lng)
		idDistMap = append(idDistMap, distMap{Id: id, Dist: hdist})
	}

	sort.Sort(byDistance(idDistMap))
	var res []int64
	for _, dMap := range idDistMap {
		res = append(res, dMap.Id)
	}
	return res
}

// CreateCustomer creates a new Customer
func (s *server) AddDriver(ctx context.Context, in *insight.Driver) (*insight.Driver, error) {

	mu.Lock()
	defer mu.Unlock()
	if obj, ok := driverMap[in.Id]; ok {
		obj.Lat = in.Lat
		obj.Lng = in.Lng
		obj.Radius = in.Radius
		obj.Timestamp = in.Timestamp
		obj.Status = in.Status
	} else {
		driverMap[in.Id] = in
	}

	addDriverToRedis(in)

	return in, nil
}

func getDriverRiderDistance(r *insight.Rider, d *insight.Driver) float64 {
	return getHaversineDist(r.SLat, r.SLng, d.Lat, d.Lng)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (s *server) GetDriver(ctx context.Context, in *insight.Rider) (*insight.Driver, error) {
	fmt.Println("Getting Driver for: ", in.Id)
	defer timeTrack(time.Now(), "GetDriver: ")

	res := getDrivers_redis(in)

	if len(res) == 0 {
		return &insight.Driver{}, nil
	}

	/*
		for _, r := range res {
			fmt.Println("Distance between driver/rider", r.Id, r.Lat, r.Lng, r.Radius, getDriverRiderDistance(in, &r))
		}
	*/

	// drawPoint
	//fmt.Println("drawRider(", in.SLat, ",", in.SLng, ")")

	driver := res[0]

	var driverlocs []cqldriverloc
	recordId++

	for _, d := range res {
		driverlocs = append(driverlocs, cqldriverloc{int(d.Id), cqllatlng{d.Lat, d.Lng}, d.Radius})
		//fmt.Println("drawDriverCircle(", d.Id, ",", d.Lat, ",", d.Lng, ",", d.Radius, ")")
	}

	// drawLine
	//fmt.Println("drawClosestDriver(", in.SLat, ",", in.SLng, ",", driver.Lat, ",", driver.Lng, ")")

	fmt.Println("Paired Rider/Driver  riderid driverid dist: ", in.Id, res[0].Id, getDriverRiderDistance(in, &driver), driver.Radius)

	tripinfo := cqltripinfo{int(in.Id), int(driver.Id), cqllatlng{in.SLat, in.SLng}, cqllatlng{in.ELat, in.ELng}, 500.0}

	/*
		mm := math.Min(5.0, float64(len(sortres)))

		for i := 0; i < int(mm); i++ {
			drv := driverMap[sortres[i]]
			driverlocs = append(driverlocs, cqldriverloc{int(drv.Id), cqllatlng{drv.Lat, drv.Lng}, drv.Radius})
			fmt.Println("Distance between driver rider: ", sortres[i], getDriverRiderDistance(in, driverMap[sortres[i]]))
		}
	*/
	addRecordToCassandra(recordId, tripinfo, driverlocs)

	return &driver, nil
}

func main() {
	var err error
	err = createCassandraSession()
	if err != nil {
		log.Fatal("Failed to create Cassandra ...")
		return
	}

	var lis net.Listener
	lis, err = net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	insight.RegisterMatcherServer(s, &server{})
	defer profile.Start(profile.CPUProfile).Stop()
	s.Serve(lis)
}
