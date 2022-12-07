package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"
)

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/process", handleSubmit)
	http.HandleFunc("/upload", handleUpload)

	fmt.Println("server running on port 3000")
	http.ListenAndServe(":3000", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles(path.Join("views", "index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		tmpl, err := template.ParseFiles(path.Join("views", "success.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := fmt.Sprintf("Hello, %s your email is %s", name, email)

		data := map[string]string{
			"name":    name,
			"message": message,
		}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fmt.Printf("File Name: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Type: %+v\n", handler.Header)

		tempFile, err := ioutil.TempFile("temp-images", "upload-*.jpg")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := tempFile.Write(fileBytes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "File uploaded successfully!")
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
