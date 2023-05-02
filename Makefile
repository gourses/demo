.DEFAULT_GOAL := help

# Variables for text transformations.
COMMA := ,
SPACE := $(subst ,, )

# Tool related definitions.
GO_VERSION		:= 1.20.3
TOOLS_DIR		:= .tools
GO 				:= ${TOOLS_DIR}/go/go${GO_VERSION}
GOLANGCI_LINT	:= ${TOOLS_DIR}/github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
COVMERGE		:= ${TOOLS_DIR}/github.com/wadey/gocovmerge@master

# Tool installation helpers.
MAJOR_VER	= $(firstword $(subst ., ,$(lastword $(subst @, ,${@}))))
LAST_PART	= $(notdir $(firstword $(subst @, ,${@})))
BIN_PATH	= ${PWD}/${@D}
BIN_NAME	= $(if $(filter ${LAST_PART},$(MAJOR_VER)),$(notdir ${BIN_PATH}),${LAST_PART})

# Build and test variables.
APP_NAME		:= demo
BUILD_TARGET	:= ./cmd/${APP_NAME}
BUILD_OUTPUT	:= target/bin/${APP_NAME}
COVER_DIR		:= target/coverage
COV_FILE_UNIT	:= ${COVER_DIR}/unit.out
COV_FILE_INT	:= ${COVER_DIR}/integration.out
COV_FILE_TOTAL	:= ${COVER_DIR}/merged.out
COVERPKG_INT	= $(subst ${SPACE},${COMMA},$(shell ${GO} list ./...))

.PHONY: build lint unit-test integration-test merge-coverages

help: ## Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<recipe>\033[0m\n\nRecipes:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

${GO}: ## Install required Go version
	@echo Installing Go ${GO_VERSION}
	@GOBIN=${PWD}/$(dir ${GO}) go install -mod=readonly golang.org/dl/go${GO_VERSION}@latest
	${GO} download
	@echo Go ${GO_VERSION} installed!

${GOLANGCI_LINT} ${COVMERGE}: ${GO} ## Install tools
	@echo Installing ${@:${TOOLS_DIR}/%=%}
	@mkdir -p ${BIN_PATH}
	@cd $(shell mktemp -d) && GOFLAGS='' GOBIN='${BIN_PATH}' ${PWD}/${GO} install ${@:${TOOLS_DIR}/%=%}
	@mv ${BIN_PATH}/${BIN_NAME} ${@}
	@echo ${@:${TOOLS_DIR}/%=%} installed!

build: ${GO} ## Build binary
	CGO_ENABLED=0 ${GO} build -o ${BUILD_OUTPUT} ${BUILD_TARGET}
	@echo Binary ${BUILD_OUTPUT} built successfully!

image: build
	docker build -t ${APP_NAME}:latest .

lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} run

unit-test: ${GO} ## Run unit tests
	$(if $(dir ${COV_FILE_UNIT}),$(shell mkdir -p $(dir ${COV_FILE_UNIT})))
	${GO} test -race -covermode atomic -coverprofile=${COV_FILE_UNIT} ./...

integration-test: ${GO} ## Run integration tests
	$(if $(dir ${COV_FILE_INT}),$(shell mkdir -p $(dir ${COV_FILE_INT})))
	${GO} test -race -covermode atomic -coverprofile=${COV_FILE_INT} -tags=integration -coverpkg=${COVERPKG_INT} ${BUILD_TARGET}

merge-coverages: ${COVMERGE} unit-test integration-test ## Merge unit and integration coverages
	${COVMERGE} ${COV_FILE_UNIT} ${COV_FILE_INT} > ${COV_FILE_TOTAL}

show-coverage: merge-coverages ## Open coverage report in browser
	${GO} tool cover -html=${COV_FILE_TOTAL}

clean: ## Remove all build and test artifacts
	rm -r target

clean-tools: ## Remove all tools
	rm -r .tools