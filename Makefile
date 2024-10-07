build: .make .make/yarndeps .make/solidity .make/solgen
	go install

clean :
	rm -rf $(GOPATH)/bin/orbit-toolkit 
	rm -rf .make

.make/yarndeps: nitro-contracts/package.json nitro-contracts/yarn.lock .make
	yarn --cwd nitro-contracts install
	yarn --cwd token-bridge-contracts install
	@touch $@

.make/solidity: nitro-contracts/src/*/*.sol .make/yarndeps .make
	yarn --cwd nitro-contracts build
	yarn --cwd token-bridge-contracts build
	@touch $@
	
.make/solgen: solgen/gen.go .make/solidity .make
	mkdir -p solgen/go
	mkdir -p solgen/abi
	go run solgen/gen.go
	@touch $@

.make:
	mkdir .make
