NAME     := $(shell basename $(CURDIR))
VERSION  := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS := $(shell find $(CURDIR) -type f -name '*.go')

GOOS   := linux
GOARCH := amd64

LDFLAGS_NAME     := -X "main.name=$(NAME)"
LDFLAGS_VERSION  := -X "main.version=v$(VERSION)"
LDFLAGS_REVISION := -X "main.revision=$(REVISION)"
LDFLAGS          := -ldflags '-s -w $(LDFLAGS_NAME) $(LDFLAGS_VERSION) $(LDFLAGS_REVISION) -extldflags -static'

RUN_KERNEL   ?= $(CURDIR)/kernel.img
RUN_INITRD   ?= $(CURDIR)/initrd.img
RUN_ROOTFS   ?= $(CURDIR)/rootfs.squashfs
RUN_METADATA ?= $(CURDIR)/meta-data
RUN_USERDATA ?= $(CURDIR)/user-data

RUN_ARGS ?=
ifneq (,$(wildcard $(RUN_KERNEL)))
RUN_ARGS += -kernel "$(RUN_KERNEL)"
endif
ifneq (,$(wildcard $(RUN_INITRD)))
RUN_ARGS += -initrd "$(RUN_INITRD)"
endif
ifneq (,$(wildcard $(RUN_ROOTFS)))
RUN_ARGS += -rootfs "$(RUN_ROOTFS)"
RUN_ARGS += -cmdline 'console=ttyS0 ds=nocloud root=file:///boot/rootfs.squashfs overlayroot=tmpfs quiet ---'
else
RUN_ARGS += -cmdline 'console=ttyS0 quiet ---'
endif
ifneq (,$(wildcard $(RUN_METADATA)))
RUN_ARGS += -metaData "$(RUN_METADATA)"
endif
ifneq (,$(wildcard $(RUN_USERDATA)))
RUN_ARGS += -userData "$(RUN_USERDATA)"
endif

.PHONY: all
all: run

.PHONY: run
run: $(CURDIR)/bin/$(NAME)
	$(CURDIR)/bin/$(NAME) $(RUN_ARGS)

.PHONY: $(NAME)
$(NAME): $(CURDIR)/bin/$(NAME)
$(CURDIR)/bin/$(NAME): $(SRCS)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $@

$(CURDIR)/bin/$(NAME).zip: $(CURDIR)/bin/$(NAME)
	cd $(CURDIR)/bin && zip $@ $(NAME)

.PHONY: test
test:
	go test -v

.PHONY: clean
clean:
	rm -rf $(CURDIR)/bin
	rm -f /tmp/disk.img
