export GO111MODULE=on

COMMIT := $(shell git rev-parse HEAD)
Version=v0.11.1
Name=adventure
GLOBAL_CONFIG_PATH := $(shell pwd)/config.toml
TX_CONFIG_PATH := $(shell pwd)/template/sample.json
# process linker flags
ifeq ($(VERSION),)
    VERSION = $(COMMIT)
endif

ldflags = -X github.com/okex/adventure/version.Version=$(Version) \
	-X github.com/okex/adventure/version.Name=$(Name) \
  	-X github.com/okex/adventure/version.Commit=$(COMMIT) \
  	-X github.com/okex/adventure/common/config.GlobalConfigPath=$(GLOBAL_CONFIG_PATH) \
  	-X github.com/okex/adventure/common/config.TxConfigPath=$(TX_CONFIG_PATH) \

ldflags := $(strip $(ldflags))
BUILD_FLAGS := -ldflags '$(ldflags)'

install:
	go install -v $(BUILD_FLAGS) .

adventure:
	go install -v $(BUILD_FLAGS) .

