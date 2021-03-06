#!/bin/sh

SOURCEDIRECTORY="github.com/bborbe/auth"

################################################################################

echo "use workspace ${WORKSPACE}"

export GOROOT=/opt/go
export PATH=/opt/go2xunit/bin/:/opt/utils/bin/:/opt/aptly_utils/bin/:/opt/aptly/bin/:/opt/debian_utils/bin/:/opt/debian/bin/:$GOROOT/bin:$PATH
export GOPATH=${WORKSPACE}
export REPORT_DIR=${WORKSPACE}/test-reports
INSTALLS=`cd src && find $SOURCEDIRECTORY/bin -name "*.go" | dirof | unique`
rm -rf $REPORT_DIR ${WORKSPACE}/*.deb ${WORKSPACE}/pkg
mkdir -p $REPORT_DIR
PACKAGES=`cd src && find $SOURCEDIRECTORY -name "*_test.go" | dirof | unique`
FAILED=false
for PACKAGE in $PACKAGES
do
  XML=$REPORT_DIR/`pkg2xmlname $PACKAGE`
  OUT=$XML.out
  go test -i $PACKAGE
  go test -v $PACKAGE | tee $OUT
  cat $OUT
  go2xunit -fail=true -input $OUT -output $XML
  rc=$?
  if [ $rc != 0 ]
  then
    echo "Tests failed for package $PACKAGE"
    FAILED=true
  fi
done

if $FAILED
then
  echo "Tests failed => skip install"
  exit 1
else
  echo "Tests success"
fi
