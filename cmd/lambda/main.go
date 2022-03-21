package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/EduardTruuvaart/geoip/dto"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oschwald/maxminddb-golang"
)

var db *maxminddb.Reader

func main() {
	lambda.Start(runLambda)
}

func runLambda(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	bucketKey := os.Getenv("BUCKET_KEY")

	var ip string = request.Headers["x-forwarded-for"]

	if db == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := s3.New(sess)
		req, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(bucketKey),
		})
		if err != nil {
			panic(err)
		}

		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Print(err)
		}

		dbFromBytes, _ := maxminddb.FromBytes(body)
		db = dbFromBytes
	}

	off, err := db.LookupOffset(net.ParseIP(ip))

	if err != nil {
		return create403APIResponse(), nil
	}

	geoResponse := new(dto.GeoResponse)
	country := new(dto.Country)

	db.Decode(off, geoResponse)

	if geoResponse.CountryOffset > 0 {
		db.Decode(geoResponse.CountryOffset, country)
	} else {
		return create403APIResponse(), nil
	}

	countryIsoCode := strings.ToLower(country.IsoCode)

	return createAPIResponse(200, fmt.Sprintf(`{"countrycode":"%v"}`, countryIsoCode)), nil
}

func createAPIResponse(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:      code,
		Body:            body,
		IsBase64Encoded: false,
	}
}

func createJsonpAPIResponse(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:      code,
		Body:            body,
		Headers:         map[string]string{"content-type": "application/javascript"},
		IsBase64Encoded: false,
	}
}

func create403APIResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 403,
	}
}
