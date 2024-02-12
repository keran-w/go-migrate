#!/bin/sh
START=${START:-1}
for i in $(seq $START 100); do
  export START=$i
  echo START
  sleep 1
done
