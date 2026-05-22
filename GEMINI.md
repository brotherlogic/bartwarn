# Agent Instructions (GEMINI.md)

This file contains custom project-specific instructions and guardrails for the Antigravity AI agent.

## Issue Management
- **DO NOT** manually close GitHub issues using the `gh` CLI (e.g., `gh issue close XX`).
- **DO** close issues by including keywords like `Fixes #XX` or `Closes #XX` in your git commit messages or PR descriptions. GitHub automation will automatically close the issue when the branch or Pull Request is merged into the default branch.
