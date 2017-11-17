preexec() { 
    export LASTCMD=$1
    }

precmd() { 
/data/home/main/go/src/github.com/ebastos/shell-history/shell-history -e $? $LASTCMD
export LASTCMD=""

}