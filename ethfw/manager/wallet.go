package manager

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/InjectiveLabs/dexterm/ethfw"
)

func (m *ethManager) Balance(ctx context.Context, account common.Address) (*ethfw.Wei, error) {
	return m.BalanceAt(ctx, account, 0)
}

func (m *ethManager) BalanceAt(ctx context.Context, account common.Address, blockNum uint64) (*ethfw.Wei, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	var blockBigNum *big.Int
	if blockNum > 0 {
		blockBigNum = big.NewInt(0)
		blockBigNum.SetUint64(blockNum)
	}
	cli := ethclient.NewClient(rpc)
	bigint, err := cli.BalanceAt(ctx, account, blockBigNum)
	if err != nil {
		return nil, err
	}
	wei := ethfw.BigWei(bigint)
	return wei, nil
}

func (m *ethManager) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return 0, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.PendingNonceAt(ctx, account)
}

func (m *ethManager) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.PendingCodeAt(ctx, account)
}

func (m *ethManager) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return 0, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.EstimateGas(ctx, msg)
}

func (m *ethManager) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.SuggestGasPrice(ctx)
}
