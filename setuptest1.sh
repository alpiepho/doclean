#!/bin/bash
ls ..
rm -rf testdir1
rm -rf testdir2
rm -f keeplist.txt

mkdir testdir1
touch testdir1/file1.txt
mkdir testdir1/testdir1.1
touch testdir1/testdir1.1/file1.1.txt
mkdir testdir1/testdir1.2
mkdir testdir1/testdir1.1/1.1.1

cp -rf testdir1 testdir2
touch testdir2/testdir1.1/1.1.1/file1.1.1.txt
mkdir testdir2/testdir1.3
mkdir testdir2/testdir1.4
touch testdir2/testdir1.4/file1.4.txt

echo "testdir2/testdir1.4/file1.4.txt" >> keeplist.txt

cd testdir1
find . > ../tree1.txt
cd ..
cd testdir2
find . > ../tree2.txt
cd ..
diff -u tree1.txt tree2.txt | grep -E "^\+\." | sed 's/\+//'
