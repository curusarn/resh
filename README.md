
[![Latest version](https://img.shields.io/github/v/tag/curusarn/resh?sort=semver)](https://github.com/curusarn/resh/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/curusarn/resh)](https://goreportcard.com/report/github.com/curusarn/resh)
[![Go test](https://github.com/curusarn/resh/actions/workflows/go.yaml/badge.svg)](https://github.com/curusarn/resh/actions/workflows/go.yaml)
[![Shell test](https://github.com/curusarn/resh/actions/workflows/sh.yaml/badge.svg)](https://github.com/curusarn/resh/actions/workflows/sh.yaml)

# REcycle SHell

Context-based replacement for `zsh` and `bash` shell history.

**Full-text search your shell history.**  
Relevant results are displayed first based on current directory, git repo, and exit status.

<!-- Contextual shell history -->
<!-- Contextual bash history -->
<!-- Contextual zsh history -->
<!-- Context-based shell history -->
<!-- Context-based bash history -->
<!-- Context-based zsh history -->
<!-- Better shell history -->
<!-- Better bash history -->
<!-- Better zsh history -->
<!-- PWD Directory -->

## Install

Install RESH with one command:

```sh
curl -fsSL https://raw.githubusercontent.com/curusarn/resh/master/scripts/rawinstall.sh | sh
```

ℹ️ You will need to have `curl` and `tar` installed.

More options on [Installation page ⇗](./installation.md)

## Search your history

Press <kbd>CTRL</kbd> + <kbd>R</kbd> to search.

<img width="906" alt="Screenshot 2023-02-25 at 18 49 07" src="https://user-images.githubusercontent.com/10132717/221371937-d4ba64e0-ede6-4bfa-8b74-529252bf73a3.png">

### IN-app key bindings

- Type to search
- <kbd>Up</kbd> / <kbd>Down</kbd> or <kbd>Ctrl</kbd> + <kbd>P</kbd> / <kbd>Ctrl</kbd> + <kbd>N</kbd> to select results
- <kbd>Enter</kbd> to execute selected command
- <kbd>Right</kbd> to paste selected command onto the command line so you can edit it before execution
- <kbd>Ctrl</kbd> + <kbd>C</kbd> or <kbd>Ctrl</kbd> + <kbd>D</kbd> to quit
- <kbd>Ctrl</kbd> + <kbd>G</kbd> to abort and paste the current query onto the command line
- <kbd>Ctrl</kbd> + <kbd>R</kbd> to search without context

## Issues & ideas

Find help on [Troubleshooting page ⇗](./troubleshooting.md)

Problem persists? [Create an issue ⇗](https://github.com/curusarn/resh/issues)
