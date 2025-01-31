// Copyright Istio Authors
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

package cluster

import "sync"

// AddressMap provides a thread-safe mapping of addresses for each Kubernetes cluster.
type AddressMap struct {
	// Addresses hold the underlying map. Visible only for testing, for the purposes of simplified construction.
	// Production code should never access this directly.
	Addresses map[ID][]string

	// NOTE: The copystructure library is not able to copy unexported fields, so the mutex will not be copied.
	mutex sync.RWMutex
}

func (m *AddressMap) IsEmpty() bool {
	if m == nil {
		return true
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.Addresses) == 0
}

func (m *AddressMap) DeepCopy() AddressMap {
	return AddressMap{
		Addresses: m.GetAddresses(),
	}
}

// GetAddresses returns the mapping of clusters to addresses.
func (m *AddressMap) GetAddresses() map[ID][]string {
	if m == nil {
		return nil
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.Addresses == nil {
		return nil
	}

	out := make(map[ID][]string)
	for k, v := range m.Addresses {
		out[k] = append([]string{}, v...)
	}
	return out
}

func (m *AddressMap) GetAddressesFor(c ID) []string {
	if m == nil {
		return nil
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.Addresses == nil {
		return nil
	}

	// Copy the Addresses array.
	return append([]string{}, m.Addresses[c]...)
}

func (m *AddressMap) SetAddressesFor(c ID, addresses []string) *AddressMap {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create the map if nil
	if m.Addresses == nil {
		m.Addresses = make(map[ID][]string)
	}

	if len(addresses) == 0 {
		delete(m.Addresses, c)
	} else {
		m.Addresses[c] = addresses
	}
	return m
}

func (m *AddressMap) AddAddressesFor(c ID, addresses []string) *AddressMap {
	if len(addresses) == 0 {
		return m
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create the map if nil.
	if m.Addresses == nil {
		m.Addresses = make(map[ID][]string)
	}

	m.Addresses[c] = append(m.Addresses[c], addresses...)
	return m
}

func (m *AddressMap) ForEach(fn func(c ID, addresses []string)) {
	if m == nil {
		return
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.Addresses == nil {
		return
	}

	for c, addresses := range m.Addresses {
		fn(c, addresses)
	}
}
