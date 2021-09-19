// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample pubsub demonstrates use of the cloud.google.com/go/pubsub package from App Engine flexible environment.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sopes/apigo/models"
	"strconv"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (
	topic *pubsub.Topic

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string

	// token is used to verify push requests.
	token = mustGetenv("GOOGLE_APPLICATION_CREDENTIALS")
)

const maxMessages = 10

func setupResponse2(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}


func main() {
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

	// Create the topic if it doesn't exist.
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Printf("Topic %v doesn't exist - creating it", topicName)
		_, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/pubsub/publish", publishHandler)
	http.HandleFunc("/pubsub/push", pushHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	// Verify the token.
	setupResponse2(&w, r)
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

func listHandler(w http.ResponseWriter, r *http.Request) {
	messagesMu.Lock()
	defer messagesMu.Unlock()

	if err := tmpl.Execute(w, messages); err != nil {
		log.Printf("Could not execute template: %v", err)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	msg := &pubsub.Message{
		Data: []byte(r.FormValue("payload")),
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		http.Error(w, fmt.Sprintf("Could not publish message: %v", err), 500)
		return
	}

	fmt.Fprint(w, "Message published.")
}

var tmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>Pub/Sub</title>
  </head>
  <body>
    <div>
      <p>Last ten messages received by this instance:</p>
      <ul>
      {{ range . }}
          <li>{{ . }}</li>
      {{ end }}
      </ul>
    </div>
    <form method="post" action="/pubsub/publish">
      <textarea name="payload" placeholder="Enter message here"></textarea>
      <input type="submit">
    </form>
    <p>Note: if the application is running across multiple instances, each
      instance will have its own list of messages.</p>
  </body>
</html>`))