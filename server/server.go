package webAscii

import (
	"fmt"
	"log"
	"net/http"

	check "webAscii/checksum"
	print "webAscii/printAscii"
	output "webAscii/readWrite"
)

var banners = map[string]string{
	"standard":   "standard.txt",
	"thinkertoy": "thinkertoy.txt",
	"shadow":     "shadow.txt",
}

func AsciiServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() %v", err), http.StatusBadRequest)
	}

	fmt.Fprintf(w, "POST request successful\n")
	input := r.FormValue("Input")
	banner := r.FormValue("Banner")

	filename, ok := banners[banner]
	if !ok {
		http.Error(w, "Invalid banner specified", http.StatusBadRequest)
		return
	}

	err := check.ValidateFileChecksum(filename)
	if err != nil {
		log.Printf("Error downloading or validating file: %v", err)
		http.Error(w, "Error generating ASCII art", http.StatusInternalServerError)
		return
	}

	asciiArtGrid, err := output.ReadAscii(filename)
	if err != nil {
		log.Fatalf("Error reading ASCII map: %v", err)
		http.Error(w, "Error generating ASCII art", http.StatusInternalServerError)
	}

	str := print.PrintArt(input, asciiArtGrid)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, str)
}
