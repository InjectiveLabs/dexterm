package manager

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxInfo struct {
	Hash        common.Hash     `json:"hash"`
	BlockNumber *hexutil.Big    `json:"blockNumber"`
	Nonce       hexutil.Uint64  `json:"nonce"`
	From        common.Address  `json:"from"`
	To          *common.Address `json:"to"`
	Value       *hexutil.Big    `json:"value"`
}

func (m *ethManager) TransactionByHash(ctx context.Context, txHex string) (*TxInfo, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	var tx TxInfo
	if err := rpc.CallContext(ctx, &tx, "eth_getTransactionByHash", common.HexToHash(txHex)); err != nil {
		logrus.WithError(err).
			WithField("fn", "TransactionByHash").
			WithField("txHex", txHex).
			Warningln("failed to retrieve tx info")
		return nil, err
	}
	return &tx, nil
}

type TxReceipt struct {
	Status            hexutil.Uint64  `json:"status"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	Logs              []*types.Log    `json:"logs"`
	Hash              common.Hash     `json:"transactionHash"`
	ContractAddress   *common.Address `json:"contractAddress"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
}

func (m *ethManager) TransactionReceiptByHash(ctx context.Context, txHex string) (*TxReceipt, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	var receipt TxReceipt
	if err := rpc.CallContext(ctx, &receipt, "eth_getTransactionReceipt", common.HexToHash(txHex)); err != nil {
		logrus.WithError(err).
			WithField("fn", "TransactionReceiptByHash").
			WithField("txHex", txHex).
			Warningln("failed to retrieve tx receipt")
		return nil, err
	}
	return &receipt, nil
}

func (m *ethManager) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.SendTransaction(ctx, tx)
}
