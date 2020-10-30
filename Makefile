.PHONY: waves

waves:
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "surfboard" 2> /dev/null || echo 'Please install sorfboard'
	npm update surfboard -g
	surfboard compile ./contracts/waves/gravity.ride > ./abi/waves/gravity.abi
	surfboard compile ./contracts/waves/nebula.ride > ./abi/waves/nebula.abi
	echo "Waves abi updated"