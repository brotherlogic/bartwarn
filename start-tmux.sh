#!/bin/bash

# Ensure the 'prod' session exists
if ! tmux has-session -t bartwarn 2>/dev/null; then
  # Create a new session named 'prod', detached
  cd /workspaces/bartwarn
  tmux new-session -d -s bartwarn
fi
