obj-m += cpu.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean

# CFLAGS += -I/lib/modules/$(shell uname -r)/build/include