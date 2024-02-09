#!/bin/sh
START=${START:-1}
for i in $(seq $START 100); do
  echo $i
  sleep 1
done
