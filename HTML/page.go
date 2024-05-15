package never

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	// "strings"
)

type HomePageData struct {
	Title string
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"dates"`
	Relations    string   `json:"datesLocations"`
}
type Location struct {
	ID        int      `json:"id"`
	LocationS []string `json:"locations"`
}
type Date struct {
	ID           int      `json:"id"`
	ConcertDates []string `json:"dates"`
}
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type ArtistWithInfo struct {
	Artist
	Locations []string            `json:"locations"`
	Dates     []string            `json:"dates"`
	Relations map[string][]string `json:"datesLocations"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		HandleNotFound(w,r)
		return
	}
	if r.Method != http.MethodGet{
		HandleMethod(w,r)
		return
	}
	// this is where we are getting the API info
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		HandleInternalError(w,r)

		return
	}
	defer resp.Body.Close()

	// this is where we are placing the decoded info from the API in the artist array
	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		HandleInternalError(w,r)
		return
	}

	// this is where it is recognizing the index.html
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleInternalError(w,r)
		return
	}

	// this is where it is executing the wanted data from artist array into the html index
	err = tmpl.Execute(w, artists)
	if err != nil {
		// http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func HandleRequest2(w http.ResponseWriter, r *http.Request) {
	// Extract the artist ID from the query parameters
	artistID := r.URL.Query().Get("id")

	// Convert the artistID to an integer
	id, err := strconv.Atoi(artistID)
	if err != nil {
		HandleNotFound(w,r)
		return
	}

	// Check if the ID is greater than 52
	if id > 52 || id <= 0 {
		HandleNotFound(w,r)
		return
	}

	// Fetch the artist's detailed information using the artist ID
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + artistID)
	if err != nil {
		HandleInternalError(w,r)
		return
	}
	defer resp.Body.Close()

	// Decode the artist's detailed information
	var artist Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		http.Error(w, "Failed to decode artist data", http.StatusBadGateway)
		return
	}

	// Get the location information
	location := getLocation(artist.ID, w, r)
	// Get the relation information
	relation := getRelation(artist.ID, w, r)
	// Get the relation information
	concertDate := getDates(artist.ID, w, r)

	// Create a new instance of ArtistWithInfo
	artistWithInfo := ArtistWithInfo{
		Artist:    artist,
		Locations: location.LocationS,
		Dates:     concertDate.ConcertDates,
		Relations: relation.DatesLocations,
	}

	// Load the info.html template
	tmpl, err := template.ParseFiles("templates/info.html")
	if err != nil {
		HandleInternalError(w,r)
		return
	}

	// Execute the template with the artist's detailed information
	err = tmpl.Execute(w, artistWithInfo)
	if err != nil {
		HandleInternalError(w,r)
		return
	}
}
