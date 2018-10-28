package types

import (
	"strconv"
	"strings"
)

// Ingredient domain
type Ingredient struct {
	OID         string `bson:"_id"`
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Quantity    int    `bson:"quantity"`
}

// GetOID function
func (i *Ingredient) GetOID() string {
	return i.OID
}

// SetOID function
func (i *Ingredient) SetOID(id string) {
	i.OID = id
}

// GetID function
func (i *Ingredient) GetID() string {
	return i.ID
}

// SetID function
func (i *Ingredient) SetID(id string) {
	i.ID = id
}

// GetName function
func (i *Ingredient) GetName() string {
	return i.Name
}

// SetName function
func (i *Ingredient) SetName(name string) {
	i.Name = name
}

// GetDescription function
func (i *Ingredient) GetDescription() string {
	return i.Description
}

// SetDescription function
func (i *Ingredient) SetDescription(description string) {
	i.Description = description
}

// GetQuantity function
func (i *Ingredient) GetQuantity() int {
	return i.Quantity
}

// SetQuantity function
func (i *Ingredient) SetQuantity(quantity int) {
	i.Quantity = quantity
}

// GetObjectInfo function
func (i *Ingredient) GetObjectInfo() string {
	info := []string{
		i.Name,
		i.Description,
	}
	resp := strings.Join(info, ": ")
	quantity := "\nQuantity: " + strconv.Itoa(i.Quantity)
	resp += quantity
	return resp
}
