# Repo root — delegates to serverBackendGo (backend lives in that directory).
BACKEND_DIR := serverBackendGo

.PHONY: help dev dev-stop dev-https dev-https-stop verify-https tunnel migrate db-up db-down test repair-apk-urls setup-studhub-https

help:
	@cd $(BACKEND_DIR) && $(MAKE) help

dev dev-stop dev-https dev-https-stop verify-https tunnel migrate db-up db-down test build lint tidy swagger repair-apk-urls setup-studhub-https:
	@$(MAKE) -C $(BACKEND_DIR) $@
