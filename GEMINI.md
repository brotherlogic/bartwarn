# Agent Instructions (GEMINI.md)

This file contains custom project-specific instructions and guardrails for the Antigravity AI agent.

## Issue Management
- **DO NOT** manually close GitHub issues using the `gh` CLI (e.g., `gh issue close XX`).
- **DO** close issues by including keywords like `Fixes #XX` or `Closes #XX` in your git commit messages or PR descriptions. GitHub automation will automatically close the issue when the branch or Pull Request is merged into the default branch.

## Pull Request Workflow
- After pushing a branch and opening a Pull Request, you must monitor the PR's status checks.
- Once all tests and status checks are successfully passing, you should add `brotherlogic` as a reviewer to the PR.
- Once you push a branch, you should monitor the PR for any comments, and address them as they appear until the PR is merged.
