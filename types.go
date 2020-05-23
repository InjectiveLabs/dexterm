package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"

	zeroex "github.com/InjectiveLabs/zeroex-go"
	sraAPI "github.com/InjectiveLabs/injective-core/api/gen/relayer_api"
)

func ro2zo(o *sraAPI.Order) (*zeroex.SignedOrder, error) {
	if o == nil {
		return nil, nil
	}
	order := zeroex.Order{
		ChainID:             big.NewInt(o.ChainID),
		ExchangeAddress:     common.HexToAddress(o.ExchangeAddress),
		MakerAddress:        common.HexToAddress(o.MakerAddress),
		TakerAddress:        common.HexToAddress(o.TakerAddress),
		MakerAssetData:      common.FromHex(o.MakerAssetData),
		TakerAssetData:      common.FromHex(o.TakerAssetData),
		MakerFeeAssetData:   common.FromHex(o.MakerFeeAssetData),
		TakerFeeAssetData:   common.FromHex(o.TakerFeeAssetData),
		SenderAddress:       common.HexToAddress(o.SenderAddress),
		FeeRecipientAddress: common.HexToAddress(o.FeeRecipientAddress),
	}
	if v, ok := math.ParseBig256(o.MakerAssetAmount); !ok {
		return nil, errors.New("makerAssetAmmount parse failed")
	} else {
		order.MakerAssetAmount = v
	}
	if v, ok := math.ParseBig256(o.MakerFee); !ok {
		return nil, errors.New("makerFee parse failed")
	} else {
		order.MakerFee = v
	}
	if v, ok := math.ParseBig256(o.TakerAssetAmount); !ok {
		return nil, errors.New("takerAssetAmmount parse failed")
	} else {
		order.TakerAssetAmount = v
	}
	if v, ok := math.ParseBig256(o.TakerFee); !ok {
		return nil, errors.New("takerFee parse failed")
	} else {
		order.TakerFee = v
	}
	if v, ok := math.ParseBig256(o.ExpirationTimeSeconds); !ok {
		return nil, errors.New("expirationTimeSeconds parse failed")
	} else {
		order.ExpirationTimeSeconds = v
	}
	if v, ok := math.ParseBig256(o.Salt); !ok {
		return nil, errors.New("salt parse failed")
	} else {
		order.Salt = v
	}
	signedOrder := &zeroex.SignedOrder{
		Order:     order,
		Signature: common.FromHex(o.Signature),
	}
	return signedOrder, nil
}
