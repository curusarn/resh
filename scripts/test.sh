#!/usr/bin/env bash
# very simple tests to catch simple errors in scripts

for f in scripts/*.sh; do
    echo "Running shellcheck on $f ..."
    shellcheck "$f" --shell=sh --severity=error || exit 1
done

for f in scripts/{shellrc,hooks}.sh; do
    echo "Checking Zsh syntax of $f ..."
    ! zsh -n "$f" && echo "Zsh syntax check failed!" && exit 1
done

if [ "$1" = "--all" ]; then
	for sh in bash zsh; do
	    echo "Running functions in scripts/shellrc.sh using $sh ..."
	    ! $sh -c ". scripts/shellrc.sh; __resh_preexec; __resh_precmd" && echo "Error while running functions!" && exit 1
	done
fi

# TODO: test installation

exit 0
