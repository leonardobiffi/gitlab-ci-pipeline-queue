package queue

import (
	"regexp"
	"time"

	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/entities"
	"github.com/sirupsen/logrus"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Service interface {
	Run(flags entities.Flags)
}

type service struct {
	logger *logrus.Logger
}

var _ Service = (*service)(nil)

func (s *service) Run(flags entities.Flags) {
	env := s.getEnv()

	// retry counter for no pipeline in queue with no pipelines found
	var retriesWithNoPipeline int

	// loop until the current pipeline is the oldest one
	for {
		pipelines, err := s.fetchPipelines(env.Token, env.ProjectID, flags)
		for _, p := range pipelines {
			s.logger.Debugf("ID: %d, Status: %s, Ref: %s, Source: %s", p.ID, p.Status, p.Ref, p.Source)
		}
		if err != nil {
			s.logger.Fatalf("Failed to fetch pipelines: %v", err)
		}

		// check if the current pipeline should be ignored
		if flags.IgnoreWhen != "" {
			s.logger.Debugf("Check if pipeline contains: %s", flags.IgnoreWhen)
			if s.ignorePipelines(pipelines, env.PipelineID, flags.IgnoreWhen) {
				s.logger.Printf("Ignoring pipelines when ref contains: %s", flags.IgnoreWhen)
				return
			}
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

		if flags.Wait {
			s.logger.Println("The current pipeline is not the oldest one, let's wait for 10 seconds and retry")
			time.Sleep(10 * time.Second)
		} else {
			s.logger.Fatalf("The current pipeline ID: %d is not the oldest...", env.PipelineID)
		}
	}
}

// ignorePipelines checks if the current pipeline contains the ignoreWhen ref
func (s *service) ignorePipelines(pipelines []*gitlab.PipelineInfo, pipelineID int, ignoreWhen string) bool {
	for _, p := range pipelines {
		if p.ID == pipelineID {
			// regex to check if the ref contains the ignoreWhen string
			r := regexp.MustCompile(ignoreWhen)
			if r.MatchString(p.Ref) {
				return true
			} else {
				return false
			}
		}
	}

	return false
}

func New() Service {
	return &service{
		logger: Logger(),
	}
}
