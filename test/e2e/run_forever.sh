#! /bin/sh


dir=$(cd $(dirname $0); pwd)
filedir=$dir/files

cd $dir

gaze -q files/*.* &

while true; do
  sleep 0.1
  touch $filedir/hello.rb
  sleep 0.1
  touch $filedir/hello.py
  sleep 0.1
  touch $filedir/hello.rs
  sleep 0.1
  touch $filedir/hello.rb
  sleep 0.1
  touch $filedir/hello.py
  sleep 0.1
  touch $filedir/hello.rs
  sleep 0.1
  touch $filedir/hello.rb
  sleep 0.1
  touch $filedir/hello.py
  sleep 0.1
  touch $filedir/hello.rs
  sleep 0.1
done
