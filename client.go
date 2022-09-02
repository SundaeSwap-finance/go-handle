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

package gohandle

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	MainnetPolicyId = "f0ff48bbb7bbe9d59a40f1ce90e9e9d0ff5002ec48f232b49ca0fb9a"
	TestnetPolicyId = "8d18d786e92776c824607fd8e193ec535c79dc61ea2405ddf3b09fe3"
)

type Environment string

const (
	Mainnet Environment = "mainnet"
	Testnet Environment = "testnet"
	Preview Environment = "preview"
)

type AssetAddress struct {
	Address string
}

type AssetQuantity struct {
	Asset    string
	Quantity string
}

type Resolver interface {
	FindAsset(ctx context.Context, policyId string, assetNameHex string) (AssetAddress, error)
	LookupAddress(ctx context.Context, address string) ([]AssetQuantity, error)
}

type Client struct {
	env      Environment
	resolver Resolver
}

func envToPolicyId(env Environment) (string, error) {
	switch env {
	case Mainnet:
		return MainnetPolicyId, nil
	case Preview:
		fallthrough
	case Testnet:
		return TestnetPolicyId, nil
	default:
		return "", fmt.Errorf("invalid environment")
	}
}

func New(environment Environment, resolver Resolver) Client {
	return Client{
		env:      environment,
		resolver: resolver,
	}
}

func (c Client) ResolveAddress(handle string) (address string, err error) {
	return c.ResolveAddressWithContext(context.Background(), handle)
}

func (c Client) ResolveAddressWithContext(ctx context.Context, handle string) (address string, err error) {
	handle = strings.TrimPrefix(handle, "$")

	policyId, err := envToPolicyId(c.env)
	if err != nil {
		return "", err
	}

	handleHex := hex.EncodeToString([]byte(handle))

	addr, err := c.resolver.FindAsset(ctx, policyId, handleHex)
	if err != nil {
		return "", fmt.Errorf("unable to resolve handle: %w", err)
	}
	return addr.Address, nil
}

func (c Client) LookupHandles(address string) (handles []string, err error) {
	return c.LookupHandlesWithContext(context.Background(), address)
}

func (c Client) LookupHandlesWithContext(ctx context.Context, address string) (handles []string, err error) {
	policyId, err := envToPolicyId(c.env)
	if err != nil {
		return nil, err
	}
	assets, err := c.resolver.LookupAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("unable to lookup address: %w", err)
	}
	for _, asset := range assets {
		if strings.HasPrefix(asset.Asset, policyId) {
			handleHex := strings.TrimPrefix(asset.Asset, policyId)
			b, err := hex.DecodeString(handleHex)
			if err != nil {
				return nil, fmt.Errorf("invalid handle %v: %w", handleHex, err)
			}
			handles = append(handles, string(b))
		}
	}
	return handles, nil
}
