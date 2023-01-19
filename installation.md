# Installation

## One command installation

Feel free to check the `rawinstall.sh` script before running it.

```sh
curl -fsSL https://raw.githubusercontent.com/curusarn/resh/master/scripts/rawinstall.sh | sh
```

You will need to have `curl` and `tar` installed.

## Update

Once installed RESH can be updated using:
```sh
reshctl update
```

## Clone & Install

```sh
git clone https://github.com/curusarn/resh.git
cd resh
scripts/rawinstall.sh
```

## Build from source

:warning: Building from source is intended for development and troubleshooting.

```sh
git clone https://github.com/curusarn/resh.git
cd resh
make install
```

## Uninstallation

You can uninstall RESH by running: `rm -rf ~/.resh/`.
Restart your terminal after uninstall.

### Installed files

Binaries and shell files are in: `~/.resh/`

Recorded history, device files, and logs are in: `~/.local/share/resh/` (or `${XDG_DATA_HOME}/resh/`)
RESH config file is in: `~/.config/resh.toml`

Also check your `~/.zshrc` and `~/.bashrc`.
RESH adds a necessary line there to load itself on terminal startup.
