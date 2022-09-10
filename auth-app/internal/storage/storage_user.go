package storage

import "auth-app/internal/model"

func (s *Store) CreateUser(user model.User) error {
	return s.userRepository.SaveUser(user)
}

func (s *Store) FetchAllUsers() (users []model.User, err error) {
	return s.userRepository.GetAllUsers()
}

func (s *Store) FetchUserByPhoneAndPassword(phone, password string) (user model.User, err error) {
	return s.userRepository.GetUserByPhoneAndPassword(phone, password)
}
