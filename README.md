
[![Latest version](https://img.shields.io/github/v/tag/curusarn/resh?sort=semver)](https://github.com/curusarn/resh/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/curusarn/resh)](https://goreportcard.com/report/github.com/curusarn/resh)
[![Go test](https://github.com/curusarn/resh/actions/workflows/go.yaml/badge.svg)](https://github.com/curusarn/resh/actions/workflows/go.yaml)
[![Shell test](https://github.com/curusarn/resh/actions/workflows/sh.yaml/badge.svg)](https://github.com/curusarn/resh/actions/workflows/sh.yaml)

# REcycle SHell

Context-based replacement for `zsh` and `bash` shell history.

Find any command you have used before.  
Releveant results are displayed first based on current directory, git repo, and exit status.

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

You will need to have `curl` and `tar` installed.

More options on [Installation page ⇗](./installation.md)

## Update

Once installed RESH can be updated using:

```sh
reshctl update
```

## Search your history

TODO: redo this

Draft:
See RESH in action - record a terminal video

Recording content:
Search your history by commands - Show searching some longer command

Get results based on current context - Show getting project-specific commands

Find any command - Show searching where the context brings the relevant command to the top

Start searching now - Show search in native shell histories


Press CTRL+R to search.
Say bye to weak standard history search.



TODO: This doesn't seem like the right place for keybindings

### In-app key bindings

- Type to search/filter
- Up/Down or Ctrl+P/Ctrl+N to select results
- Right to paste selected command onto the command line so you can edit it before execution
- Enter to execute
- Ctrl+C/Ctrl+D to quit
- Ctrl+G to abort and paste the current query onto the command line
- Ctrl+R to switch between RAW and NORMAL mode

## Issues & Ideas

Find help on [Troubleshooting page ⇗](./troubleshooting.md)

Problem persists? [Issues ⇗](https://github.com/curusarn/resh/issues)
