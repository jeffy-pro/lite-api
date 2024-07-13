# CI/CD Workflow

This directory contains GitHub Actions workflows for our CI/CD pipeline. The workflows are designed to run tests, deploy to ECS, and deploy to EC2 based on specific triggers.

## Workflow Overview

1. **Reusable Workflow: Run Tests**
    - Location: `.github/workflows/tests.yml`
    - Purpose: Executes all CI tests
    - This workflow is designed to be reusable and is called by other workflows

2. **CI Workflow**
    - Trigger: Push to `main` branch or any Pull Request against `main`
    - Steps:
        1. Call the reusable "Run Tests" workflow

3. **Release Workflow**
    - Trigger: Push of a tag matching the pattern `v*` (e.g., v1.0.0)
    - Steps:
        1. Call the reusable "Run Tests" workflow
        2. If tests pass, run the following jobs in parallel:
            - Deploy to ECS
            - Deploy to EC2

## Workflow Details

### Reusable Workflow: Run Tests
This workflow contains all the steps necessary to run our CI tests. It's designed to be reusable so that we can maintain our testing process in one place and call it from multiple workflows.

### CI Workflow
This workflow ensures that all pushes to the `main` branch and all Pull Requests are tested before merging.

### Release Workflow
This workflow is triggered when a new version tag is pushed. It first runs all tests, and if they pass, it proceeds to deploy the application to both ECS and EC2 in parallel.

## Usage

- To run tests on your branch, simply create a Pull Request.
- To deploy a new version:
    1. Ensure all your changes are merged to `main`
    2. Create and push a new tag with the format `v*` (e.g., `git tag v1.0.0 && git push origin v1.0.0`)

## Configuration

Make sure to configure the necessary secrets and environment variables in your GitHub repository settings for the ECS and EC2 deployment jobs.

## Note
The implementation on AWS deployment is not complete. 
This Readme describes the CI/CD architecture.