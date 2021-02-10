package sync

import (
	"errors"
	"fmt"

	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/pkg/util/array"
)

type ldapSyncJob struct {
	sync *Service
}

func (s *Service) startLdap() error {
	if s.ldap == nil {
		return fmt.Errorf("ldap sync not initialized - check the initialization")
	}

	go s.runTimer(ldapSync)
	return nil
}

func ldapSync(s *Service) {
	job := &ldapSyncJob{
		sync: s,
	}

	err := job.syncUser()

	if err != nil {
		s.logger.Error(fmt.Sprintf("error while syncing user %s", err))
		return
	}
}

func (s *ldapSyncJob) syncUser() error {
	dbUsers, err := s.sync.repo.GetLdapUsers()

	if err != nil {
		return err
	}

	for _, user := range *dbUsers {
		ldapUser, err := s.sync.ldap.FindUser(user.Username, true)
		if err != nil {
			return err
		}

		if ldapUser == nil {
			err = s.sync.repo.DeleteUser(user.ID)
			if err != nil {
				return err
			}
			continue
		}

		for _, group := range ldapUser.Groups {
			dbGroup, err := s.sync.repo.GetGroupByName(group)
			if err != nil {
				if !errors.Is(err, repository.ErrNotFound) {
					return err
				}
			}

			if dbGroup == nil {
				newGroup := &models.Group{
					Name: group,
					Type: "ldap",
				}
				dbGroup, err = s.sync.repo.CreateGroup(newGroup)

				if err != nil {
					return err
				}
			}

			err = s.sync.repo.AddUserToGroup(user.ID, dbGroup.ID)

			if err != nil {
				return err
			}
		}

		dbGroups, err := s.sync.repo.GetGroupsByUser(user.ID)

		if err != nil {
			return err
		}

		for _, group := range *dbGroups {
			ldapGroup, err := s.sync.ldap.GetGroup(group.Name, false)

			if err != nil {
				return err
			}

			if ldapGroup == nil {
				err = s.sync.repo.DeleteGroup(group.ID)

				if err != nil {
					return err
				}
				continue
			}

			if !array.ContainsEntry(ldapUser.Groups, group.Name) {
				err = s.sync.repo.RemoveUserFromGroup(user.ID, group.ID)

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
