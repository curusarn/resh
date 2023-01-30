# Installation

## One command installation

Feel free to check the `rawinstall.sh` script before running it.

```sh
curl -fsSL https://raw.githubusercontent.com/curusarn/resh/master/scripts/rawinstall.sh | sh
```

You will need to have `curl` and `tar` installed.

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

## Update

Once installed RESH can be updated using:
```sh
reshctl update
```

## Disabling RESH

If you have a persistent issue with RESH you can temporarily disable it and then enable it later.  
You won't lose your history nor configuration.

Go to `~/.zshrc` and `~/.bashrc` and comment out following lines:
```sh
[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc
[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # bashrc only
```
The second line is bash-specific so you won't find it in `~/.zshrc`

You can re-enable RESH by uncommenting the lines above or by re-installing it.

## Uninstallation

You can uninstall RESH by running: `rm -rf ~/.resh/`.  
Restart all open terminals after uninstall!

### Installed files

Binaries and shell files are in `~/.resh/`.

Recorded history, device files, and logs are in one of:
- `~/.local/share/resh/`
- `$XDG_DATA_HOME/resh/` (if set)

RESH config file is read from one of:
- `~/.config/resh.toml`
- `$XDG_CONFIG_HOME/resh.toml` (if set)

RESH also adds a following lines to `~/.zshrc` and `~/.bashrc` to load itself on terminal startup:
```sh
[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc
[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # bashrc only
```

:information_source: RESH follows [XDG directory specification ⇗](https://maex.me/2019/12/the-power-of-the-xdg-base-directory-specification/)
