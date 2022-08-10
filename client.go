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
)

type AssetAddress struct {
	Address string
}

type Resolver interface {
	FindAsset(ctx context.Context, policyId string, assetNameHex string) (AssetAddress, error)
}

type Client struct {
	env      Environment
	policyId string
	resolver Resolver
}

func envToPolicyId(env Environment) string {
	switch env {
	case Mainnet:
		return MainnetPolicyId
	case Testnet:
		return TestnetPolicyId
	default:
		return ""
	}
}

func New(environment Environment, resolver Resolver) Client {
	return Client{
		env:      environment,
		policyId: envToPolicyId(environment),
		resolver: resolver,
	}
}

func (c Client) ResolveAddress(handle string) (address string, err error) {
	return c.ResolveAddressWithContext(context.Background(), handle)
}

func (c Client) ResolveAddressWithContext(ctx context.Context, handle string) (address string, err error) {
	handle = strings.TrimPrefix(handle, "$")

	if c.policyId == "" {
		return "", fmt.Errorf("unrecognized environment: %v", c.env)
	}

	handleHex := hex.EncodeToString([]byte(handle))

	addr, err := c.resolver.FindAsset(ctx, c.policyId, handleHex)
	if err != nil {
		return "", fmt.Errorf("unable to resolve handle: %w", err)
	}
	return addr.Address, nil
}
