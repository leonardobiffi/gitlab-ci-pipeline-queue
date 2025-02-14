package queue

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	Token      string
	PipelineID int
	ProjectID  int
}

func (s *service) getEnv() Env {
	apiToken := getEnvVar("CI_JOB_TOKEN")
	pipelineIDStr := getEnvVar("CI_PIPELINE_ID")
	projectIDStr := getEnvVar("CI_PROJECT_ID")

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		s.logger.Fatalf("Failed to convert project ID to int: %v", err)
	}

	pipelineID, err := strconv.Atoi(pipelineIDStr)
	if err != nil {
		s.logger.Fatalf("Failed to convert pipeline ID to int: %v", err)
	}

	return Env{
		Token:      apiToken,
		PipelineID: pipelineID,
		ProjectID:  projectID,
	}
}

func getEnvVar(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalf("Env var %s not present", name)
	}
	return value
}
