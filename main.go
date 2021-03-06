package main

import (
	"bytes"
	context "context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/jesseokeya/go-httplogger"
	"github.com/joho/godotenv"
	grpc "google.golang.org/grpc"
)

func createCaso(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Se han enviado datos incorrectos para crear un caso")
		return
	}

	nuevos := &CasoRequest{}
	if err := jsonpb.Unmarshal(bytes.NewBuffer(reqBody), nuevos); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Se han enviado datos incorrectos para crear un caso")
		return
	}

	URL := getVariable("URL_GRPC")

	conn, err := grpc.Dial(URL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "No se pudo conectar al servidor GRPC")
		return
	}
	defer conn.Close()

	nuevoCliente := NewCasoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := nuevoCliente.CrearCasos(ctx, nuevos)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "No se ha podido enviar el caso.")
		return
	}
	message := response.GetMensaje()
	w.WriteHeader(http.StatusOK)
	log.Println(message)
	fmt.Fprintf(w, "Mensaje: %s", message)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bienvenido a mi API en Golang")
}

func getVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	return os.Getenv(key)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute).Methods("GET")
	router.HandleFunc("/", createCaso).Methods("POST")
	router.Use(mux.CORSMethodMiddleware(router))
	log.Println("Starting server. Listening on port 4000.")
	log.Fatal(http.ListenAndServe(":4000", httplogger.Golog(router)))
}
