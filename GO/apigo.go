package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sopes/apigo/models"
	ps "sopes/apigo/services"
	"sync"
)

var (
	topic *pubsub.Topic

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string

)


func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/publicar", publicar)
	http.HandleFunc("/finalizarCarga", pushH)
	http.HandleFunc("/iniciarCarga", iniciarCarga)
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

func iniciarCarga(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/iniciarCarga" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Carga Iniciada")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}


func publicar(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	var d models.Publicacion

	var x models.Mensaje
	x.Api = "go"
	x.Cosmos = false
	x.Sql = false

	if json.NewDecoder(r.Body).Decode(&d) != nil {
		x.Cosmos = false
		x.Sql = false
		json.NewEncoder(w).Encode(x)
		fmt.Println(" El insert fracaso")
	} else {
		//firefox id 10845 / ps -ef | grep firefox
		if len(d.Nombre) > 0 && len(d.Fecha) > 0 && len(d.Comentario) > 0 && d.Downvotes >= 0 &&
			d.Upvotes >= 0 {

			err2 := ps.CreateS(d)
			err := ps.CreateM(d)

			if err2 != nil{
				x.Sql = false
				fmt.Println("El insert SQL fracaso")
			}else{
				x.Sql = true
				fmt.Println(" El insert SQL fue un exito")
			}
			if err != nil {
				x.Cosmos = false
				fmt.Println("El insert Cosmos fracaso")
			}else {
				x.Cosmos = true
				fmt.Println(" El insert Cosmos fue un exito")
			}
		} else{
				x.Sql = false
				x.Cosmos = false
			json.NewEncoder(w).Encode(x)
			fmt.Println(" El insert fracaso")
		}
		fmt.Println("metodo publicar finalizado")
	}
	json.NewEncoder(w).Encode(x)
}



func pushH(w http.ResponseWriter, r *http.Request) {
	// Verify the token.
	setupResponse(&w, r)
	var d models.Notificacion

	ctx := context.Background()
	projectID := "savvy-hull-325303"

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	topicName := "sopes1p1"
	topic = client.Topic(topicName)

	if err != nil {
		fmt.Print("Error :( ")
		fmt.Print(err)
		http.Error(w, fmt.Sprintf("Could not publish message: %v", err), 500)
		return
	}

	if json.NewDecoder(r.Body).Decode(&d) != nil {
		var x = models.Notificacion{
			Go_cosmos:    0,
			Go_sql: 0,
			Python_cosmos: 0,
			Python_sql: 0,
			Rust_cosmos: 0,
			Rust_sql: 0,
		}
		json.NewEncoder(w).Encode(x)
		fmt.Println(" El push fracaso")
	}


	p, err2 := json.Marshal(d)

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	var y = string(p)

	msg := &pubsub.Message{
		Data: []byte(y),
	}

	result := topic.Publish(ctx, msg)

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Print("error")
		fmt.Print(err)
		var x = models.Notificacion{
			Go_cosmos:    0,
			Go_sql: 0,
			Python_cosmos: 0,
			Python_sql: 0,
			Rust_cosmos: 0,
			Rust_sql: 0,
		}
		json.NewEncoder(w).Encode(x)
		fmt.Println(" El push fracaso")
	} else {
		fmt.Println("Publicado: %v", id)
	}

	json.NewEncoder(w).Encode(d)

}
