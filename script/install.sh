#!/usr/bin/env bash
set -e

PKG_NAME=sea-server-go

if [[ -z "${SAVE_PKG_DIR}" ]]; then
    echo "empty env value of SAVE_PKG_DIR"
    exit 1
fi

if [[ -z "${INSTALL_DIR}" ]]; then
    echo "empty env value of INSTALL_DIR"
    exit 1
fi

# check conf args

if [[ -z "${PROGRAM}" ]]; then
    echo "empty env value of PROGRAM"
    exit 1
fi

if [[ -z "${CMD}" ]]; then
    echo "empty env value of CMD"
    exit 1
fi

if [[ -z "${DIR}" ]]; then
    echo "empty env value of DIR for app directory"
    exit 1
fi

if [[ -z "${STDERR}" ]]; then
    echo "empty env value of STDERR"
    exit 1
fi

if [[ -z "${STDOUT}" ]]; then
    echo "empty env value of STDOUT"
    exit 1
fi

rm -f $INSTALL_DIR/$PKG_NAME.zip
rm -fr $INSTALL_DIR/$PKG_NAME

cp $SAVE_PKG_DIR/$PKG_NAME.zip $INSTALL_DIR
cwd=`pwd`
cd $INSTALL_DIR
unzip $PKG_NAME.zip
cd $cwd

if supervisorctl status $PROGRAM > /dev/null 2>&1; then
    echo "服务 $PROGRAM 已安装，直接重启"
    supervisorctl restart $PROGRAM
else
    cat script/app.supervisor.conf /etc/supervisor/conf.d/$PROGRAM.conf
    supervisorctl start $PROGRAM
fi