package main

import (
	"context"
	"log"
	"os"

	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/queue"
	"github.com/leonardobiffi/gitlab-ci-pipeline-queue/version"
	"github.com/urfave/cli/v3"
)

func main() {
	var ref string

	cmd := &cli.Command{
		Name:        "gitqueue",
		Description: "GitLab Pipeline Queue",
		Version:     version.String(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ref",
				Destination: &ref,
				Usage:       "Branch or tag name to filter pipelines",
			},
		},
		Action: func(context.Context, *cli.Command) error {
			queue := queue.New()
			queue.Run(ref)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
