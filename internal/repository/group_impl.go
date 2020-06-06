package repository

import (
	"time"

	"github.com/cynt4k/wygops/internal/event"
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/internal/util/gormutil"
	"github.com/jinzhu/gorm"
)

// CreateGroup : Create an group
func (repo *GormRepository) CreateGroup(group *models.Group) (*models.Group, error) {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if exist, err := gormutil.RecordExists(tx, &group); err != nil {
			return err
		} else if exist {
			return ErrAlreadyExists
		}

		if err := tx.Create(group).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	repo.bus.Publish(event.GroupCreated, event.GroupCreatedEvent{
		GroupID: group.ID,
		Name:    group.Name,
	})

	return group, nil
}

// AddUserToGroup : Add user to an group
func (repo *GormRepository) AddUserToGroup(userID uint, groupID uint) error {
	_, err := repo.GetUser(userID)

	if err != nil {
		return err
	}

	var added bool
	err = repo.db.Transaction(func(tx *gorm.DB) error {
		var group models.Group
		if err := tx.Preload("Users").First(&group, &models.Group{ID: groupID}).Error; err != nil {
			return convertError(err)
		}

		if !group.IsMember(userID) {
			if err := tx.Create(&models.UserGroup{UserID: userID, GroupID: groupID}).Error; err != nil {
				return err
			}
			added = true
			return tx.Model(&group).UpdateColumn("updated_at", time.Now()).Error
		}
		return nil
	})

	if err != nil {
		return err
	}

	if added {
		repo.bus.Publish(event.UserAddedToGroup, event.UserAddedToGroupEvent{
			UserID:  userID,
			GroupID: groupID,
		})
	}
	return nil
}

// GetGroup : Get the group
func (repo *GormRepository) GetGroup(groupID uint) (*models.Group, error) {
	group := models.Group{ID: groupID}
	repo.db.Find(&group)
	return &group, nil
}
