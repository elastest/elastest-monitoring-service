#!/bin/bash

clearpath() {
	cd $1
	local PREFIXPATH=$2
	REPLACEPATH="..\/..\/home\/felipe.gorostiaga\/pruebaswagger\/"
	# pwd
	for f in *.go; do

		## Check if the glob gets expanded to existing files.
		## If not, f here will be exactly the pattern above
		## and the exists test will evaluate to false.
		[ -e "$f" ] && sed -i "s/$REPLACEPATH/$PREFIXPATH/g" *.go

		## This is all we needed to know, so we can break after the first iteration
		break
	done
	for d in */; do
		[ $d = '*/' ] || clearpath $d "..\\/$PREFIXPATH"
	done
	cd ..
}

clearpath . ""
