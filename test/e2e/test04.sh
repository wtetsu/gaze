#! /bin/sh

gaze="${1:-gaze}"

dir=$(cd $(dirname $0); pwd)
filedir=$dir/files/
nested=$filedir/deep/path

cd $dir
rm -f test.*.log
rm -rf $filedir/deep

sleep 1.0

timeout -sKILL 6 ${gaze} -v "files/**/*.*" --debug | tee test.log &

sleep 1.0

# Test new deep path after gaze started
mkdir -p $nested
echo "*" > $nested/.gitignore

sleep 1.0
cp $filedir/hello.rb $nested/hello.rb
sleep 0.3
cp $filedir/hello.py $nested/hello.py
sleep 0.3
cp $filedir/hello.rb $nested/hello.rb
sleep 0.3
cp $filedir/hello.py $nested/hello.py
sleep 0.3

wait

num=`cat test.log | grep "hello, world!" | wc -l`

if [ $num -ne 4 ]; then
  echo "Failed:${num}"
  exit 1
fi

git checkout -- files

echo "OK"
exit 0
