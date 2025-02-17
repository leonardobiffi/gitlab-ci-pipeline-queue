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

	// retry counter for no pipeline in queue with no pipelines found
	var retriesWithNoPipeline int

	// loop until the current pipeline is the oldest one
	for {
		pipelines, err := s.fetchPipelines(env.Token, env.ProjectID, ref)
		if err != nil {
			s.logger.Fatalf("Failed to fetch pipelines: %v", err)
		}

		if len(pipelines) <= 1 {
			if retriesWithNoPipeline >= 3 {
				s.logger.Printf("No other pipelines in queue after 3 retries, continue to run pipeline...")
				return
			}

			retriesWithNoPipeline++
			s.logger.Printf("No other pipelines in queue, retrying in 5 seconds... (retry %d)", retriesWithNoPipeline)
			time.Sleep(5 * time.Second)

			continue
		}

		s.logger.Debugf("Found %d pipelines in queue", len(pipelines))
		s.logger.Debugf("Old pipeline ID: %d", pipelines[0].ID)

		// check if the current pipeline is the oldest one
		if pipelines[0].ID == env.PipelineID {
			s.logger.Println("The current pipeline is the oldest one, ready to continue!")
			return
		}

		s.logger.Println("The current pipeline is not the oldest one, let's wait for 10 seconds and retry")
		time.Sleep(10 * time.Second)
	}
}

func New() Service {
	return &service{
		logger: Logger(),
	}
}
