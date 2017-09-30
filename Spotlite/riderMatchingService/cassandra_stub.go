package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

const (
	KEY_SPACE        = "demo_insight"
	TABLE_NAME_TRIPS = "trips"
)

var cqlSession *gocql.Session

//////////////////////////////////////  structs ///////////////////////////////

type cqllatlng struct {
	Lat float64 `cql:"lat"`
	Lng float64 `cql:"lng"`
}
type cqldriverloc struct {
	Id     int       `cql:"id"`
	LatLng cqllatlng `cql:"latlng"`
	Radius float64   `cql:"radius"`
}

type cqltripinfo struct {
	DriverId int       `cql:"did"`
	RiderId  int       `cql:"rid"`
	Startloc cqllatlng `cql:"sloc"`
	Endloc   cqllatlng `cql:"eloc"`
	Duration float64   `cql:"duration"`
}

type cqltriprecord struct {
	TripId     int            `cql:"tripid"`
	TripInfo   cqltripinfo    `cql:"tripinfo"`
	DriverLocs []cqldriverloc `cql:"driverlocs"`
}

func createTable(s *gocql.Session, table string) error {

	if err := s.Query(table).RetryPolicy(nil).Exec(); err != nil {
		log.Printf("error creating table table=%q err=%v\n", table, err)
		return err
	}

	return nil
}

func buildTypesAndTables() error {

	var err error

	// create Keyspace

	err = createTable(cqlSession, `CREATE KEYSPACE IF NOT EXISTS demo_insight WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' :1};`)
	// Create the types
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = createTable(cqlSession, `CREATE TYPE IF NOT EXISTS latlng (
					lat double,
					lng double);`)

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Created latlng")

	err = createTable(cqlSession, `CREATE TYPE IF NOT EXISTS driverloc (
					id int,
					latlng frozen<latlng>,
					radius double);`)

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Created driverloc")

	err = createTable(cqlSession, `CREATE TYPE IF NOT EXISTS tripinfo (
					did int,
					rid int,
					sloc frozen<latlng>,
					eloc frozen<latlng>,
					duration double);`)

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Created tripinfo")

	err = createTable(cqlSession, `DROP TABLE IF EXISTS alltrips;`)
	err = createTable(cqlSession, `CREATE TABLE IF NOT EXISTS alltrips (
					id int,
					starttime timestamp,
					tripinfo frozen<tripinfo>,
					driverlocs list<frozen <driverloc>>,
					primary key(id)
				);`)

	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Created table alltrips")
	return nil
}

func addRecordToCassandra(recordId int, tripinfo cqltripinfo, driverlocs []cqldriverloc) error {
	fmt.Println("Adding record to cassandra")
	/*
		var tinfo = cqltripinfo{1, 2, cqllatlng{11.2, 22.2}, cqllatlng{44.4, 55.5}, 540}
		var dlocs = []cqldriverloc{
			cqldriverloc{1, cqllatlng{1.2, 1.2}, 1.2}, cqldriverloc{2, cqllatlng{2.3, 2.3}, 2.3},
		}
	*/
	ctime := time.Now().Truncate(time.Millisecond).UTC()
	err := cqlSession.Query("INSERT INTO alltrips(id, starttime, tripinfo, driverlocs) VALUES (?, ?, ?, ?)", recordId, ctime, tripinfo, driverlocs).Exec()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Successfully inserted")
	return nil
}

func createKeyspace() error {
	var err error

	err = createTable(cqlSession, `DROP KEYSPACE IF EXISTS `+KEY_SPACE)
	if err != nil {
		panic(fmt.Sprintf("unable to drop keyspace: %v", err))
	}

	err = createTable(cqlSession, fmt.Sprintf(`CREATE KEYSPACE %s
        WITH replication = {
                'class' : 'SimpleStrategy',
                'replication_factor' : %d
        }`, KEY_SPACE, 1))

	if err != nil {
		panic(fmt.Sprintf("unable to create keyspace: %v", err))
	}
	return nil
}

func createCassandraSession() error {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = KEY_SPACE
	cluster.ProtoVersion = 3
	cluster.Consistency = gocql.Quorum

	cqlSession, err = cluster.CreateSession()

	if err != nil {
		return err
	}

	err = createKeyspace()
	if err != nil {
		fmt.Println("Got and error creating keyspace: ", err)
	}

	if err = buildTypesAndTables(); err != nil {
		fmt.Printf("Failed to setup tables.  Error:", err)
		return err
	}
	return nil
}

/*********************************************************************************
 * Objects
 ********************************************************************************/
