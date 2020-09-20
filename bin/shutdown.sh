#!/bin/sh

name="questionsandanswers"
bindir=`dirname $0`
root="${bindir}/.."
cd "$root"
root=`pwd`
pidfile="${root}/var/run/${name}.pid"
if [ ! -e "$pidfile" ] ; then
    echo "${name} not running"
    exit 0
fi
if [ ! -w "$pidfile" ] ; then
    echo "${name} not started by you"
    exit 1
fi
pid=`cat "$pidfile"`
kill $pid
for i in 1 1 2 3 5 8 ; do
    if kill -0 $pid >>/dev/null 2>&1 ; then
        sleep $i
    else
        rm "$pidfile"
        exit 0
    fi
done
echo "${name} hasn't exited, giving up"
exit 1
