#!/bin/sh

name="questionsandanswers"
bindir=`dirname $0`
root="${bindir}/.."
cd "$root"
root=`pwd`

if "${root}/bin/shutdown.sh" ; then
    pidfile="${root}/var/run/${name}.pid"
    export GODEBUG="http2server=0"
    mkdir -p "${root}/var/log"
    mkdir -p "${root}/var/run"
    "${root}/bin/${name}" -port 8080 -document-root "${root}/htdocs" -log-dir "${root}/var/log" -database "${root}/var/${name}.db" >> "${root}/var/log/${name}.log" 2>&1 &
    echo $! > "$pidfile"
    chmod 644 "$pidfile"
    disown
else
    exit $?
fi
