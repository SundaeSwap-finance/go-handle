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

package mock

import (
	"context"
	"encoding/hex"
	"fmt"

	gohandle "github.com/SundaeSwap-finance/go-handle"
)

/// Mock resolver for tests
type MockResolver struct {
	Environment gohandle.Environment
	Handles     map[string]string
}

func New(env gohandle.Environment, pairs [][]string) MockResolver {
	mock := MockResolver{
		Environment: env,
		Handles:     map[string]string{},
	}
	for _, p := range pairs {
		mock.Handles[p[0]] = p[1]
	}
	return mock
}

func (m MockResolver) SetAddress(handle string, address string) {
	m.Handles[handle] = address
}

func (m MockResolver) FindAsset(ctx context.Context, policyId string, assetHex string) (gohandle.AssetAddress, error) {
	switch m.Environment {
	case gohandle.Mainnet:
		if policyId != gohandle.MainnetPolicyId {
			return gohandle.AssetAddress{}, fmt.Errorf("wrong policyId (%v) for environment (%v)", policyId, m.Environment)
		}
	case gohandle.Testnet:
		if policyId != gohandle.TestnetPolicyId {
			return gohandle.AssetAddress{}, fmt.Errorf("wrong policyId (%v) for environment (%v)", policyId, m.Environment)
		}
	default:
		return gohandle.AssetAddress{}, fmt.Errorf("unrecognized environment %v", m.Environment)
	}

	handle, err := hex.DecodeString(assetHex)
	if err != nil {
		return gohandle.AssetAddress{}, fmt.Errorf("invalid asset name: %w", err)
	}
	if addr, ok := m.Handles[string(handle)]; ok {
		return gohandle.AssetAddress{Address: addr}, nil
	}
	return gohandle.AssetAddress{}, fmt.Errorf("handle not found: %w", err)
}

func (m MockResolver) LookupAddress(ctx context.Context, address string) ([]gohandle.AssetQuantity, error) {
	var ret []gohandle.AssetQuantity
	policyId := ""
	if m.Environment == gohandle.Mainnet {
		policyId = gohandle.MainnetPolicyId
	} else if m.Environment == gohandle.Testnet {
		policyId = gohandle.TestnetPolicyId
	}
	for handle, addr := range m.Handles {
		if addr == address {
			ret = append(ret, gohandle.AssetQuantity{
				Asset:    fmt.Sprintf("%v%v", policyId, hex.EncodeToString([]byte(handle))),
				Quantity: "1",
			})
		}
	}
	return ret, nil
}
