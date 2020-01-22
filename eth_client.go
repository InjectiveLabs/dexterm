package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/InjectiveLabs/dexterm/ethfw"
	"github.com/InjectiveLabs/dexterm/ethfw/gasmeter"
	"github.com/InjectiveLabs/dexterm/ethfw/keystore"
	"github.com/InjectiveLabs/dexterm/ethfw/manager"
	"github.com/InjectiveLabs/dexterm/wrappers"
)

type EthClient struct {
	keystore          keystore.EthKeyStore
	ethManager        manager.EthManager
	gasStation        gasmeter.GasStation
	networkId         uint64
	contractAddresses map[EthContract]common.Address
	nonceCache        ethfw.NonceCache

	ercWrappers    map[common.Address]*wrappers.ERC20
	ercWrappersMux *sync.RWMutex
	weth9          *wrappers.WETH9
}

type EthContract string

const (
	EthContractERC20Proxy EthContract = "erc20proxy"
	EthContractWETH9      EthContract = "weth9"
)

func NewEthClient(
	ks keystore.EthKeyStore,
	ethManager manager.EthManager,
	defaultFromAddress common.Address,
	contractAddresses map[EthContract]common.Address,
	allowGasOracles bool,
) (*EthClient, error) {
	cli := &EthClient{
		keystore:          ks,
		ethManager:        ethManager,
		networkId:         ethManager.ChainID(),
		nonceCache:        ethfw.NewNonceCache(),
		contractAddresses: contractAddresses,
		ercWrappers:       make(map[common.Address]*wrappers.ERC20),
		ercWrappersMux:    new(sync.RWMutex),
	}

	if err := cli.initContractWrappers(); err != nil {
		return nil, err
	}

	if allowGasOracles && cli.networkId == 1 {
		// we're on Ethereum MainNet
		gasStation, err := gasmeter.NewGasStation("https://ethgasstation.info/json/ethgasAPI.json", time.Minute)
		if err != nil {
			err = errors.New("failed to connect to MainNet gas station")
			return nil, err
		}

		cli.gasStation = gasStation
	}

	cli.nonceCache.Sync(defaultFromAddress, func() (uint64, error) {
		nonce, err := cli.ethManager.PendingNonceAt(context.TODO(), defaultFromAddress)
		return nonce, err
	})

	return cli, nil
}

func (cli *EthClient) EthBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	wei, err := cli.ethManager.Balance(ctx, address)
	if err != nil {
		return nil, err
	}

	return wei.ToInt(), nil
}

func (cli *EthClient) initContractWrappers() error {
	weth9, err := wrappers.NewWETH9(cli.contractAddresses[EthContractWETH9], cli.ethManager)
	if err != nil {
		err = errors.Wrap(err, "failed to init WETH9 contract wrapper")
		return err
	}

	cli.weth9 = weth9

	return nil
}

type CallArgs struct {
	Gas      *ethfw.Wei
	From     common.Address
	FromPass string
	Context  context.Context
}

// UnlimitedAllowance is uint constant MAX_UINT = 2**256 - 1
var UnlimitedAllowance = big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))

func (cli *EthClient) erc20Wrapper(asset common.Address) (*wrappers.ERC20, error) {
	cli.ercWrappersMux.RLock()
	wrapper, ok := cli.ercWrappers[asset]
	cli.ercWrappersMux.RUnlock()

	if ok {
		return wrapper, nil
	}

	cli.ercWrappersMux.Lock()
	defer cli.ercWrappersMux.Unlock()

	wrapper, err := wrappers.NewERC20(asset, cli.ethManager)
	if err != nil {
		err = errors.Wrap(err, "failed to init ERC20 contract wrapper")
		return nil, err
	}
	cli.ercWrappers[asset] = wrapper

	return wrapper, nil
}

func (cli *EthClient) Contracts() map[EthContract]common.Address {
	return cli.contractAddresses
}

func (cli *EthClient) BalanceOf(ctx context.Context, owner, asset common.Address) (amount *big.Int, err error) {
	assetContract, err := cli.erc20Wrapper(asset)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{
		From:    owner,
		Context: ctx,
	}

	return assetContract.BalanceOf(opts, owner)
}

func (cli *EthClient) Allowance(ctx context.Context, from, spender, asset common.Address) (amount *big.Int, err error) {
	assetContract, err := cli.erc20Wrapper(asset)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{
		From:    from,
		Context: ctx,
	}

	return assetContract.Allowance(opts, from, spender)
}

func (cli *EthClient) AllowancesMap(
	ctx context.Context,
	from, spender common.Address,
	assets []common.Address,
) (map[common.Address]*big.Int, error) {

	results := make(map[common.Address]*big.Int, len(assets))
	resultsMux := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	for _, asset := range assets {
		wg.Add(1)

		go func(asset common.Address) {
			defer wg.Done()

			val, err := cli.Allowance(ctx, from, spender, asset)
			if err != nil {
				logrus.WithError(err).Warningf("unable to get allowance of %s", asset.Hex())
				return
			}

			resultsMux.Lock()
			results[asset] = val
			resultsMux.Unlock()
		}(asset)
	}

	return results, nil
}

func (cli *EthClient) BalancesMap(
	ctx context.Context,
	owner common.Address,
	assets []common.Address,
) (map[common.Address]*big.Int, error) {

	results := make(map[common.Address]*big.Int, len(assets))
	resultsMux := new(sync.Mutex)

	wg := new(sync.WaitGroup)
	wg.Add(len(assets))
	defer wg.Wait()

	for _, asset := range assets {
		go func(asset common.Address) {
			defer wg.Done()

			val, err := cli.BalanceOf(ctx, owner, asset)
			if err != nil {
				logrus.WithError(err).Warningf("unable to get balance of %s", asset.Hex())
				return
			}

			resultsMux.Lock()
			results[asset] = val
			resultsMux.Unlock()
		}(asset)
	}

	return results, nil
}

func (cli *EthClient) Approve(call *CallArgs, asset, spender common.Address, value *big.Int) (txHash common.Hash, err error) {
	opts := cli.transactOpts(call)

	var erc20Wrapper *wrappers.ERC20
	if erc20Wrapper, err = cli.erc20Wrapper(asset); err != nil {
		return
	}

	err = cli.nonceCache.Serialize(opts.From, func() error {
		nonce := cli.nonceCache.Incr(opts.From)
		var resyncUsed bool

		for {
			opts.Nonce = big.NewInt(nonce)
			opts.Context, _ = context.WithTimeout(context.Background(), 30*time.Second)

			tx, err := erc20Wrapper.Approve(opts, spender, value)
			if err != nil {
				resyncUsed, err = cli.handleTxError(err, opts.From, resyncUsed)
				if err != nil {
					// unhandled error
					return err
				}

				// try again with new nonce
				nonce = cli.nonceCache.Incr(opts.From)
				continue
			}

			txHash = tx.Hash()
			return nil
		}
	})

	return txHash, err
}

func (cli *EthClient) transactOpts(call *CallArgs) *bind.TransactOpts {
	signerFn := cli.keystore.SignerFn(call.From, call.FromPass)
	opts := &bind.TransactOpts{
		From:     call.From,
		Signer:   signerFn,
		GasLimit: cli.ethManager.GasLimit(),
		Context:  call.Context,
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	if cli.gasStation != nil {
		wei, _ := cli.gasStation.Estimate(gasmeter.GasPriorityFast)
		if wei.Gwei() == 0 {
			wei = ethfw.Gwei(10)
		}

		opts.GasPrice = wei.ToInt()
	} else if call.Gas != nil {
		opts.GasPrice = call.Gas.ToInt()
	} else {
		wei, err := cli.ethManager.SuggestGasPrice(context.TODO())
		if err != nil || wei.Int64() == 0 {
			wei = ethfw.Gwei(10).ToInt()
		}

		opts.GasPrice = wei
	}

	return opts
}

func (cli *EthClient) resyncNonces(from common.Address) {
	cli.nonceCache.Sync(from, func() (uint64, error) {
		return cli.ethManager.PendingNonceAt(context.TODO(), from)
	})
}

// handleTxError returns (bool) that indicates if nonce has been re-synced, (error) that is not nil
// incdicates a fatal error.
func (cli *EthClient) handleTxError(err error, from common.Address, resyncUsed bool) (bool, error) {
	switch {
	case err.Error() == "invalid sender":
		cli.nonceCache.Decr(from)
		err = errors.Wrap(err, "failed to sign transaction")
		return false, err

	case err.Error() == "nonce too low",
		err.Error() == "nonce too high",
		strings.HasPrefix(err.Error(), "the tx doesn't have the correct nonce"):
		if resyncUsed {
			err = errors.Wrap(err, "nonce mismatch and cannot fix by resync")
			return false, err
		}

		cli.resyncNonces(from)
		return true, nil
	default:
		if strings.HasPrefix(err.Error(), "known transaction") {
			// skip one nonce step, try to send again
			return false, nil
		}

		if strings.HasPrefix(err.Error(), "VM Exception") {
			// a VM execution consumes gas and nonce is increasing
			return false, err
		}

		cli.nonceCache.Decr(from)
		return false, err
	}
}

type AccountCreateArgs struct {
	Password       string
	PasswordRepeat string
}

func (a *AccountCreateArgs) check() error {
	if a.Password != a.PasswordRepeat {
		return errors.New("password repeat don't match")
	}

	if len(a.Password) < 8 {
		return errors.New("password must be at least 8 symbols long")
	}

	return nil
}

func ethCreateAccount(keystorePath string, args *AccountCreateArgs) (accounts.Account, error) {
	if err := args.check(); err != nil {
		return accounts.Account{}, err
	}

	return ethKeystore.StoreKey(
		keystorePath,
		args.Password,
		ethKeystore.StandardScryptN,
		ethKeystore.StandardScryptP,
	)
}

type AccountImportArgs struct {
	FilePath string
}

func ethImportAccount(keystorePath string, args *AccountImportArgs) (common.Address, error) {
	keyfilePath, err := homedir.Expand(args.FilePath)
	if err != nil {
		return common.Address{}, err
	}

	keyfileJSON, err := ioutil.ReadFile(keyfilePath)
	if err != nil {
		err = errors.Wrap(err, "unable to read keyfile")

		return common.Address{}, err
	}

	var spec *WalletSpec
	if err = json.Unmarshal(keyfileJSON, &spec); err != nil {
		err = errors.Wrap(err, "unable to read keyfile")
		return common.Address{}, err
	}

	if len(spec.Address) == 0 {
		err := errors.Errorf("failed to load address from %s", keyfilePath)
		return common.Address{}, err
	}

	info, err := os.Stat(keystorePath)
	if err != nil {
		err = errors.Wrap(err, "failed to check keystore")
		return common.Address{}, err
	} else if !info.IsDir() {
		err = errors.Errorf("keystore path is not a dir: %s", keystorePath)
		return common.Address{}, err
	}

	fileName := filepath.Base(keyfilePath)
	err = ioutil.WriteFile(filepath.Join(keystorePath, fileName), keyfileJSON, 0600)
	if err != nil {
		err = errors.Wrap(err, "failed to copy keyfile into keystore")
		return common.Address{}, err
	}

	return spec.HexToAddress(), nil
}

type AccountImportPrivKeyArgs struct {
	PrivateKeyHex  string
	Password       string
	PasswordRepeat string
}

func ethImportPrivKey(keystorePath string, args *AccountImportPrivKeyArgs) (common.Address, error) {
	privKeyHex := strings.TrimPrefix(args.PrivateKeyHex, "0x")
	pk, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		err = errors.Wrap(err, "private key failed to parse")
		return common.Address{}, err
	}

	if err := (&AccountCreateArgs{
		Password:       args.Password,
		PasswordRepeat: args.PasswordRepeat,
	}).check(); err != nil {
		return common.Address{}, err
	}

	ks := ethKeystore.NewKeyStore(
		keystorePath,
		ethKeystore.StandardScryptN,
		ethKeystore.StandardScryptP,
	)

	acc, err := ks.ImportECDSA(pk, args.Password)
	if err != nil {
		err = errors.Wrap(err, "failed to import ECDSA private key")
		return common.Address{}, err
	}

	return acc.Address, nil
}

type WalletSpec struct {
	Address string `json:"address"`
	ID      string `json:"id"`
	Version int    `json:"version"`
	Path    string `json:"-"`
}

func (spec *WalletSpec) HexToAddress() common.Address {
	return common.HexToAddress(spec.Address)
}

type AccountUseArgs struct {
	Address string
}

func ethParseAccount(args *AccountUseArgs) (common.Address, error) {
	addr := common.HexToAddress(args.Address)
	if bytes.Equal(addr.Bytes(), common.Address{}.Bytes()) {
		err := errors.Errorf("failed to parse address: %s", args.Address)
		return common.Address{}, err
	}

	return addr, nil
}
