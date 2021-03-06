#!/usr/bin/env bash
set -e

PROJ_NAME=sea-server-go
OUTPUT_DIR=build

versionDir=github.com/hzxiao/goutil/version
PWD=`pwd`
if [[ ${PWD} = ${GOPATH}/* ]]; then
    if [[ -d vendor ]]; then
        gp=${GOPATH//\//\\\/}
        proj_path=`echo ${PWD} | sed "s/${gp}\/src\///g"`
        versionDir=${proj_path}/vendor/github.com/hzxiao/goutil/version
    fi
fi
echo "version dir: ${versionDir}"

# build with version info
gitTag=$(if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate=$(TZ=Asia/Shanghai date +%FT%T%z)
gitCommit=$(git log --pretty=format:'%H' -n 1)
gitTreeState=$(if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
#
#ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR/$PROJ_NAME
mkdir -p $OUTPUT_DIR/$PROJ_NAME/script

go build -v -ldflags "-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}" -o ${PROJ_NAME}
cp -R script/* $OUTPUT_DIR/$PROJ_NAME/script
mv $PROJ_NAME $OUTPUT_DIR/$PROJ_NAME
cp Makefile $OUTPUT_DIR/$PROJ_NAME
cd $OUTPUT_DIR
zip -r -q $PROJ_NAME.zip $PROJ_NAME

if [ -n "${SAVE_PKG_DIR}" ]; then
    cp $PROJ_NAME.zip ${SAVE_PKG_DIR}
fi
