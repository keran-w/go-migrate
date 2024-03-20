#!/bin/sh
START="${START:-1}"
END="${END:-100}"

for i in $(seq "$START" "$END"); do
  CURR="$i"
  echo "$CURR"
  sleep 1
done

