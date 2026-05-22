# Agent Instructions (GEMINI.md)

This file contains custom project-specific instructions and guardrails for the Antigravity AI agent.

## Issue Management
- **DO NOT** manually close GitHub issues using the `gh` CLI (e.g., `gh issue close XX`).
- **DO** close issues by including keywords like `Fixes #XX` or `Closes #XX` in your git commit messages or PR descriptions. GitHub automation will automatically close the issue when the branch or Pull Request is merged into the default branch.

## Pull Request Workflow
- After pushing a branch, a Pull Request will be automatically created. You must proactively monitor the PR's status checks and CI/CD results.
- You must monitor the PR for any automated comments (e.g., from the Gemini Code Review workflow) and address them by pushing fixes to the branch.
- **CRITICAL**: Only once all tests are passing AND all PR comments have been addressed (meaning you consider the PR 100% done and ready to merge), you should assign `brotherlogic` as a reviewer to the PR. Continue to address any further feedback from `brotherlogic` until merged.
