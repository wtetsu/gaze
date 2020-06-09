#! /bin/sh


dir=$(cd $(dirname $0); pwd)
filedir=$dir/files

cd $dir
rm -f test.*.log

timeout -sKILL 3 gaze -v -c "ruby {{file}} 1" -r files/*.*  | tee test.log &

touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1
touch $filedir/hello.rb
sleep 0.1

wait

num=`cat test.log | grep "hello, world!" | wc -l`

if [ $num -ne 1 ]; then
  echo "Failed:${num}"
  exit 1
fi

echo "OK"
exit 0
