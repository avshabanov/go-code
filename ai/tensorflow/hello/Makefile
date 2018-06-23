
TFBASE=${HOME}/opt/tensorflow
CFLAGS=-I${TFBASE}/include
LFLAGS=-L${TFBASE}/lib -ltensorflow

.PHONY: all
all: target/hello_tf

target/hello_tf: target hello_tf.c
	gcc ${CFLAGS} ${LFLAGS} hello_tf.c -o target/hello_tf

target:
	mkdir target

.PHONY: clean
clean:
	rm -rf target

