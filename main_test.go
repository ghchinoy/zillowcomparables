package main


import (
	"testing"
	"io/ioutil"

	"fmt"
	"os"
)

func TestParseSampleResults(t *testing.T) {

	xmlbytes, err := ioutil.ReadFile("test/test_results.xml")
	if err != nil {
		t.Error(err.Error())
	}

	r, err := parseXMLToResponse(xmlbytes)
	properties := r.Comparables

	// length of comps
	//log.Println("Comps",len(comps))

	// list them
	listProperties(properties)

	if len(properties) != 25 {
		t.Error(fmt.Sprintf("Expected 25 comparables, received %v", len(properties)))
	}
}

func TestCreateGoogleMapsMatrixURL(t *testing.T) {
	mapsapikey := os.Getenv("MAPSAPI")

	xmlbytes, err := ioutil.ReadFile("test/test_results.xml")
	if err != nil {
		t.Error(err.Error())
	}
	r, err := parseXMLToResponse(xmlbytes)

	if mapsapikey != "" {
		fmt.Println(googleMapsDistanceMatrixURL(r, mapsapikey))
	}
}
