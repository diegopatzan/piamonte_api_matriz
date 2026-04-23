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

		// Forzar a quitar el .txt si el sistema de origen lo agrega por defecto
		filename := header.Filename
		if len(filename) > 4 && filename[len(filename)-4:] == ".txt" {
			filename = filename[:len(filename)-4]
		}

		log.Printf("Received file from Odoo: %s (%d bytes) - Will save as: %s", header.Filename, header.Size, filename)

		err = sftpclient.ForwardToSFTP(cfg, file, filename)
		if err != nil {
			log.Printf("SFTP transfer failed: %v", err)
			http.Error(w, "Internal server error during remote transfer.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success: File %s was uploaded to Piamonte.", header.Filename)
	}
}
