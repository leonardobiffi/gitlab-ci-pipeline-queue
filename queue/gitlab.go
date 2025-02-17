package queue

import (
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// change use gitlab api to fetch pipelines
func (s *service) fetchPipelines(token string, projectID int, ref string) ([]*gitlab.PipelineInfo, error) {
	git, err := gitlab.NewClient(token)
	if err != nil {
		s.logger.Fatalf("Failed to create client: %v", err)
	}

	pipelines, _, err := git.Pipelines.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
		Scope: gitlab.Ptr("running"),
		Sort:  gitlab.Ptr("desc"),
	})
	if err != nil {
		s.logger.Fatalf("Failed to list pipelines: %v", err)
	}

	// revert the order of pipelines
	for i, j := 0, len(pipelines)-1; i < j; i, j = i+1, j-1 {
		pipelines[i], pipelines[j] = pipelines[j], pipelines[i]
	}

	// return all pipelines if ref is empty
	if ref == "" {
		for _, p := range pipelines {
			s.logger.Debugf("ID: %d, Status: %s, Ref: %s", p.ID, p.Status, p.Ref)
		}

		return pipelines, nil
	}

	var filteredPipelines []*gitlab.PipelineInfo
	for _, p := range pipelines {
		// filter pipelines by ref if ref is not empty
		if ref == p.Ref {
			filteredPipelines = append(filteredPipelines, p)
			s.logger.Debugf("ID: %d, Status: %s, Ref: %s", p.ID, p.Status, p.Ref)
		}
	}

	return filteredPipelines, nil
}
