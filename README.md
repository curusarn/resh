# Rich Enhanced Shell History

Context-based replacement/enhancement for zsh and bash shell history

:warning: *Work in progress*

## Motivation

When you exectue a command in zsh or bash following gets recorded to your shell history:

- Command itself
- Date
- Duration of the command (only in zsh and only if enabled)

But shell usage is contextual - you probably use different commands based on additional context:

- Current directory
- Current git repository/origin
- Previously executed commands
- etc ...

Additionally it's annoying to not have your shell history follow you accros your devices.
Have you lost you history when reinstalling? I personally think this is unacceptable in 2019.

Why not synchronize your shell history accross your devices and add some metadata to know where it came from:

- Hostname
- Username
- Machine ID

## What this project does

| | Legend |
| --- | --- |
| :heavy_check_mark: | Implemented |
| :white_check_mark: | Implemented but there are issues |
| :x: | Not implemented |

*NOTE: Features can change in the future*

- :heavy_check_mark: Record shell history with metadata
  - :heavy_check_mark: save it as JSON to `~/.resh_history.json`

- :white_check_mark: Provide bindings for arrow keys
  - :heavy_check_mark: imitate default behaviour
  - :heavy_check_mark: save additional metadata (e.g. command was recalled using arrow keys)
  - :x: provide enhanced behaviour
  - :heavy_check_mark: for zsh
  - :white_check_mark: for bash

- :x: Provide an app to search the history (similar to [hstr](https://github.com/dvorka/hstr/))
  - :x: provide binding for Control+R
  - :x: allow searchnig by metadata
  - :x: app contians different search modes
  
- :x: Synchronize recorded history between devices

- :x: Provide an API to make resh extendable

- :white_check_mark: Show cool graphs based on shell history

- :heavy_check_mark: Provide a tool to sanitize the recorded history

- :heavy_check_mark: Be compatible with zsh and bash

- :heavy_check_mark: Be compatible with Linux and macOS

## Data sanitization and analysis

In order to be able to develop a good history tool I will need to get some insight into real life shell and shell history usage patterns.

This project is also my Master thesis so I need to be a bit scientific and base my design decisions on evidence/data.

Running `reshctl sanitize` creates a sanitized version of recorded history.  
In sanitized history, all sensitive information is replaced with its SHA1 hashes.

If you tried sanitizing your history and you think the result is not sanitized enough then please create an issue or message me.

If you would consider supporting my research/thesis by sending me a sanitized version of your history then please give me some contact info using this form: https://forms.gle/227SoyJ5c2iteKt98

## Prereqisities

Standard stuff: `bash`, `curl`, `tar`, ...

Additional prerequisities: `bash-completion` (if you use bash)

## Installation

### Simplest

Run this command.

```sh
curl -s https://raw.githubusercontent.com/curusarn/resh/master/scripts/rawinstall.sh | bash
```

### Simple

1. Run `git clone https://github.com/curusarn/resh.git && cd resh`
2. Run `scripts/rawinstall.sh`

## Examples

Resh history is saved to `~/.resh_history.json`

Each line is a JSON that represents one executed command line.

This is how I view it `tail -f ~/.resh_history.json | jq` or `jq < ~/.resh_history.json`.  

You can install `jq` using your favourite package manager or you can use other JSON parser to view the history.

![screenshot](img/screen.png)

## Known issues

### Q: I use bash on macOS and resh doesn't work

**A:** You have to add `[ -f ~/.bashrc ] && . ~/.bashrc` to your `~/.bash_profile`.  

**Long Answer:** Under macOS bash shell only loads `~/.bash_profile` because every shell runs as login shell. I will definitely work around this in the future but since this doesn't affect many people I decided to not solve this issue at the moment.

## Issues

You are welcome to create issues: https://github.com/curusarn/resh/issues

## Uninstallation

You can uninstall this project at any time by running `rm -rf ~/.resh/`

You won't lose any recorded history by removing `~/.resh` directory because history is saved in `~/.resh_history.json`.
