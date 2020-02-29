#! /bin/sh


dir=$(cd $(dirname $0); pwd)
filedir=$dir/files

cd $dir
rm -f test.*.log

timeout -sKILL 3 gaze -v files/*.* | tee test.log &

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

wait

num=`cat test.log | grep "hello, world!" | wc -l`

if [ $num -ne 9 ]; then
  echo "Failed:${num}"
  exit 1
fi

echo "OK"
exit 0
