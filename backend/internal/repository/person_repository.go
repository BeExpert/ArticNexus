package repository

import (
	"articnexus/backend/internal/domain"

	"gorm.io/gorm"
)

// PersonRepository defines the persistence contract for Person entities.
type PersonRepository interface {
	FindByID(id int64) (*domain.Person, error)
	Create(person *domain.Person) error
	Update(person *domain.Person) error
	Delete(id int64) error
}

type personRepository struct {
	db *gorm.DB
}

// NewPersonRepository returns a GORM-backed PersonRepository.
func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{db: db}
}

func (r *personRepository) FindByID(id int64) (*domain.Person, error) {
	var person domain.Person
	if err := r.db.First(&person, id).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *personRepository) Create(person *domain.Person) error {
	return r.db.Create(person).Error
}

func (r *personRepository) Update(person *domain.Person) error {
	return r.db.Save(person).Error
}

func (r *personRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Person{}, id).Error
}
