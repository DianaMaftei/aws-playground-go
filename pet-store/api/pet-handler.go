package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"pet-store/database"
	"pet-store/model"
	"pet-store/service"
	"strconv"
)

func CreatePet(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var pet model.Pet
	err := json.Unmarshal(reqBody, &pet)
	if handleError(w, err, "invalid request body", http.StatusBadRequest) {
		return
	}

	pet, err = database.CreatePet(pet)
	if handleError(w, err, "unable to save pet", http.StatusInternalServerError) {
		return
	}
	json.NewEncoder(w).Encode(pet)
}

func UpdatePet(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var pet model.Pet
	err := json.Unmarshal(reqBody, &pet)
	if handleError(w, err, "invalid request body", http.StatusBadRequest) {
		return
	}

	pet, err = database.UpdatePet(pet)
	if handleError(w, err, "unable to update pet", http.StatusInternalServerError) {
		return
	}
	json.NewEncoder(w).Encode(pet)
}

func FindPetsByStatus(w http.ResponseWriter, r *http.Request) {
	status := getRequestParam(r, "status")
	petStatus := model.PetStatus(status)
	err := validatePetStatus(petStatus)
	if handleError(w, err, "invalid pet status supplied", http.StatusBadRequest) {
		return
	}

	pets, err := database.FindPetsByStatus(petStatus)
	if handleError(w, err, "unable to fetch pets by status", http.StatusInternalServerError) {
		return
	}
	json.NewEncoder(w).Encode(pets)
}

func FindPetById(w http.ResponseWriter, r *http.Request) {
	key := getRequestParam(r, "petId")
	id, err := strconv.Atoi(key)
	if handleError(w, err, "invalid petId supplied", http.StatusBadRequest) {
		return
	}

	pet, err := database.FindPetByID(id)
	if handleError(w, err, "pet not found", http.StatusNotFound) {
		return
	}

	json.NewEncoder(w).Encode(pet)
}

func DeletePet(w http.ResponseWriter, r *http.Request) {
	key := getRequestParam(r, "petId")
	id, err := strconv.Atoi(key)
	if handleError(w, err, "invalid petId supplied", http.StatusBadRequest) {
		return
	}

	err = database.DeletePet(id)
	handleError(w, err, "pet not found", http.StatusNotFound)
}

func UploadPetImage(w http.ResponseWriter, r *http.Request) {
	maxSize := int64(1024000) // allow only 1MB of file size
	err := r.ParseMultipartForm(maxSize)
	if handleError(w, err, fmt.Sprintf("image too large. Max Size: %v", maxSize), http.StatusBadRequest) {
		return
	}

	file, fileHeader, err := r.FormFile("pet_image")
	if handleError(w, err, fmt.Sprintf("invalid form file pet_image"), http.StatusBadRequest) {
		return
	}
	defer file.Close()

	session, err := service.ConnectToAws()
	if handleError(w, err, fmt.Sprintf("could not upload file"), http.StatusInternalServerError) {
		return
	}

	fileName, err := service.UploadFileToS3(session, file, fileHeader, "pictures")
	if handleError(w, err, fmt.Sprintf("could not upload file"), http.StatusInternalServerError) {
		return
	}

	fmt.Printf("Image uploaded successfully: %v", fileName)
}

func validatePetStatus(status model.PetStatus) error {
	validStatuses := []model.PetStatus{model.Available, model.Pending, model.Sold}
	for _, val := range validStatuses {
		if val == status {
			return nil
		}
	}
	return fmt.Errorf("invalid pet status %s", status)
}

func getRequestParam(r *http.Request, param string) string {
	vars := mux.Vars(r)
	key := vars[param]
	return key
}

func handleError(w http.ResponseWriter, err error, message string, status int) bool {
	if err != nil {
		fmt.Printf(message, err.Error())
		w.Write([]byte(message))
		w.WriteHeader(status)
		return true
	}
	return false
}
