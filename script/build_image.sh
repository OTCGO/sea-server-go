#!/usr/bin/env sh

versionDir=github.com/hzxiao/goutil/version

# build with version info
gitTag=$(if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate=$(TZ=Asia/Shanghai date +%FT%T%z)
gitCommit=$(git log --pretty=format:'%H' -n 1)
gitTreeState=$(if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
#

CGO_ENABLED=0 GOOS=linux go build -v -ldflags "-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}" -o sea-server-go
