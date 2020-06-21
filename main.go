package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jesseokeya/go-httplogger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

// Caso : Estructura para guardar un nuevo caso
type Caso struct {
	Nombre        string `json:"nombre"`
	Departamento  string `json:"departamento"`
	Edad          int32  `json:"edad"`
	FormaContagio string `json:"formaContagio"`
	Estado        string `json:"estado"`
}

func createCaso(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Se han enviado datos incorrectos para crear un caso")
		return
	}

	HOST := getVariable("HOST_GRPC")

	conn, err := grpc.Dial(HOST, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "No se pudo conectar al servidor GRPC")
		return
	}
	defer conn.Close()

	var nuevoCaso Caso
	json.Unmarshal(reqBody, &nuevoCaso)

	nuevoCliente := NewCasoClient(conn)

	response, err := nuevoCliente.CrearCaso(
		context.Background(), &CasoRequest{
			Nombre:        nuevoCaso.Nombre,
			Departamento:  nuevoCaso.Departamento,
			Edad:          nuevoCaso.Edad,
			FormaContagio: nuevoCaso.FormaContagio,
			Estado:        nuevoCaso.Estado})

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "No se ha podido enviar el caso.")
		return
	}
	message := response.GetMensaje()
	log.Println(message)
	fmt.Fprintf(w, "Mensaje: %s", message)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi API")
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
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/casos", createCaso).Methods("POST")
	log.Println("Starting server. Listening on port 4000.")
	log.Fatal(http.ListenAndServe(":4000", httplogger.Golog(router)))
}
