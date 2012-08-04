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
	"strings"
)

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/list", list)
	http.HandleFunc("/search/", search)
	http.HandleFunc("/show/", show)
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
func getAllAPN(r *http.Request) []APN {
	// Check if we have the item in the cache
	cachedItem, cacheStatus := cache.GetCache(r, "getAllAPN")
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
	cache.AddCache(r, "getAllAPN", mModels.Bytes())

	return apn[0].APN
}

// Function to get a single APN
func getAPN(r *http.Request, MCC string, MNC string) []APN {
	// Check if we have the item in the cache
	cachedItem, cacheStatus := cache.GetCache(r, "getAPN-"+MCC+"-"+MNC)
	if cacheStatus == true {
		var apn []APN
		pAPN := bytes.NewBuffer(cachedItem)
		decAPN := gob.NewDecoder(pAPN)
		decAPN.Decode(&apn)
		return apn
	}

	// Read out 
	data := getAllAPN(r)
	i := 0
	for key := range data {
		if data[key].MCC == MCC && data[key].MNC == MNC {
			i++
		}
	}
	apn := make([]APN, i)
	i = 0
	for key := range data {
		if data[key].MCC == MCC && data[key].MNC == MNC {
			apn[i] = data[key]
			i++
		}
	}

	// Add to cache
	mModels := new(bytes.Buffer) //initialize a *bytes.Buffer
	encModels := gob.NewEncoder(mModels)
	encModels.Encode(apn)
	cache.AddCache(r, "getAPN-"+MCC+"-"+MNC, mModels.Bytes())

	return apn
}

// Home
func home(w http.ResponseWriter, r *http.Request) {
	// Check if we are on the root page, and if not return a 404
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		passedTemplate := new(bytes.Buffer)
		template.Must(template.ParseFiles("templates/404.html")).Execute(passedTemplate, nil)
		render.Render(w, r, passedTemplate, http.StatusNotFound)
		return
	}

	passedTemplate := new(bytes.Buffer)
	template.Must(template.ParseFiles("apnconf/templates/home.html")).Execute(passedTemplate, getAllAPN(r))
	render.Render(w, r, passedTemplate)
}

// List of all APNs
func list(w http.ResponseWriter, r *http.Request) {
	passedTemplate := new(bytes.Buffer)
	template.Must(template.ParseFiles("apnconf/templates/list.html")).Execute(passedTemplate, getAllAPN(r))
	render.Render(w, r, passedTemplate)
}

// Search function
func search(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("APN") == "" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data := getAllAPN(r)
	for key := range data {
		if data[key].Carrier == r.FormValue("APN") {
			http.Redirect(w, r, "/show/"+data[key].MCC+"-"+data[key].MNC, http.StatusTemporaryRedirect)
			return
		}
	}
}

// Shows a single APN
func show(w http.ResponseWriter, r *http.Request) {
	// Get the full identifier
	identifier := strings.Replace(r.URL.Path, "/show/", "", 1)
	split := strings.SplitAfter(identifier, "-")

	MCC := strings.Replace(split[0], "-", "", 1)
	MNC := split[1]
	data := getAPN(r, MCC, MNC)

	passedTemplate := new(bytes.Buffer)
	template.Must(template.ParseFiles("apnconf/templates/show.html")).Execute(passedTemplate, data)
	render.Render(w, r, passedTemplate)
}