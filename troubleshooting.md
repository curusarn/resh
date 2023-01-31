# Troubleshooting

## First help

Run RESH doctor to detect common issues:
```sh
reshctl doctor
```  

## Restarting RESH daemon

Sometimes restarting RESH daemon can help:
```sh
resh-daemon-restart
```

You can also start and stop RESH daemon with:
```sh
resh-daemon-start
resh-daemon-stop
```

:warning: You will get error messages in your shell when RESH daemon is not running.

## Recorded history

Your RESH history is saved in one of:
- `~/.local/share/resh/history.reshjson`
- `$XDG_DATA_HOME/resh/history.reshjson`

The format is JSON prefixed by version. Display it as json using:

```sh
cat ~/.local/share/resh/history.reshjson | sed 's/^v[^{]*{/{/' | jq .
```

You will need `jq` installed.

## Configuration

RESH config is read from one of:
- `~/.config/resh.toml` 
- `$XDG_CONFIG_HOME/resh.toml`

## Logs

Logs can be useful for troubleshooting issues.

Find RESH logs in one of:
- `~/.local/share//resh/log.json`
- `$XDG_DATA_HOME/resh/log.json`

### Log verbosity

Get more detailed logs by setting `LogLevel = "debug"` in [RESH config](#configuration).  
Restart RESH daemon for the config change to take effect: `resh-daemon-restart`

## Common issues

### Using RESH with bash on macOS

ℹ️ It is recommended to use zsh on macOS.

MacOS comes with really old bash (`bash 3.2`).  
Update it using: `brew install bash`

On macOS, bash shell does not load `~/.bashrc` because every shell runs as login shell.  
Fix it by running: `echo '[ -f ~/.bashrc ] && . ~/.bashrc' >> ~/.bash_profile`

## Github issues

Problem persists? [Create an issue ⇗](https://github.com/curusarn/resh/issues)
