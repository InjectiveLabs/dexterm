package sol

import (
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	Name            string
	SourcePath      string
	CompilerVersion string
	Address         common.Address

	ABI []byte
	Bin string
}
