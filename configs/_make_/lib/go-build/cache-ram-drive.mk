
define CacheRAMDrive
    $(eval GOCACHE := /dev/shm/$(1)/go-cache)
    $(eval TMPDIR := /dev/shm/$(1)/tmp)
    $(eval GOBILDLOG := /dev/shm/$(1)/tmp/go-build.log)

    @rm -rf "$(GOCACHE)" "$(TMPDIR)"
    @mkdir -p "$(GOCACHE)" "$(TMPDIR)"

    $(eval GO_CACHE_ENV = GOCACHE=$(GOCACHE))
    $(eval GO_TMPDIR_ENV = TMPDIR=$(TMPDIR))
    $(eval GO_BUILD_LOG = $(GOBILDLOG))
endef
