# Path to the environment file
ENV_FILE := .env
# Base command for running tests
GO_TEST := go test -v
# Allure
ALLURE := allure
ALLURE_OUTPUT_PATH  ?= /
ALLURE_REPORT_PATH  ?= internal/tests/allure-report
ALLURE_RESULTS_PATH ?= internal/tests/allure

.PHONY: unit integration e2e allure

update_allure_history:
	rm -rf "$(ALLURE_RESULTS_PATH)/history"
	if [ -d "$(ALLURE_REPORT_PATH)/history" ]; then \
		cp -R "$(ALLURE_REPORT_PATH)/history" "$(ALLURE_RESULTS_PATH)/history"; \
	fi
	find $(ALLURE_RESULTS_PATH) -maxdepth 1 -type f -delete

unit: update_allure_history
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -short ./...'

integration: update_allure_history
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=integration ./...'

e2e: update_allure_history
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=e2e ./...'

allure:
	$(ALLURE) generate $(ALLURE_RESULTS_PATH) --clean -o $(ALLURE_REPORT_PATH)
	$(ALLURE) open $(ALLURE_REPORT_PATH)
