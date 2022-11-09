.PHONY: test 
test: check-test-env
	go test ./... -v -parallel 8

check-test-env:
ifndef UPCLOUD_GO_SDK_TEST_USER
	$(error UPCLOUD_GO_SDK_TEST_USER is undefined)
endif
ifndef UPCLOUD_GO_SDK_TEST_PASSWORD
	$(error UPCLOUD_GO_SDK_TEST_PASSWORD is undefined)
endif
