package queue

import (
	"strings"

	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/entities"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// change use gitlab api to fetch pipelines
func (s *service) fetchPipelines(token string, projectID int, flags entities.Flags) ([]*gitlab.PipelineInfo, error) {
	git, err := gitlab.NewClient(token)
	if err != nil {
		s.logger.Fatalf("Failed to create client: %v", err)
	}

	pipelines, _, err := git.Pipelines.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
		// Scope: gitlab.Ptr("running"),
		Sort: gitlab.Ptr("desc"),
	})
	if err != nil {
		s.logger.Fatalf("Failed to list pipelines: %v", err)
	}

	// revert the order of pipelines to first by oldest
	for i, j := 0, len(pipelines)-1; i < j; i, j = i+1, j-1 {
		pipelines[i], pipelines[j] = pipelines[j], pipelines[i]
	}

	// return all pipelines if filters flags is empty
	if flags.Ref == "" && flags.RefContains == "" && flags.Source == "" && flags.RefPriority == "" {
		return pipelines, nil
	}

	var filteredPipelines []*gitlab.PipelineInfo
	for _, p := range pipelines {
		pipeline := filterPipelines(p, flags)
		if pipeline != nil {
			filteredPipelines = append(filteredPipelines, p)
		}
	}

	// order pipelines by oldest and highest priority ref if informed
	if flags.RefPriority != "" {
		var priorityPipelines []*gitlab.PipelineInfo
		var nonPriorityPipelines []*gitlab.PipelineInfo

		for _, p := range filteredPipelines {
			if p.Ref == flags.RefPriority {
				priorityPipelines = append(priorityPipelines, p)
			} else {
				nonPriorityPipelines = append(nonPriorityPipelines, p)
			}
		}

		filteredPipelines = append(priorityPipelines, nonPriorityPipelines...)
	}

	return filteredPipelines, nil
}

// filterPipelines filters pipelines by ref, ref contains and source
// and returns if match with the all flags informed
func filterPipelines(p *gitlab.PipelineInfo, flags entities.Flags) *gitlab.PipelineInfo {
	if flags.Ref != "" && flags.Ref != p.Ref {
		return nil
	}
	if flags.RefContains != "" && !strings.Contains(p.Ref, flags.RefContains) {
		return nil
	}
	if flags.Source != "" && flags.Source != p.Source {
		return nil
	}
	return p
}
