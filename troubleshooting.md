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

Two more useful commands:
```sh
resh-daemon-start
resh-daemon-stop
```

:warning: You will get error messages in your shell when RESH daemon is not running.

## Recorded history

Your RESH history is saved in one of:
- `~/.local/share/resh/history.reshjson`
- `${XDG_DATA_HOME/resh/history.reshjson`

The format is JSON prefixed by version. Display it as json using:

```sh
cat ~/.local/share/resh/history.reshjson | sed 's/^v[^{]*{/{/' | jq .
```

You will need `jq` installed.

## Logs

Logs can be useful for troubleshooting issues.

Find RESH logs in one of:
- `~/.local/share//resh/log.json`
- `${XDG_DATA_HOME}/resh/log.json`

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

## Common issues

### Using RESH in bash on macOS

MacOS comes with really old bash (`bash 3.2`).  
Update it using: `brew install bash`

On macOS, bash shell does not load `~/.bashrc` because every shell runs as login shell.  
Run  `echo '[ -f ~/.bashrc ] && . ~/.bashrc' >> ~/.bash_profile`
