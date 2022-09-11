package repositories

import (
	voerrors "auth-app/internal/autherrors"
	"auth-app/internal/config"
	"auth-app/internal/model"
	"encoding/csv"
	"errors"
	"log"
	"os"
)

type UserRepository struct {
	cfg config.Config
}

func NewUserRepository(cfg config.Config) *UserRepository {
	return &UserRepository{
		cfg: cfg,
	}
}

func (r *UserRepository) SaveUser(user model.User) (err error) {

	file, err := os.OpenFile(r.cfg.UserDataFilePath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{user.Name, user.Phone, user.Role, user.Password, user.Timestampz})

	csvWriter.Flush()

	return nil
}

func (r *UserRepository) GetAllUsers() (users []model.User, err error) {

	// Opens the csv file
	file, err := os.Open(r.cfg.UserDataFilePath())
	if err != nil {

		return users, err
	}
	defer file.Close()

	// Read and parse the csv file into [][]string
	lines, _ := csv.NewReader(file).ReadAll()

	// Parse the result to new struct
	for _, line := range lines {
		user := model.User{
			Name:       line[0],
			Phone:      line[1],
			Role:       line[2],
			Password:   line[3],
			Timestampz: line[4],
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserByPhoneAndPassword(phone, password string) (user model.User, err error) {
	users, err := r.GetAllUsers()
	if err != nil {

		return user, errors.New("Failed to get all users")
	}

	for _, val := range users {
		if val.Phone == phone && val.Password == password {
			return val, nil
		}
	}

	return user, voerrors.ErrNotFound
}

func (r *UserRepository) GetUserByPhone(phone string) (user model.User, err error) {
	users, err := r.GetAllUsers()
	if err != nil {

		return user, errors.New("Failed to get all users")
	}

	for _, val := range users {
		if val.Phone == phone {
			return val, nil
		}
	}

	return user, voerrors.ErrNotFound
}
