#!/bin/bash

name=$1
set -ex
#upx --brute dist/$name*/$name
upx dist/$name*/$name
