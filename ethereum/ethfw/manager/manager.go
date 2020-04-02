package manager

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/InjectiveLabs/dexterm/ethereum/ethfw"
)

var (
	errNodeUnavailable = errors.New("no geth node available")
)

var _ bind.ContractBackend = &ethManager{}

type EthManager interface {
	bind.ContractCaller
	bind.ContractFilterer

	Balance(ctx context.Context, account common.Address) (*ethfw.Wei, error)
	BalanceAt(ctx context.Context, account common.Address, blockNum uint64) (*ethfw.Wei, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	TransactionByHash(ctx context.Context, txHex string) (*TxInfo, error)
	TransactionReceiptByHash(ctx context.Context, txHex string) (*TxReceipt, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error

	ChainID() uint64
	GasLimit() uint64

	Nodes() []string
	Close() error
}

func NewManager(nodes []string, gasLimit uint64) EthManager {
	m := &ethManager{
		chainID:  1,
		session:  NewSessionID(),
		ring:     hashring.New(nodes),
		ringMux:  new(sync.RWMutex),
		fails:    make(map[string]int),
		closeC:   make(chan struct{}, 1),
		wg:       new(sync.WaitGroup),
		gasLimit: gasLimit,
	}
	if len(nodes) > 0 {
		if cli, _, ok := m.rpcClient(context.Background()); ok {
			var id hexutil.Uint64
			if err := cli.Call(&id, "eth_chainId"); err != nil {
				logrus.WithError(err).
					WithField("fn", "NewManager").
					Warningln("failed to get chainID from Geth node")
			} else if uint64(id) != m.chainID {
				m.chainID = uint64(id)
			}
		}
	}
	return m
}

type ethManager struct {
	session  string
	chainID  uint64
	ring     *hashring.HashRing
	ringMux  *sync.RWMutex
	fails    map[string]int
	closeC   chan struct{}
	wg       *sync.WaitGroup
	gasLimit uint64
}

func (m *ethManager) Nodes() []string {
	m.ringMux.RLock()
	nodes, _ := m.ring.GetNodes("", m.ring.Size())
	m.ringMux.RUnlock()
	return nodes
}

func (m *ethManager) Close() error {
	close(m.closeC)
	m.wg.Wait()
	return nil
}

func (m *ethManager) ChainID() uint64 {
	return m.chainID
}

func (m *ethManager) GasLimit() uint64 {
	return m.gasLimit
}

func (m *ethManager) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.CodeAt(ctx, contract, blockNumber)
}

func (m *ethManager) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.CallContract(ctx, call, blockNumber)
}

func (m *ethManager) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.FilterLogs(ctx, query)
}

func (m *ethManager) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, m.closeC)
	defer cancelFn()
	rpc, _, ok := m.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}
	cli := ethclient.NewClient(rpc)
	return cli.SubscribeFilterLogs(ctx, query, ch)
}

func (m *ethManager) rpcClient(ctx context.Context) (cli *rpc.Client, addr string, ok bool) {
	for {
		m.ringMux.RLock()
		addr, ok = m.ring.GetNode(m.session)
		m.ringMux.RUnlock()
		if !ok {
			logrus.WithField("fn", "rpcClient").
				Warningln("no available geth nodes in pool, all dead x_X")
			return nil, "", false
		}
		newCli, err := rpc.Dial(addr)
		if err == nil {
			cli = newCli
			break
		}
		logrus.WithError(err).
			WithField("addr", addr).
			WithField("fn", "rpcClient").
			Warningln("failed to connect to geth node")
		m.failNode(addr)
		time.Sleep(3 * time.Second)
	}

	return cli, addr, ok
}

func (m *ethManager) failNode(addr string) {
	m.ringMux.Lock()
	defer m.ringMux.Unlock()
	if m.fails[addr] < 0 {
		// node been removed
		return
	}
	m.fails[addr]++
	if m.fails[addr] < 3 {
		return
	}
	m.fails[addr] = -1
	m.ring = m.ring.RemoveNode(addr)
	logrus.WithField("addr", addr).
		WithField("fn", "failNode").
		Warningln("geth node has been removed from pool and will be checked again in 5min")
	go func() {
		// schedule a revival
		time.Sleep(5 * time.Minute)
		m.reviveNode(addr)
	}()
}

func (m *ethManager) reviveNode(addr string) {
	m.ringMux.Lock()
	defer m.ringMux.Unlock()
	if m.fails[addr] >= 0 {
		// node been restored
		return
	}
	logrus.WithField("addr", addr).
		WithField("fn", "reviveNode").
		Warningln("geth node has been added back into pool")

	m.ring = m.ring.AddNode(addr)
	m.fails[addr] = 0
}
