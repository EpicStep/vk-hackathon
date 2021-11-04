#!/bin/bash

addr="http://localhost:8181"

if [ -n "$1" ]
then
  echo "Server address parameter passed"
  addr=$1
fi

echo "Upload photo"
curl -F "file=@_examples/test.jpg" "${addr}/image"

echo "Get photo"
