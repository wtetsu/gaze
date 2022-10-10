#! /bin/sh

dir=$(cd $(dirname $0); pwd)
filedir=$dir/files

cd $dir
rm -f test.*.log

cp $filedir/hello.py "$filedir/he'llo.py"
cp $filedir/hello.py "$filedir/he&llo.py"

timeout -sKILL 3 ./main -v files/*.* | tee test.log &

sleep 0.1
touch "$filedir/he'llo.py"
sleep 0.1
touch "$filedir/he&llo.py"
sleep 0.1
touch "$filedir/he'llo.py"
sleep 0.1
touch "$filedir/he&llo.py"
sleep 0.1
touch "$filedir/he'llo.py"
sleep 0.1
touch "$filedir/he&llo.py"

wait

rm "$filedir/he'llo.py"
rm "$filedir/he&llo.py"

num=`cat test.log | grep "hello, world!" | wc -l`

if [ $num -ne 6 ]; then
  echo "Failed:${num}"
  exit 1
fi

echo "OK"
exit 0
