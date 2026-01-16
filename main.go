package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// --- Structs for Specific JSON Format (New) ---

type ExperienceItem struct {
	Role     string   `json:"Role"`
	Company  string   `json:"Company"`
	Location string   `json:"Location"`
	Link     string   `json:"Link"`
	Period   string   `json:"Period"`
	Points   []string `json:"Points"`
}

type EducationItem struct {
	Degree      string `json:"Degree"`
	Institution string `json:"Institution"`
	Location    string `json:"Location"`
	Period      string `json:"Period"`
	GPA         string `json:"GPA"`
}

// --- Structs for Generic/Legacy Format (Backward Compatible) ---

type Entry struct {
	Title       string   `json:"title,omitempty"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Period      string   `json:"period,omitempty"`
	Location    string   `json:"location,omitempty"`
	Description []string `json:"description,omitempty"`
}

type Section struct {
	Title   string  `json:"title"`
	Entries []Entry `json:"entries"`
}

// --- Main Resume Struct ---

type Resume struct {
	// Header Info
	Name     string `json:"Name"`
	Title    string `json:"Title"` // Added
	Phone    string `json:"Phone"`
	Email    string `json:"Email"`
	Location string `json:"Location"`
	Linkedin string `json:"Linkedin"`
	Github   string `json:"Github"`
	Website  string `json:"Website"` // Added
	Summary  string `json:"Summary"`

	// Specific Arrays (Matches your JSON)
	Skills     []string         `json:"Skills"`
	Experience []ExperienceItem `json:"Experience"`
	Education  []EducationItem  `json:"Education"`

	// Generic Arrays (Backward Compatibility)
	Sections []Section `json:"sections"`
}

const OUTPUT_DIR = "./output"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run main.go <path-to-json> [path-to-template]\n")
		return
	}

	templ_path := "./templates/default.html"
	if len(os.Args) >= 3 {
		templ_path = os.Args[2]
		fmt.Println("Using template:", templ_path)
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(OUTPUT_DIR); os.IsNotExist(err) {
		os.Mkdir(OUTPUT_DIR, 0755)
	}

	// Parse Template
	t, err := template.ParseFiles(templ_path)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	out_path := filepath.Join(OUTPUT_DIR, "resume.html")
	out_file, err := os.Create(out_path)
	if err != nil {
		panic(err)
	}
	defer out_file.Close()

	data := parse_resume(os.Args[1])

	err = t.Execute(out_file, data)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Resume has been generated to:", out_path)
}

func parse_resume(path string) Resume {
	var resume Resume
	byteValue, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading json file: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(byteValue, &resume)
	if err != nil {
		fmt.Printf("Error unmarshalling json: %v\n", err)
		os.Exit(1)
	}
	return resume
}
