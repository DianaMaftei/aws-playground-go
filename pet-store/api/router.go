package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)

	// PET
	router.HandleFunc("/pet/findByStatus", FindPetsByStatus).Methods("GET")
	router.HandleFunc("/pet", CreatePet).Methods("POST")
	router.HandleFunc("/pet/{petId}", UpdatePet).Methods("PUT")
	router.HandleFunc("/pet/{petId}", DeletePet).Methods("DELETE")
	router.HandleFunc("/pet/{petId}", FindPetById).Methods("GET")
	router.HandleFunc("/pet/{petId}/uploadImage", UploadPetImage).Methods("POST")

	// STORE

	// USER

	return router
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}
