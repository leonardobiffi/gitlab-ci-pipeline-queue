package queue

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Run(ref string)
}

type service struct {
	logger *logrus.Logger
}

var _ Service = (*service)(nil)

func (s *service) Run(ref string) {
	env := s.getEnv()

	// loop until the current pipeline is the oldest one
	for {
		pipelines, err := s.fetchPipelines(env.Token, env.ProjectID, ref)
		if err != nil {
			panic(err)
		}

		if len(pipelines) <= 1 {
			s.logger.Println("No other pipelines in queue, ready to continue!")
			return
		}

		s.logger.Debugf("Found %d pipelines in queue", len(pipelines))
		s.logger.Debugf("Latest pipeline ID: %d", pipelines[0].ID)

		// check if the current pipeline is the oldest one
		if pipelines[0].ID == env.PipelineID {
			s.logger.Println("The current pipeline is the oldest one, ready to continue!")
			return
		}

		s.logger.Println("The current pipeline is not the oldest one, let's wait for 5 seconds and retry")
		time.Sleep(5 * time.Second)
	}
}

func New() Service {
	return &service{
		logger: Logger(),
	}
}
