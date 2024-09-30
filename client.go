package main

import (
	"context"
	"fmt"
	"log"

	petstore "<project>/petstore"
)

func main() {
	client, err := petstore.NewClient("http://localhost:8080/v3")
	if err != nil {
		log.Fatalf("Tạo client thất bại: %v", err)
	}

	res, err := client.GetPetById(context.Background(), petstore.GetPetByIdParams{
		PetId: 1,
	})
	if err != nil {
		log.Fatalf("Lấy pet thất bại: %v", err)
	}

	switch p := res.(type) {
	case *petstore.Pet:
		fmt.Printf("Pet tìm thấy: %+v\n", p)
	case *petstore.GetPetByIdNotFound:
		fmt.Println("Pet không tìm thấy")
	}
}
