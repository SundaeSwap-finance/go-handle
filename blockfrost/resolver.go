// MIT License
//
// Copyright (c) 2022 SundaeSwap Labs, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package blockfrost

import (
	"context"
	"fmt"
	"strings"

	gohandle "github.com/SundaeSwap-finance/go-handle"
	bfrost "github.com/blockfrost/blockfrost-go"
)

type BlockfrostResolver struct {
	client bfrost.APIClient
}

func New(key string) BlockfrostResolver {
	client := bfrost.NewAPIClient(bfrost.APIClientOptions{
		ProjectID: key,
	})
	return BlockfrostResolver{
		client: client,
	}
}

func (b BlockfrostResolver) FindAsset(ctx context.Context, policyId string, assetNameHex string) (gohandle.AssetAddress, error) {
	addresses, err := b.client.AssetAddresses(ctx, fmt.Sprintf("%v%v", policyId, assetNameHex), bfrost.APIQueryParams{})
	if err != nil {
		return gohandle.AssetAddress{}, fmt.Errorf("unable to fetch asset addresses from blockfrost: %w", err)
	}
	if len(addresses) != 1 {
		return gohandle.AssetAddress{}, fmt.Errorf("multiple assets with policyId %v and assetName %v", policyId, assetNameHex)
	}
	return gohandle.AssetAddress{
		Address: addresses[0].Address,
	}, nil
}

func (b BlockfrostResolver) LookupAddress(ctx context.Context, address string) ([]gohandle.AssetQuantity, error) {
	switch {
	case strings.HasPrefix(address, "stake"):
		assets := b.client.AccountAssociatedAssetsAll(ctx, address)
		var assetQty []gohandle.AssetQuantity
		for resp := range assets {
			if resp.Err != nil {
				return nil, fmt.Errorf("unable to fetch associated assets: %w", resp.Err)
			}
			for _, each := range resp.Res {
				assetQty = append(assetQty, gohandle.AssetQuantity{
					Asset:    each.Unit,
					Quantity: each.Quantity,
				})
			}
		}
		return assetQty, nil
	case strings.HasPrefix(address, "addr"):
		utxos := b.client.AddressUTXOsAll(ctx, address)
		var assetQty []gohandle.AssetQuantity
		for resp := range utxos {
			if resp.Err != nil {
				return nil, fmt.Errorf("unable to fetch utxos for address: %w", resp.Err)
			}
			for _, row := range resp.Res {
				for _, each := range row.Amount {
					assetQty = append(assetQty, gohandle.AssetQuantity{
						Asset:    each.Unit,
						Quantity: each.Quantity,
					})
				}
			}
		}
		return assetQty, nil
	}
	return nil, fmt.Errorf("unrecognized address %v", address)
}
