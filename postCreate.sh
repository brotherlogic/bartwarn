#!/bin/zsh
set -e

export GOPATH=/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

sudo apt-get update
sudo apt-get install -y protobuf-compiler xdg-utils tmux emacs
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Account for Ghostty
tic -x ghostty.terminfo

# Install antigravity
curl -fsSL -o /tmp/install_agy.sh https://antigravity.google/cli/install.sh
bash /tmp/install_agy.sh
rm /tmp/install_agy.sh

# Set git identity
git config --global user.email "brotherlogicautomation@gmail.com"
git config --global user.name "Brotherlogic Automation"

# Install gh extension for sub-issues
gh extension install yahsan2/gh-sub-issue

TMUX_BLOCK=$(cat << 'EOF'
if [ -z "$TMUX" ] && [ -n "$PS1" ]; then
  cd /workspaces/bartwarn
  /workspaces/bartwarn/start-tmux.sh && tmux attach-session -t bartwarn
fi
EOF
)

grep -q "tmux attach-session" ~/.zshrc || echo "$TMUX_BLOCK" >> ~/.zshrc
grep -q "tmux attach-session" ~/.bashrc || echo "$TMUX_BLOCK" >> ~/.bashrc

# Ensure the session is created
/workspaces/bartwarn/start-tmux.sh
