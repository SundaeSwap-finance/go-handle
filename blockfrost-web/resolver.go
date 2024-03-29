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

package blockfrostweb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	gohandle "github.com/SundaeSwap-finance/go-handle"
)

// Temporary resolver that uses http, waiting for https://github.com/blockfrost/blockfrost-go/pull/60 to get merged
// Doesn't bother with pagination for now

type BlockfrostWebResolver struct {
	key string
}

func New(key string) BlockfrostWebResolver {
	return BlockfrostWebResolver{
		key: key,
	}
}

func (b BlockfrostWebResolver) BaseUrl() (string, error) {
	if strings.HasPrefix(b.key, "mainnet") {
		return "https://cardano-mainnet.blockfrost.io/api/v0", nil
	} else if strings.HasPrefix(b.key, "testnet") {
		return "https://cardano-testnet.blockfrost.io/api/v0", nil
	} else if strings.HasPrefix(b.key, "preview") {
		return "https://cardano-preview.blockfrost.io/api/v0", nil
	} else if strings.HasPrefix(b.key, "preprod") {
		return "https://cardano-preprod.blockfrost.io/api/v0", nil
	} else {
		return "", fmt.Errorf("unrecognized environment")
	}
}

func (b BlockfrostWebResolver) FindAsset(ctx context.Context, policyId string, assetNameHex string) (gohandle.AssetAddress, error) {
	base, err := b.BaseUrl()
	if err != nil {
		return gohandle.AssetAddress{}, err
	}
	url := fmt.Sprintf("%v/assets/%v%v/addresses", base, policyId, assetNameHex)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return gohandle.AssetAddress{}, fmt.Errorf("failed to query blockfrost: %w", err)
	}
	req.Header.Add("project_id", b.key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gohandle.AssetAddress{}, fmt.Errorf("failed to query blockfrost: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return gohandle.AssetAddress{}, fmt.Errorf("failed to read blockfrost response: %w", err)
	}

	type AssetAddr struct {
		Address string `json:"address"`
	}
	var assetAddr []AssetAddr
	if err := json.Unmarshal(respBytes, &assetAddr); err != nil {
		return gohandle.AssetAddress{}, fmt.Errorf("failed to unmarshal blockfrost response: %w", err)
	}

	if len(assetAddr) != 1 {
		return gohandle.AssetAddress{}, fmt.Errorf("multiple assets with policyId %v and assetName %v", policyId, assetNameHex)
	}
	return gohandle.AssetAddress{
		Address: assetAddr[0].Address,
	}, nil
}

func (b BlockfrostWebResolver) LookupAddress(ctx context.Context, address string) ([]gohandle.AssetQuantity, error) {
	base, err := b.BaseUrl()
	if err != nil {
		return nil, err
	}
	results := []gohandle.AssetQuantity{}
	page := 1
outer:
	for {
		switch {
		case strings.HasPrefix(address, "stake"):
			url := fmt.Sprintf("%v/accounts/%v/addresses/assets?page=%v", base, address, page)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to query blockfrost: %w", err)
			}
			req.Header.Add("project_id", b.key)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("failed to query blockfrost: %w", err)
			}
			defer resp.Body.Close()

			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read blockfrost response: %w", err)
			}

			type AssetAmt struct {
				Unit     string `json:"unit"`
				Quantity string `json:"quantity"`
			}
			var assetAmts []AssetAmt
			if err := json.Unmarshal(respBytes, &assetAmts); err != nil {
				return nil, fmt.Errorf("failed to unmarshal blockfrost response: %w", err)
			}

			if len(assetAmts) == 0 {
				break outer
			}

			page += 1

			for _, each := range assetAmts {
				results = append(results, gohandle.AssetQuantity{
					Asset:    each.Unit,
					Quantity: each.Quantity,
				})
			}
		case strings.HasPrefix(address, "addr"):
			url := fmt.Sprintf("%v/addresses/%v/utxos?page=%v", base, address, page)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to query blockfrost: %w", err)
			}
			req.Header.Add("project_id", b.key)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("failed to query blockfrost: %w", err)
			}
			defer resp.Body.Close()

			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read blockfrost response: %w", err)
			}

			type AddressUTXOs struct {
				Amount []struct {
					Unit     string `json:"unit"`
					Quantity string `json:"quantity"`
				} `json:"amount"`
			}
			var utxos []AddressUTXOs
			if err := json.Unmarshal(respBytes, &utxos); err != nil {
				return nil, fmt.Errorf("failed to unmarshal blockfrost response: %w", err)
			}

			if len(utxos) == 0 {
				break outer
			}

			page += 1

			for _, utxo := range utxos {
				for _, each := range utxo.Amount {
					results = append(results, gohandle.AssetQuantity{
						Asset:    each.Unit,
						Quantity: each.Quantity,
					})
				}
			}
		default:
			return nil, fmt.Errorf("unrecognized address format %v", address)
		}
	}
	return results, nil
}
