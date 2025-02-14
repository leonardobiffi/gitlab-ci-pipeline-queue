# Gitlab CI Pipeline Queue

A simple tool for GitLab CI that queues pipelines to prevent concurrent deployments. This workaround addresses the issue outlined in [GitLab Issue #20481](https://gitlab.com/gitlab-org/gitlab-ce/issues/20481).

It polls GitLab's API and only proceeds when the current pipeline is ready to run. Note that this method is not optimal, as it utilizes runner instances for waiting and polling. Use at your own risk.

> Inspiration from https://github.com/flegall/gitlab-ci-pipeline-queue

## Usage

In your .gitlab-ci.yml file add the following stage and job:

```yaml
stages:
  - wait

wait:
  stage: wait
  image: alpine:3.21
  script:
    - curl -fsSL https://raw.githubusercontent.com/leonardobiffi/gitlab-ci-pipeline-queue/master/scripts/install.sh | sh
    - gitqueue --ref $CI_COMMIT_REF_NAME
```
