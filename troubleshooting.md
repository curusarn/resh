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

## Logs



## Disabling RESH

If you have a persistent issue with RESH you can temporarily disable it.

Go to `~/.zshrc` and `~/.bashrc` and comment out following lines:
```sh
[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc
[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # bashrc only
```
The second line is bash-specific so you won't find it in `~/.zshrc`


### RESH in bash on macOS doesn't work

**A:** Add line `[ -f ~/.bashrc ] && . ~/.bashrc` to your `~/.bash_profile`.

**Long Answer:** Under macOS bash shell only loads `~/.bash_profile` because every shell runs as login shell.
