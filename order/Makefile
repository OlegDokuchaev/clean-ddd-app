# Environment file
ENV_FILE := .env

# Base commands
GO_TEST := go test -v
GO_VET := go vet
GO_LINT := golangci-lint
GO_IMPORTS := goimports

# Allure
ALLURE := allure
ALLURE_OUTPUT_PATH  ?= /
ALLURE_REPORT_PATH  ?= internal/tests/allure-report
ALLURE_RESULTS_PATH ?= internal/tests/allure

# Thresholds
CYCLO_LIMIT ?= 10
MAINT_LIMIT ?= 20

# Tools detection
COMPLEXITY := $(shell command -v complexity)

.PHONY: unit integration e2e allure cycloover maintunder lint imports

update_allure_history:
	rm -rf "$(ALLURE_RESULTS_PATH)/history"
	if [ -d "$(ALLURE_REPORT_PATH)/history" ]; then \
		cp -R "$(ALLURE_REPORT_PATH)/history" "$(ALLURE_RESULTS_PATH)/history"; \
	fi
	if [ -d "$(ALLURE_RESULTS_PATH)" ]; then \
    	find "$(ALLURE_RESULTS_PATH)" -maxdepth 1 -type f -delete; \
    fi

unit: update_allure_history
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -short ./...'

integration: update_allure_history
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=integration ./...'

e2e: update_allure_history
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=e2e ./...'

allure:
	$(ALLURE) generate $(ALLURE_RESULTS_PATH) --clean -o $(ALLURE_REPORT_PATH)
	$(ALLURE) open $(ALLURE_REPORT_PATH)

cycloover:
	$(GO_VET) -vettool="$(COMPLEXITY)" --cycloover $(CYCLO_LIMIT) ./...

maintunder:
	$(GO_VET) -vettool="$(COMPLEXITY)" --maintunder $(MAINT_LIMIT) ./...

lint:
	$(GO_LINT) run

imports:
	$(GO_IMPORTS) -w .
