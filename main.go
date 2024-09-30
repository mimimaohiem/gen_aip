package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	petstore "project_1/petstore"
)

// petsService là cấu trúc để giữ dữ liệu về pets.
type petsService struct {
	pets map[int64]petstore.Pet
	id   int64
	mux  sync.Mutex
}

// AddPet thêm một pet mới vào store.
func (p *petsService) AddPet(ctx context.Context, req *petstore.Pet) (*petstore.Pet, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Tăng ID và gán ID cho pet mới dưới dạng OptInt64.
	p.id++
	req.ID = petstore.NewOptInt64(p.id)
	p.pets[p.id] = *req
	return req, nil
}

// DeletePet xóa một pet bằng ID.
func (p *petsService) DeletePet(ctx context.Context, params petstore.DeletePetParams) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	delete(p.pets, params.PetId)
	return nil
}

// GetPetById trả về pet theo ID.
func (p *petsService) GetPetById(ctx context.Context, params petstore.GetPetByIdParams) (petstore.GetPetByIdRes, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	pet, ok := p.pets[params.PetId]
	if !ok {
		// Return Not Found nếu không có pet.
		return &petstore.GetPetByIdNotFound{}, nil
	}
	return &pet, nil
}

// UpdatePet cập nhật thông tin của pet.
func (p *petsService) UpdatePet(ctx context.Context, params petstore.UpdatePetParams) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Kiểm tra xem pet có tồn tại không.
	pet, ok := p.pets[params.PetId]
	if !ok {
		return fmt.Errorf("Pet not found")
	}

	// Cập nhật trạng thái và tên (nếu có).
	pet.Status = params.Status
	if val, ok := params.Name.Get(); ok {
		pet.Name = val
	}
	p.pets[params.PetId] = pet

	return nil
}

func main() {
	// Tạo một service instance.
	service := &petsService{
		pets: make(map[int64]petstore.Pet),
	}
	// Tạo server từ code được sinh ra.
	srv, err := petstore.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	// Chạy server tại cổng 8080.
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
