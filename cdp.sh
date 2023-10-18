#!/usr/bin/env bash

nvim_set=0
github_set=0

while getopts "ngh" opt; do
  case "$opt" in
    n)
      echo "nvim"
      nvim_set=1
      ;;
    g)
      echo "github"
      github_set=1
      ;;
    h*)
      echo "help"
      ;;

  esac
done

echo $nvim_set $github_set

# move to a project directory

# fd is way faster so it will attempt to use it if available
if [[ $(command -v fd) ]]; then
  repos=$(fd -t d -g -H  "\.git" "$HOME/Repos" | sed  's/\/\.git\///')
else
  repos=$(find "$HOME/Repos" -type d -name ".git" -exec echo {} \;)
fi

repo=$(echo "$repos" | fzf)

cd "$repo" || exit

if [[ $nvim_set -eq 1 ]]; then
  nvim .
fi

if [[ $github_set -eq 1 ]]; then
  gh repo view --web
fi
