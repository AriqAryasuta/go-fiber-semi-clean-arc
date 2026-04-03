package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(model *UserModel) error {
	if r.db == nil {
		return nil
	}
	return r.db.Create(model).Error
}

func (r *Repository) FindAll() ([]UserModel, error) {
	if r.db == nil {
		return []UserModel{}, nil
	}
	var users []UserModel
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
