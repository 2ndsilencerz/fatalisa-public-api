#! /bin/sh
# shellcheck disable=SC2006
#build=`date +"%A %F %T"`
# shellcheck disable=SC2034
#build=`date +"%A %Y/%m/%d %H:%M:%S %Z"`
build=`date +"%Y%m%d%H%M%S"`
export BUILD_DATE=build
echo ${build} > build-date.txt