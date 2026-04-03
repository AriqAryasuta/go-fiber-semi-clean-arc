package user

import "backend-boiler/internal/shared/utils"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(input CreateUserRequest) (UserResponse, error) {
	model := &UserModel{
		ID:    utils.NewID(),
		Name:  input.Name,
		Email: input.Email,
	}
	if err := s.repository.Create(model); err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		ID:    model.ID,
		Name:  model.Name,
		Email: model.Email,
	}, nil
}

func (s *Service) List() ([]UserResponse, error) {
	models, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]UserResponse, 0, len(models))
	for _, item := range models {
		result = append(result, UserResponse{
			ID:    item.ID,
			Name:  item.Name,
			Email: item.Email,
		})
	}
	return result, nil
}
