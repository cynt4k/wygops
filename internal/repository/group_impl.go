package repository

import (
	"time"

	"github.com/cynt4k/wygops/internal/event"
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/pkg/util/gormutil"
	"github.com/cynt4k/wygops/pkg/util/structs"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
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

	msg, _ := structs.StructToMap(event.GroupCreatedEvent{
		GroupID: group.ID,
		Name:    group.Name,
	})

	repo.hub.Publish(hub.Message{
		Name:   event.GroupCreated,
		Fields: msg,
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
		msg, _ := structs.StructToMap(event.UserAddedToGroupEvent{
			UserID:  userID,
			GroupID: groupID,
		})
		repo.hub.Publish(hub.Message{
			Name:   event.UserAddedToGroup,
			Fields: msg,
		})
	}
	return nil
}

// RemoveUserFromGroup : Remove an user from an group
func (repo *GormRepository) RemoveUserFromGroup(userID uint, groupID uint) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var group models.Group
		if err := tx.Scopes(userGroupPreloads).First(&group, &models.Group{ID: groupID}).Error; err != nil {
			return convertError(err)
		}

		if group.IsMember(userID) {
			if err := tx.Delete(&models.UserGroup{UserID: userID, GroupID: groupID}).Error; err != nil {
				return err
			}
		}
		return tx.Model(&group).UpdateColumn("updated_at", time.Now()).Error
	})
	return err
}

// DeleteGroup : Remove an group
func (repo *GormRepository) DeleteGroup(groupID uint) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&models.UserGroup{GroupID: groupID}).Delete(&models.UserGroup{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.Group{ID: groupID})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
	if err != nil {
		return err
	}
	msg, _ := structs.StructToMap(event.GroupDeletedEvent{
		GroupID: groupID,
	})
	repo.hub.Publish(hub.Message{
		Name:   event.GroupDeleted,
		Fields: msg,
	})
	return nil
}

// GetGroup : Get the group
func (repo *GormRepository) GetGroup(groupID uint) (*models.Group, error) {
	var group models.Group
	if err := repo.db.Scopes(userGroupPreloads).First(&group, &models.Group{ID: groupID}).Error; err != nil {
		return nil, convertError(err)
	}
	return &group, nil
}

func userGroupPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Users")
}
