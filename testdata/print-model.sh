#!/bin/bash

marg=""
for arg in $*; do
	if [[ "$arg" == "--model" ]]; then
		marg="yes"
	elif [[ "$marg" == "yes" ]]; then
		echo "$arg"
		exit 0
	fi
done
exit 1
