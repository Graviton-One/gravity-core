.PHONY: ethereum waves

ethereum:
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "abigen" 2> /dev/null || echo 'Please install abigen'
	# Core
	# Gravity
	abigen --pkg=gravity --sol="./contracts/ethereum/Gravity/Gravity.sol" --out="./abi/ethereum/gravity/gravity.go"
	abigen --pkg=nebula --sol="./contracts/ethereum/Nebula/Nebula.sol" --out="./abi/ethereum/nebula/nebula.go"
	echo "Ethereum ABI for Nebula & Gravity contracts updated!"

waves:
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "surfboard" 2> /dev/null || echo 'Please install sorfboard'
	surfboard compile ./contracts/waves/gravity.ride > ./abi/waves/gravity.abi
	surfboard compile ./contracts/waves/nebula.ride > ./abi/waves/nebula.abi
	echo "Waves ABI for Nebula and Gravity contracts updated!"
