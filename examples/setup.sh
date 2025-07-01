#!/bin/sh

git clone https://github.com/lfi-project/lfi-runtime
cd lfi-runtime
meson setup --prefix=$PWD/install --libdir=lib build -Dbuildtype=debugoptimized
cd build
ninja install
