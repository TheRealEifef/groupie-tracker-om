package never

import (
	"encoding/json"
	"html/template"
	"net/http"
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
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// this is where we are getting the API info
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// this is where we are placing the decoded info from the API in the artist array
	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		http.Error(w, "Failed to decode artist data", http.StatusInternalServerError)
		return
	}

	// this is where it is recognizing the imdex.html
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
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

	// Fetch the artist's detailed information using the artist ID
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + artistID)
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	resploc, err2 := http.Get("https://groupietrackers.herokuapp.com/api/locations" + artistID)
	if err2 != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}
	defer resploc.Body.Close()

	respdate, err3 := http.Get("https://groupietrackers.herokuapp.com/api/dates" + artistID)
	if err3 != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}
	defer respdate.Body.Close()

	// Decode the artist's detailed information
	var artist Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		http.Error(w, "Failed to decode artist data", http.StatusInternalServerError)
		return
	}

	// Load the info.html template
	tmpl, err := template.ParseFiles("templates/info.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the artist's detailed information
	err = tmpl.Execute(w, artist)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

}
