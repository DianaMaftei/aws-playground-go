package database

import (
	"github.com/stretchr/testify/require"
	"pet-store/model"
	"testing"
)

func TestFindPetsByStatus(t *testing.T) {
	pets, err := FindPetsByStatus(model.Available)

	require.NoError(t, err)
	require.NotEmpty(t, pets)
	require.Equal(t, 1, len(pets))
}

func TestFindPetByID(t *testing.T) {
	pet, err := FindPetByID(1)

	require.NoError(t, err)
	require.Equal(t, "Fluffy", pet.Name)
}

func TestCreatePet(t *testing.T) {
	pet := model.Pet{Name: "TestPet"}
	pet, err := CreatePet(pet)

	require.NoError(t, err)
	require.Equal(t, 2, pet.Id)
}

func TestUpdatePet(t *testing.T) {
	pet := model.Pet{Id: 1, Status: model.Sold}
	pet, err := UpdatePet(pet)

	require.NoError(t, err)

	pet, err = FindPetByID(1)
	require.NoError(t, err)
	require.Equal(t, model.Sold, pet.Status)
}

func TestDeletePet(t *testing.T) {
	err := DeletePet(1)

	require.NoError(t, err)

	pet, err := FindPetByID(1)

	require.Error(t, err)
	require.Empty(t, pet)
}
