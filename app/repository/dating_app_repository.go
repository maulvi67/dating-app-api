package repository

import (
	"dating-apps/app/model/entity"
	db "dating-apps/helper/database"
	"time"
)

type DatingAppRepository interface {
	CreateUser(data entity.User) (*entity.User, error)
	Swipe(data entity.Swipe) (*entity.Swipe, error)
	CountUser(email string) (int64, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CheckSwipedUser(userID uint, targetUserID uint, today time.Time) (*entity.Swipe, error)
	CountSwipedUser(userID uint, today time.Time) (int64, error)
	SaveUser(data entity.User) (*entity.User, error)
}

type datingAppRepository struct {
	BaseRepository
}

func NewDatingAppRepository(db *db.Database) DatingAppRepository {
	return &datingAppRepository{NewBaseRepository(db)}
}

func (r *datingAppRepository) CreateUser(data entity.User) (*entity.User, error) {
	err := r.GetDB().Create(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *datingAppRepository) Swipe(data entity.Swipe) (*entity.Swipe, error) {
	err := r.GetDB().Create(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *datingAppRepository) CountUser(email string) (int64, error) {
	var count int64
	err := r.GetDB().Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *datingAppRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *datingAppRepository) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.GetDB().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *datingAppRepository) CheckSwipedUser(userID uint, targetUserID uint, today time.Time) (*entity.Swipe, error) {
	var existing entity.Swipe
	err := r.GetDB().Where("user_id = ? AND target_user_id = ? AND swipe_date = ?", userID, targetUserID, today).
		First(&existing).Error
	if err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *datingAppRepository) CountSwipedUser(userID uint, today time.Time) (int64, error) {
	var count int64
	err := r.GetDB().Model(&entity.Swipe{}).Where("user_id = ? AND swipe_date = ?", userID, today).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *datingAppRepository) SaveUser(data entity.User) (*entity.User, error) {
	err := r.GetDB().Save(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
