GO_FILES = $(shell go list ./...)

.PHONY: test-unit
test-unit: 
	@printf $(COLOR) "Unit test for Pakakeh... \n"
	for s in $(GO_FILES); do if ! go test -failfast -v -race $$s; then break; fi; done
	