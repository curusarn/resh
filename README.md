# master thesis 

## What

This project is the first phase of my master thesis.

It records shell history with rich set of metadata.

## Why

The ultimate point of my thesis is to provide a drop-in replacement for bash and zsh shell history.

The idea is to provide following:
- Context-based history
- Simple way to search whole history by command itself and/or metadata
- Synchronization across devices
- And more ...

## Prereqisities

- `git`
- `golang` (>1.11 if possible but we can deal with old ones as well )

## Installation

### Simplest
Just run `bash -c "$(wget -O - https://raw.githubusercontent.com/curusarn/resh/master/rawinstall.sh)"` from anywhere.

### Simple
1. Run `git clone https://github.com/curusarn/resh.git && cd resh`
2. Run `make autoinstall` for assisted build & instalation. OR Run `make install` if you know how to build Golang projects.

## Compatibility

Works in `bash` and `zsh`.

Tested on:
- Arch
- Ubuntu
- MacOS
