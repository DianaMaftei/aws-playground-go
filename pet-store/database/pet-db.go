package database

import (
	"pet-store/model"
)

func FindPetsByStatus(status model.PetStatus) ([]model.Pet, error) {
	var pets []model.Pet
	results, err := DBCon.Query("SELECT * FROM pet where status = ?", status)
	if err != nil {
		return pets, err
	}
	for results.Next() {
		var pet model.Pet
		err = results.Scan(&pet.Id, &pet.Name, &pet.Status)
		if err != nil {
			return pets, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func FindPetByID(id int) (model.Pet, error) {
	var pet model.Pet
	row := DBCon.QueryRow("SELECT * FROM pet WHERE id = ?", id)
	if err := row.Scan(&pet.Id, &pet.Name, &pet.Status); err != nil {
		return pet, err
	}
	return pet, nil
}

func CreatePet(pet model.Pet) (model.Pet, error) {
	res, err := DBCon.Exec("INSERT INTO pet (name, status) VALUES (?, ?)", pet.Name, model.Available)
	if err != nil {
		return pet, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return pet, err
	}
	pet.Id = int(id)
	return pet, nil
}

func UpdatePet(pet model.Pet) (model.Pet, error) {
	if pet.Id == 0 {
		return CreatePet(pet)
	}

	res, err := DBCon.Exec("UPDATE pet SET  name = ?, status = ? WHERE id = ?;", pet.Name, pet.Status, pet.Id)
	if err != nil {
		return pet, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return pet, err
	}
	pet.Id = int(id)
	return pet, nil
}

func DeletePet(id int) error {
	_, err := DBCon.Exec("DELETE FROM pet WHERE id = ?", id)
	return err
}
