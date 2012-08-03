package app

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"library/cache"
	"library/render"
	"net/http"
)

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/list", list)
}

type APN struct {
	Carrier  string `xml:"carrier,attr"`
	APN      string `xml:"apn,attr"`
	User     string `xml:"user,attr"`
	Password string `xml:"password,attr"`
	Proxy    string `xml:"proxy,attr"`
	Port     int    `xml:"port,attr"`
	MMSC     string `xml:"mmsc,attr"`
	MMSProxy string `xml:"mmsproxy,attr"`
	MMSPort  string `xml:"mmsport,attr"`
	MCC      string `xml:"mcc,attr"`
	MNC      string `xml:"mnc,attr"`
}
type APNStruct struct {
	APN []APN `xml:"apn"`
}

// Function to get all APN
func getAPN(r *http.Request) []APN {
	// Check if we have the item in the cache
	cachedItem, cacheStatus := cache.GetCache(r, "getAPN")
	if cacheStatus == true {
		var apn []APNStruct
		pAPN := bytes.NewBuffer(cachedItem)
		decAPN := gob.NewDecoder(pAPN)
		decAPN.Decode(&apn)
		return apn[0].APN
	}

	// Read out 
	data, _ := ioutil.ReadFile("apns-conf.xml")
	var apn []APNStruct
	xml.Unmarshal(data, &apn)

	// Add to cache
	mModels := new(bytes.Buffer) //initialize a *bytes.Buffer
	encModels := gob.NewEncoder(mModels)
	encModels.Encode(apn)
	cache.AddCache(r, "getAPN", mModels.Bytes())

	return apn[0].APN
}

// Home
func home(w http.ResponseWriter, r *http.Request) {
	passedTemplate := new(bytes.Buffer)
	template.Must(template.ParseFiles("apnconf/templates/home.html")).Execute(passedTemplate, getAPN(r))
	render.Render(w, r, passedTemplate)
}

// List of all APNs
func list(w http.ResponseWriter, r *http.Request) {
	passedTemplate := new(bytes.Buffer)
	template.Must(template.ParseFiles("apnconf/templates/list.html")).Execute(passedTemplate, getAPN(r))
	render.Render(w, r, passedTemplate)
}
