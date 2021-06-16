package main

import (
	"database/sql"
	"errors"
	"flag"
	"github.com/svetlakoff88/dblib/connect"
	"github.com/svetlakoff88/dblib/drivers"
	"log"
)

func main() {
	driver, err := drivers.InstalledDrivers()
	if err != nil {
		log.Fatal(errors.New("error main installed drivers"))
	}
	for _, v := range driver {
		log.Printf("found driver: %s\n", v)
	}
	best, err := drivers.BestDriver()
	if err != nil {
		log.Fatal(errors.New("best driver main error"))
	}
	log.Printf("best driver: %s\n", best)
	fqdn := flag.String("fqdn", "", "fqdn to test connecting")
	flag.Parse()
	if *fqdn == "" {
		return
	}
	log.Printf("connecting to: %s\n", *fqdn)
	cxn := connect.Connection{
		Server:  *fqdn,
		Trusted: true,
	}
	s, err := cxn.ConnectionString()
	if err != nil {
		log.Fatal(errors.New("connection string error"))
	}
	db, err := sql.Open("odbc", s)
	if err != nil {
		log.Fatal("db connect error")
	}
	defer db.Close()
	var serverName string
	err = db.QueryRow("SELECT @@SERVERNAME").Scan(&serverName)
	if err != nil {
		log.Fatal("query executing error")
	}
	log.Printf("@@SERVERNAME: %s\n", serverName)
}
