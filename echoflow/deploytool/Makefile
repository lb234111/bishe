CHAINMAKER_CRYPTOGEN_PATH = $(shell dirname ${PWD})/chainmaker-cryptogen
CHAINMAKER_GO_PATH = $(shell dirname ${PWD})/chainmaker-csv-adapter
PARAM_ADAPTER_PATH = $(shell dirname ${PWD})/param_adapter
FRAMEWORK_PATH = $(shell dirname ${PWD})/framework

chainmaker: environment-check chainmaker-go-build 

environment-check:
	@./scripts/env_check.sh

chainmaker-go-build: chainmaker-cryptogen-build
	@echo "build chainmaker-go..."
	@cd $(CHAINMAKER_GO_PATH)/tools && ln -fs ../../chainmaker-cryptogen/ .
	@cd $(CHAINMAKER_GO_PATH)/scripts && echo -e "1\nERROR\nYES\n\n\n" | ./prepare.sh 4 1 11301 12301
	@cd $(CHAINMAKER_GO_PATH)/scripts && ./build_release.sh


chainmaker-cryptogen-build:
	@echo "build chainmaker-cryptogen..."
	@cd $(CHAINMAKER_CRYPTOGEN_PATH) && make


build-link:
	@echo "build Symbolic Link..."
	@cd $(FRAMEWORK_PATH)/worker/chainmaker/config_files/chain1 && ln -fs $(CHAINMAKER_GO_PATH)/build/crypto-config/ .
fileExist = $(shell if [ -d ${CHAINMAKER_GO_PATH}/build ]; then echo "exist"; else echo "noexist"; fi)
.PHONY:clean
clean:
ifeq ("$(fileExist)", "exist")
	@echo $(CHAINMAKER_GO_PATH)
	@cd $(CHAINMAKER_GO_PATH)/scripts && ./cluster_quick_stop.sh
	@echo "y" | docker container prune
	@cd $(CHAINMAKER_GO_PATH) && sudo rm -rf $(CHAINMAKER_GO_PATH)/build
endif
