package webAscii

import (
	"fmt"
	"net/http"

	check "webAscii/checksum"
	print "webAscii/printAscii"
	output "webAscii/readWrite"
	send "webAscii/utils"
)

var banners = map[string]string{
	"standard":   "public/standard.txt",
	"thinkertoy": "public/thinkertoy.txt",
	"shadow":     "public/shadow.txt",
}

func AsciiServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		send.SendError(w, fmt.Sprintf("ParseForm() %v", err), http.StatusBadRequest)
		return
	}

	text := r.FormValue("Text")
	banner := r.FormValue("Banner")
	for param := range r.Form {
		if param != "Text" && param != "Banner" {
			send.SendError(w, "Error 404: Bad request", http.StatusBadRequest)
			break
		}
	}
	str := ""

	all := []string{"standard", "thinkertoy", "shadow"}

	if banner == "all" {
		for i, bn := range all {
			if i != 0 {
				str += "\n"
			}
			str += writeAscii(w, bn, text)
		}
	} else {
		str += writeAscii(w, banner, text)
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, str)
}

func writeAscii(w http.ResponseWriter, banner, text string) string {
	filename, ok := banners[banner]
	if !ok {
		send.SendError(w, "Error 404: Not Found: Invalid banner specified\n", http.StatusNotFound)
		return ""
	}

	err := check.ValidateFileChecksum(filename)
	if err != nil {
		send.SendError(w, fmt.Sprintf("Error 404: Error downloading or validating file: %v", err), http.StatusNotFound)
		return ""
	}

	asciiArtGrid, err := output.ReadAscii(filename, w)
	if err != nil {
		send.SendError(w, fmt.Sprintf("Error 500: Internal Server Error: Error reading ASCII art:%v", err), http.StatusInternalServerError)
		return ""
	}

	str := print.PrintArt(w, text, asciiArtGrid)
	return str
}
