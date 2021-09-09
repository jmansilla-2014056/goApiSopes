package main

import (
	"ApiGo/models"
	ps "ApiGo/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/publicar", publicar)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func publicar(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	var d models.Publicacion

	if json.NewDecoder(r.Body).Decode(&d) != nil {
		//println("CLIENTE: error 3")
	} else {
		//firefox id 10845 / ps -ef | grep firefox
		err := ps.Create(d)

		if err != nil{
			var x  = models.Mensaje{
				Msg : "Ocurrio un error insertando publicacion",
			}
			fmt.Println("El insert fracaso")
			json.NewEncoder(w).Encode(x)
		}else {
			var x  = models.Mensaje{
				Msg : "Se inserto publicacion con exito",
			}
			json.NewEncoder(w).Encode(x)
			fmt.Println(" El insert fue un exito")
		}

		fmt.Println("metodo publicar finalizado")

	}
}

