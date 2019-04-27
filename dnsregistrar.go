// Copyright 2019 Weald Technology Trading
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ens

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wealdtech/go-ens/v2/contracts/dnsregistrar"
)

// DNSRegistrar is the structure for the registrar
type DNSRegistrar struct {
	client   *ethclient.Client
	domain   string
	Contract *dnsregistrar.Contract
}

// NewDNSRegistrar obtains the registrar contract for a given domain
func NewDNSRegistrar(client *ethclient.Client, domain string) (*DNSRegistrar, error) {
	address, err := RegistrarContractAddress(client, domain)
	if err != nil {
		return nil, err
	}

	if address == UnknownAddress {
		return nil, fmt.Errorf("no registrar for domain %s", domain)
	}

	contract, err := dnsregistrar.NewContract(address, client)
	if err != nil {
		return nil, err
	}

	// Ensure this really is a DNS registrar.  To do this confirm that it supports
	// the expected interface.
	supported, err := contract.SupportsInterface(nil, [4]byte{0x1a, 0xa2, 0xe6, 0x41})
	if err != nil {
		return nil, err
	}
	if !supported {
		return nil, fmt.Errorf("purported registrar for domain %s does not support DNS registrar functionality", domain)
	}

	return &DNSRegistrar{
		client:   client,
		domain:   domain,
		Contract: contract,
	}, nil
}

// TODO claim
