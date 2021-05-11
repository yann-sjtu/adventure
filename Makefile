export GO111MODULE=on

COMMIT := $(shell git rev-parse HEAD)
Version=dev
Name=adventure

# process linker flags
ifeq ($(VERSION),)
    VERSION = $(COMMIT)
endif

ldflags = -X github.com/okex/adventure/tools/version.Version=$(Version) \
	-X github.com/okex/adventure/tools/version.Name=$(Name) \
  	-X github.com/okex/adventure/tools/version.Commit=$(COMMIT)

ldflags := $(strip $(ldflags))
BUILD_FLAGS := -ldflags '$(ldflags)'

install:
	go install -v $(BUILD_FLAGS) .

adventure:
	go install -v $(BUILD_FLAGS) .

