#!/bin/bash
set -e

if [[ -z $1 ]]; then
  echo "Usage: $0 deps_dir"
  exit 1
fi

DEPS_DIR=$1

LUA_LIB="LuaJIT-2.1.0-beta3"
LUA_LIB_TGZ="${LUA_LIB}.tar.gz"
LUA_LIB_URL="http://luajit.org/download/$LUA_LIB_TGZ"

if [[ -d $DEPS_DIR/usr/local/lib ]]; then
  echo "Already built"
  exit 0
fi
mkdir -p $DEPS_DIR/tmp
wget -O "$DEPS_DIR/tmp/$LUA_LIB_TGZ" $LUA_LIB_URL
pushd $DEPS_DIR/tmp
tar -zxvf $LUA_LIB_TGZ
popd
pushd $DEPS_DIR/tmp/$LUA_LIB
DESTDIR=$DEPS_DIR make install
popd
