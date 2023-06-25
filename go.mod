module github.com/meshplus/eth-kit

go 1.13

replace (
	github.com/meshplus/bitxhub-kit => github.com/TaiChiChain/bitxhub-kit v1.20.1-0.20230625064843-b7ae66571e52
	github.com/meshplus/bitxhub-model => github.com/TaiChiChain/bitxhub-model v1.20.2-0.20230625065636-b5b1ab540d61
)

require (
	github.com/ethereum/go-ethereum v1.10.8
	github.com/holiman/uint256 v1.2.0
	github.com/meshplus/bitxhub-kit v1.2.1-0.20210902085548-07f4fa85bfc9
	github.com/meshplus/bitxhub-model v1.2.1-0.20210805020100-41c27414badf
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)
