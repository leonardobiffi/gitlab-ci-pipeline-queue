package main

import (
	"context"
	"log"
	"os"

	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/entities"
	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/queue"
	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/version"
	"github.com/urfave/cli/v3"
)

func main() {
	var flags entities.Flags

	cmd := &cli.Command{
		Name:        "gitqueue",
		Description: "GitLab Pipeline Queue",
		Version:     version.String(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ref",
				Destination: &flags.Ref,
				Usage:       "Branch or tag name to filter pipelines",
			},
			&cli.StringFlag{
				Name:        "ref-contains",
				Destination: &flags.RefContains,
				Usage:       "Branch or tag name to filter pipelines by contains",
			},
			&cli.StringFlag{
				Name:        "ref-priority",
				Destination: &flags.RefPriority,
				Usage:       "Branch or tag name to filter pipelines and set Higher priority",
			},
			&cli.StringFlag{
				Name:        "source",
				Destination: &flags.Source,
				Usage:       "Source of the pipeline",
			},
			&cli.BoolFlag{
				Name:        "wait",
				Destination: &flags.Wait,
				Value:       true,
				Usage:       "Wait for be the oldest pipeline",
			},
			&cli.StringFlag{
				Name:        "ignore-when",
				Destination: &flags.IgnoreWhen,
				Usage:       "Ignore when the pipeline Ref contains the informed value",
			},
		},
		Action: func(context.Context, *cli.Command) error {
			queue := queue.New()
			queue.Run(flags)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
