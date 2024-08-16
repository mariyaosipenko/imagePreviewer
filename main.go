package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := mux.NewRouter()
	r.HandleFunc("/getimage/{width:[0-9]+}/{height:[0-9]+}/{url:[0-9a-zA-Z.\\/_-]+}", fillHandler).
		Methods("GET").
		Schemes("http")
	log.Fatal(http.ListenAndServe(":8005", r))
}

func fillHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	for k, v := range vars {
		log.Printf("%v: %v\n", k, v)
	}

	url := vars["url"]
	//height := vars["height"]
	//width := vars["width"]

	err := downloadFile(url)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := os.ReadFile("img.jpg")

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/jpg")
	w.Write(buf)
}

func downloadFile(url string) error {
	//Get the response bytes from the url
	response, err := http.Get("http://" + url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != 200 {
		log.Fatal("Received non 200 response code")
	}
	//Create an empty file
	file, err := os.Create("img.jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
