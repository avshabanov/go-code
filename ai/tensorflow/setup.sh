
TFBASE=$HOME/opt/tensorflow

export LIBRARY_PATH=$LIBRARY_PATH:$TFBASE/lib # For both Linux and Mac OS X
#export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$TFBASE/lib # For Linux only
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:$TFBASE/lib # For Mac OS X only

# Verify:
# source setup.sh
# go get github.com/tensorflow/tensorflow/tensorflow/go
# cd $GOPATH/src/github.com/tensorflow/tensorflow/tensorflow/go/
# go test

