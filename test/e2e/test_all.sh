#! /bin/sh

# dry-run
sh test01.sh ./main

sh test01.sh ./main
r01=$?

sh test02.sh ./main
r02=$?

sh test03.sh ./main
r03=$?

sh test04.sh ./main
r04=$?


echo "test01.sh: $r01"
echo "test02.sh: $r02"
echo "test03.sh: $r03"
echo "test04.sh: $r04"

if [ $r01 -ne 0 ] || [ $r02 -ne 0 ] || [ $r03 -ne 0 ] || [ $r04 -ne 0 ]; then
  exit 1
fi