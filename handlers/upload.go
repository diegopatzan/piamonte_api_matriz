package handlers

import (
	"fmt"
	"log"
	"net/http"

	"piamonte-proxy/config"
	"piamonte-proxy/sftpclient"
)

func UploadHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed. Please use POST.", http.StatusMethodNotAllowed)
			return
		}

		// aguanta hasta 10MB en memoria por si odoo manda un doc pesado
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing multipart form data.", http.StatusBadRequest)
			return
		}

		// extraemos el doc que mandaron
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File payload not found in request.", http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Received file from Odoo: %s (%d bytes)", header.Filename, header.Size)

		err = sftpclient.ForwardToSFTP(cfg, file, header.Filename)
		if err != nil {
			log.Printf("SFTP transfer failed: %v", err)
			http.Error(w, "Internal server error during remote transfer.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success: File %s was uploaded to Piamonte.", header.Filename)
	}
}
