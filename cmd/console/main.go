package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/EduardTruuvaart/geoip/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oschwald/maxminddb-golang"
)

func main() {
	ip := "1.2.3.4"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)
	req, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("eduardtruuvaart-geoip"),
		Key:    aws.String("GeoIP2-Country.mmdb"),
	})
	if err != nil {
		panic(err)
	}

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
	}

	db, _ := maxminddb.FromBytes(body)
	off, err := db.LookupOffset(net.ParseIP(ip))

	if err != nil {
		log.Fatalln("IP lookup error", err)
	}

	geoResponse := new(dto.GeoResponse)
	country := new(dto.Country)

	db.Decode(off, geoResponse)

	if geoResponse.CountryOffset > 0 {
		db.Decode(geoResponse.CountryOffset, country)
	} else {
		fmt.Printf("Country not found for IP: %v", ip)
	}

	fmt.Println(country.IsoCode)
}
