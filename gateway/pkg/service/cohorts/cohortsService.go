package cohorts

import (
	"PingMeMaybe/libs/db/models"
	"github.com/hibiken/asynq"
)

type CohortsService struct {
	asynq             *asynq.Client
	cohortsRepository models.IUserCohortRepository
}

type CohortsServiceInterface interface {
	GetUserCohorts(userId int) ([]models.UserCohort, error)
}

func NewCohortsService(asynq *asynq.Client, cohortsRepository models.IUserCohortRepository) CohortsServiceInterface {
	return &CohortsService{
		asynq:             asynq,
		cohortsRepository: cohortsRepository,
	}
}

func (c *CohortsService) GetUserCohorts(userId int) ([]models.UserCohort, error) {
	return nil, nil
	// TODO: WIP
	//cohorts, err := c.cohortsRepository.GetCohortUsers()
	//if err != nil {
	//	return nil, err
	//}
	//return cohorts, nil
}
