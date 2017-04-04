package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/xml"
	"errors"
	"strings"
	"net/url"
	"strconv"
)

const (
	baseurlformat = "http://www.zillow.com/webservice/GetDeepComps.htm?zpid=%v&zws-id=%s&rentzestimate=%v&count=%v"
	//propertyurlformat ="http://www.zillow.com/webservice/GetDeepSearchResults.htm"
	googlematrixurlformat = "https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&origins=%s&destinations=%s&key=%s"
)

var (
	zwsid string
	mapsapi string
)

type Response struct {
	XMLName xml.Name `xml:"comps"`
	Request Request `xml:"request"`
	Message Message `xml:"message"`
	Principal Property `xml:"response>properties>principal"`
	Comparables []Property `xml:"response>properties>comparables>comp"`
}

type Request struct {
	ZPID int `xml:"zpid"`
	Count int `xml:"count"`
}
type Message struct {
	Text string `xml:"text"`
	Code int `xml:"code"`
}

type Property struct {
	Score float64 `xml:"score,attr"`
	ZPID int `xml:"zpid"`
	HomeDetailsURL string `xml:"links>homedetails"`
	GraphsAndDataURL string `xml:"graphsanddata"`
	MapThisLocationURL string `xml:"maothishome"`
	ComparablesURL string `xml:"comparables"`
	Address Address `xml:"address"`
	TaxAssessmentYear int `xml:"taxAssessmentYear"`
	TaxAssessment float64 `xml:"taxAssessment"`
	SquareFeet int `xml:"finishedSqFt"`
	Bathrooms float64 `xml:"bathrooms"`
	Bedrooms int `xml:"bedrooms"`
	RentInfo RentInfo `xml:"rentzestimate"`
}
type RentInfo struct {
	Amount int `xml:"amount"`
	Currency string `xml:"amount,attr"`
	LastUpdateDate string `xml:"last-updated"` // should be time.Time
}
type Address struct {
	Street string `xml:"street"`
	ZipCode string `xml:"zipcode"`
	City string `xml:"city"`
	State string `xml:"state"`
	Latitude float64 `xml:"latitude"`
	Longitude float64 `xml:"longitude"`
}


func main() {

	zwsid = os.Getenv("ZWSID")
	mapsapi = os.Getenv("MAPSAPI")

	if zwsid == "" {
		fmt.Println("Please provide an env var ZWSID")
		os.Exit(1)
	}

	propertystring := os.Args[1]
	if propertystring == "" {
		fmt.Println("Please provide a propertyid on the commandline")
		os.Exit(1)
	}

	propertyid, err := strconv.Atoi(propertystring)
	if err != nil {
		fmt.Println("Please provide a numeric propertyid")
		os.Exit(1)
	}
	response, err := getResponseForZillowID(propertyid)
	properties := response.Comparables
	if err != nil {
		log.Fatal("Unable to even", err.Error())
	}

	listProperties(properties)

	if mapsapi != "" {
		mapsurl := googleMapsDistanceMatrixURL(response, mapsapi)
		fmt.Print(mapsurl)
	}


}

func listProperties(properties []Property) {
	for _, v := range properties {
		fmt.Printf("%v,%v,$%v,%s,%v,%v,%v\n",
			v.Score,
			v.ZPID,
			v.RentInfo.Amount,
			v.Address.Street,
			v.Bedrooms,
			v.Bathrooms,
			v.SquareFeet,
		)
	}
}

func googleMapsDistanceMatrixURL(response Response, apikey string) string {

	var matrixURL string

	var locations []string

	for _, v := range response.Comparables {
		locations = append(locations, fmt.Sprintf("%v,%v", v.Address.Latitude, v.Address.Longitude))
	}

	matrixURL = fmt.Sprintf(googlematrixurlformat,
		url.PathEscape(fmt.Sprintf("%v,%v", response.Principal.Address.Latitude, response.Principal.Address.Longitude)),
		url.PathEscape(strings.Join(locations, "|")),
		apikey)

	return matrixURL
}

/*
https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&origins=40.6655101,-73.89188969999998&destinations=40.6905615%2C-73.9976592%7C40.6905615%2C-73.9976592%7C40.6905615%2C-73.9976592%7C40.6905615%2C-73.9976592%7C40.6905615%2C-73.9976592%7C40.6905615%2C-73.9976592%7C40.659569%2C-73.933783%7C40.729029%2C-73.851524%7C40.6860072%2C-73.6334271%7C40.598566%2C-73.7527626%7C40.659569%2C-73.933783%7C40.729029%2C-73.851524%7C40.6860072%2C-73.6334271%7C40.598566%2C-73.7527626&key=YOUR_API_KEY
 */

func getResponseForZillowID(id int) (Response, error) {

	var response Response
	var err error

	count := 25
	rentzestimate := true


	url := fmt.Sprintf(baseurlformat,
		id,
		zwsid,
		rentzestimate,
		count,
	)

	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()



	xmldata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	response, err = parseXMLToResponse(xmldata)


	return response, nil
}


func parseXMLToResponse(xmldata []byte) (Response, error) {
	var response Response

	err := xml.Unmarshal(xmldata, &response)
	if err != nil {
		return response, err
	}
	log.Println("Info for ZPID:",response.Request.ZPID)
	log.Println("Comps in response:", response.Request.Count)
	//log.Println("Message:", comps.Message.Text)
	if response.Request.ZPID != response.Principal.ZPID {
		return response, errors.New("ZPID request and ZPID principle properties don't match")
	}
	//log.Printf("%v, %v", response.Request.ZPID, response.Principal.ZPID)
	//log.Printf("%s", response.Principal.HomeDetails)

	return response, nil
}
