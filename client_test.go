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

package gohandle_test

import (
	"testing"

	gohandle "github.com/SundaeSwap-finance/go-handle"
	"github.com/SundaeSwap-finance/go-handle/mock"
	"github.com/tj/assert"
)

func Test_Lookup(t *testing.T) {
	m := mock.New(gohandle.Mainnet, [][]string{{"abc", "addr1xyz"}})
	c := gohandle.New(gohandle.Mainnet, m)
	addr, err := c.ResolveAddress("abc")
	assert.Nil(t, err)
	assert.Equal(t, "addr1xyz", addr)
	_, err = c.ResolveAddress("xyz")
	assert.NotNil(t, err)
}

func Test_LookupTestnet(t *testing.T) {
	m := mock.New(gohandle.Testnet, [][]string{{"abc", "addr1xyz"}})
	c := gohandle.New(gohandle.Testnet, m)
	addr, err := c.ResolveAddress("abc")
	assert.Nil(t, err)
	assert.Equal(t, "addr1xyz", addr)
	_, err = c.ResolveAddress("xyz")
	assert.NotNil(t, err)
}

func Test_CrossEnv(t *testing.T) {
	m := mock.New(gohandle.Mainnet, [][]string{{"abc", "addr1xyz"}})
	c := gohandle.New(gohandle.Testnet, m)
	_, err := c.ResolveAddress("abc")
	assert.NotNil(t, err)
}

func Test_Reverse(t *testing.T) {
	m := mock.New(gohandle.Mainnet, [][]string{{"abc", "addr1xyz"}, {"xyz", "addr1xyz"}, {"www", "addr1www"}})
	c := gohandle.New(gohandle.Mainnet, m)
	handles, err := c.LookupHandles("addr1xyz")
	assert.Nil(t, err)
	assert.ElementsMatch(t, handles, []string{"abc", "xyz"})
}
