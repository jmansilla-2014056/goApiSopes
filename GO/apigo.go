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
	"strconv"
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

	if json.NewDecoder(r.Body).Decode(&d) != nil {
		var x = models.Mensaje{
			Api:    "go",
			Result: false,
		}
		json.NewEncoder(w).Encode(x)
		fmt.Println(" El insert fracaso")
	} else {
		//firefox id 10845 / ps -ef | grep firefox
		if len(d.Nombre) > 0 && len(d.Fecha) > 0 && len(d.Comentario) > 0 && d.Downvotes >= 0 &&
			d.Upvotes >= 0 {

			err2 := ps.CreateS(d)
			err := ps.CreateM(d)

			if err2 != nil{
				var x  = models.Mensaje{
					Api : "go",
					Result: false,
				}
				fmt.Println("El insert fracaso")
				json.NewEncoder(w).Encode(x)
			}

			if err != nil || err2 != nil{
				var x  = models.Mensaje{
					Api : "go",
					Result: false,
				}
				fmt.Println("El insert fracaso")
				json.NewEncoder(w).Encode(x)
			}else {
				var x  = models.Mensaje{
					Api : "go",
					Result: true,
				}
				json.NewEncoder(w).Encode(x)
				fmt.Println(" El insert fue un exito")
			}
		} else{
			var x = models.Mensaje{
				Api:    "go",
				Result: false,
			}
			json.NewEncoder(w).Encode(x)
			fmt.Println(" El insert fracaso")
		}
		fmt.Println("metodo publicar finalizado")
	}
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
			Api:    "go",
			Estado: "Error",
			Numero: 0,
		}
		json.NewEncoder(w).Encode(x)
		fmt.Println(" El push fracaso")
	}

	var y = "{\"api\": \"" + d.Api + "\",\"estado\":\"" + d.Estado + "\",\"numero\":" + strconv.Itoa(d.Numero) + "\"}"

	msg := &pubsub.Message{
		Data: []byte(y),
	}

	result := topic.Publish(ctx, msg)

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Print("error")
		fmt.Print(err)
		var x = models.Notificacion{
			Api:    "go",
			Estado: "Error",
			Numero: 0,
		}
		json.NewEncoder(w).Encode(x)
		fmt.Println(" El push fracaso")
	} else {
		fmt.Print("Publicado: %v", id)
	}

	json.NewEncoder(w).Encode(d)

}
