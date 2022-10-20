#! /bin/sh

sh test01.sh ./main
r01=$?

sh test02.sh ./main
r02=$?

sh test03.sh ./main
r03=$?


echo "test01.sh: $r01"
echo "test02.sh: $r02"
echo "test03.sh: $r03"
