#!/bin/bash

if [ -n "$1" ]
then
  echo "Передан параметр адреса сервера"
else
  $1 = "http://localhost:8181"
fi

