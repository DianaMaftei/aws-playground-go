package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"pet-store/database"
	"pet-store/model"
	"strconv"
)

func CreatePet(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var pet model.Pet
	err := json.Unmarshal(reqBody, &pet)
	handleError(w, err, "Invalid Request Body", http.StatusBadRequest)

	pet, err = database.CreatePet(pet)
	if handleError(w, err, "Unable to save pet", http.StatusInternalServerError) {
		return
	}
	json.NewEncoder(w).Encode(pet)
}

func UpdatePet(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var pet model.Pet
	err := json.Unmarshal(reqBody, &pet)
	handleError(w, err, "Invalid Request Body", http.StatusBadRequest)

	pet, err = database.UpdatePet(pet)
	if handleError(w, err, "Unable to update pet", http.StatusInternalServerError) {
		return
	}
	json.NewEncoder(w).Encode(pet)
}

func FindPetsByStatus(w http.ResponseWriter, r *http.Request) {
	status := getRequestParam(r, "status")
	petStatus := model.PetStatus(status)
	err := validatePetStatus(petStatus)
	if handleError(w, err, "Invalid pet status supplied", http.StatusBadRequest) {
		return
	}

	pets, err := database.FindPetsByStatus(petStatus)
	if handleError(w, err, "Unable to fetch pets by status", http.StatusInternalServerError) {
		return
	}
	json.NewEncoder(w).Encode(pets)
}

func FindPetById(w http.ResponseWriter, r *http.Request) {
	key := getRequestParam(r, "petId")
	id, err := strconv.Atoi(key)

	if handleError(w, err, "Invalid ID supplied", http.StatusBadRequest) {
		return
	}

	pet, err := database.FindPetByID(id)

	if handleError(w, err, "Pet not found", http.StatusNotFound) {
		return
	}

	json.NewEncoder(w).Encode(pet)
}

func DeletePet(w http.ResponseWriter, r *http.Request) {
	key := getRequestParam(r, "petId")
	id, err := strconv.Atoi(key)

	if handleError(w, err, "Invalid ID supplied", http.StatusBadRequest) {
		return
	}

	err = database.DeletePet(id)
	handleError(w, err, "Pet not found", http.StatusNotFound)
}

// TODO
// update with PUT
// findByStatus - GET
// uploadImage

func validatePetStatus(status model.PetStatus) error {
	validStatuses := []model.PetStatus{model.Available, model.Pending, model.Sold}
	for _, val := range validStatuses {
		if val == status {
			return nil
		}
	}
	return errors.New("invalid pet status")
}

func getRequestParam(r *http.Request, param string) string {
	vars := mux.Vars(r)
	key := vars[param]
	return key
}

func handleError(w http.ResponseWriter, err error, message string, status int) bool {
	if err != nil {
		w.Write([]byte(message))
		w.WriteHeader(status)
		return true
	}
	return false
}
