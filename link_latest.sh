#!/bin/bash

BASENAME=$(find . -type l -ls | awk '{print $11}')

if [ -z "$BASENAME" ]; then
  read -p "What is the basename for this service: " BASENAME
  if [ -z "$BASENAME" ]; then
    echo "Exiting."
    exit 1
  fi
fi

LATEST=$(ls ${BASENAME}_* | sort -V | tail -1)

rm $BASENAME > /dev/null 2>&1
ln -s $LATEST $BASENAME

while true; do
  read -p "Do you want to restart $BASENAME? [yes/NO]:" yn
  case $yn in
    yes ) ./stop.sh && ./start.sh && ./tailLog.sh; break;;
    [Nn]* ) exit;;
    * ) echo "Please answer yes or no.";;
  esac
done
