// Copyright 2026 The HuaTuo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"context"
	"fmt"
	"net"
)

// Resolver resolves host names while preparing a filter for compilation.
type Resolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
}

type hostResolutionError struct {
	host string
}

func (e hostResolutionError) Error() string {
	return fmt.Sprintf("unknown host: %s", e.host)
}

func (e hostResolutionError) Unwrap() error {
	return ErrHostResolution
}

func prepareFilter(f Filter, resolver Resolver) (Filter, error) {
	switch v := f.(type) {
	case primitive:
		return v.prepare(resolver)
	case arithmeticComparison:
		if err := validateArithmeticDepth(v.left, 0); err != nil {
			return nil, err
		}
		if err := validateArithmeticDepth(v.right, 0); err != nil {
			return nil, err
		}
		if 1+max(arithmeticScratch(v.left), arithmeticScratch(v.right)) > 16 {
			return nil, fmt.Errorf("arithmetic expression requires more than 16 scratch slots")
		}
		return v, nil
	case negated:
		inner, err := prepareFilter(v.inner, resolver)
		if err != nil {
			return nil, err
		}
		v.inner = inner
		return v, nil
	case composite:
		if len(v.filters) == 0 {
			return nil, ErrInvalidFilter
		}
		prepared := make(Filters, 0, len(v.filters))
		for _, child := range v.filters {
			filter, err := prepareFilter(child, resolver)
			if err != nil {
				return nil, err
			}
			prepared = append(prepared, filter)
		}
		v.filters = prepared
		return v, nil
	default:
		return nil, ErrInvalidFilter
	}
}

func (p primitive) prepare(resolver Resolver) (primitive, error) {
	if err := p.validate(); err != nil {
		return primitive{}, err
	}
	if p.kind != filterKindHost || p.protocol == filterProtocolEther {
		return p, nil
	}

	addrs, err := p.resolve(resolver)
	if err != nil {
		return primitive{}, err
	}
	if !addressesMatchProtocol(addrs, p.protocol) {
		return primitive{}, hostResolutionError{host: p.id}
	}
	p.resolved = addrs
	return p, nil
}

func addressesMatchProtocol(addrs []net.IP, protocol filterProtocol) bool {
	if protocol == filterProtocolUnset {
		return len(addrs) > 0
	}
	wantIPv4 := protocol == filterProtocolIP || protocol == filterProtocolArp || protocol == filterProtocolRarp
	for _, addr := range addrs {
		if wantIPv4 && addr.To4() != nil {
			return true
		}
		if protocol == filterProtocolIP6 && addr.To4() == nil {
			return true
		}
	}
	return false
}

func (p primitive) resolve(resolver Resolver) ([]net.IP, error) {
	if addr := net.ParseIP(p.id); addr != nil {
		return []net.IP{addr}, nil
	}
	if resolver == nil {
		resolver = net.DefaultResolver
	}

	resolved, err := resolver.LookupHost(context.Background(), p.id)
	if err != nil {
		return nil, hostResolutionError{host: p.id}
	}
	addrs := make([]net.IP, 0, len(resolved))
	for _, value := range resolved {
		if addr := net.ParseIP(value); addr != nil {
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		return nil, hostResolutionError{host: p.id}
	}
	return addrs, nil
}

func unsupportedFeature(name string) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedFeature, name)
}
