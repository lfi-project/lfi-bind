ROOT=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))/..

ARCH ?= $(shell uname -m)
LFIBIND ?= $(ROOT)/lfi-bind
LFICC ?= $(ARCH)-lfi-linux-musl-clang
CC ?= cc
LIBLFI ?= $(ROOT)/examples/lfi-runtime/install/lib/liblfi.so
INCLUDELFI ?= $(ROOT)/examples/lfi-runtime/install/include
