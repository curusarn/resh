
# RESH Roadmap

| | Legend |
| --- | --- |
| :heavy_check_mark: | Implemented |
| :white_check_mark: | Implemented but I'm not happy with it |
| :x: | Not implemented |

*NOTE: Features can change in the future*

- :heavy_check_mark: Record shell history with metadata
  - :heavy_check_mark: save it as JSON to `~/.resh_history.json`

- :white_check_mark: Provide an app to search the history
  - :heavy_check_mark: launch with CTRL+R (enable it using `reshctl enable ctrl_r_binding_global`)
  - :heavy_check_mark: search by keywords
  - :heavy_check_mark: relevant results show up first based on context (host, directory, git, exit status)
  - :heavy_check_mark: allow searching completely without context ("raw" mode)
  - :heavy_check_mark: import and search history from before RESH was installed
  - :white_check_mark: include a help with keybindings
  - :x: allow listing details for individual commands
  - :x: allow explicitly searching by metadata

- :heavy_check_mark: Provide a `reshctl` utility to control and interact with the project
  - :heavy_check_mark: turn on/off resh key bindings
  - :heavy_check_mark: zsh completion
  - :heavy_check_mark: bash completion

- :x: Multi-device history
  - :x: Synchronize recorded history between devices
  - :x: Allow proxying history when ssh'ing into remote servers

- :x: Provide a stable API to make resh extensible

- :heavy_check_mark: Support zsh and bash

- :heavy_check_mark: Support Linux and macOS

- :white_check_mark: Require only essential prerequisite software
  - :heavy_check_mark: Linux
  - :white_check_mark: MacOS *(requires coreutils - `brew install coreutils`)*

- :heavy_check_mark: Provide a tool to sanitize the recorded history

