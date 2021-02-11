package control

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	// "fmt"
	"sync"

	// "github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	// "github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	// "github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/providerapi"
)

// ProviderManager managers the pool of providers and the connections to them.
type ProviderManager struct {
	settings     settings.ClientSettings
	providers     []ActiveProvider
	providersLock sync.RWMutex

	// Registered Providers
	RegisteredProviders []register.ProviderRegister
}

// ActiveProvider contains information for a single provider
type ActiveProvider struct {
	info  register.ProviderRegister
	comms *providerapi.Comms
}

// NewProviderManager returns an initialised instance of the provider manager.
func NewProviderManager(settings settings.ClientSettings) *ProviderManager {
	g := ProviderManager{}
	g.settings = settings
	g.providerManagerRunner()
	return &g
}

// TODO this should be in a go routine and loop for ever.
func (g *ProviderManager) providerManagerRunner() {
	logging.Info("Provider Manager: Management thread started")

	// Call this once each hour or maybe day.
	providers, err := register.GetRegisteredProviders(g.settings.RegisterURL())
	if err != nil {
		logging.Error("Unable to get registered providers: %v", err)
	}
	g.RegisteredProviders = providers

	// TODO this loop is where the managing of providers that the client is using happens.
	logging.Info("Provider Manager: GetProviders returned %d providers", len(providers))
	for _, provider := range providers {
		logging.Info("Setting-up comms with: %+v", provider)
		comms, err := providerapi.NewProviderAPIComms(&provider, &g.settings)
		if err != nil {
			panic(err)
		}

		activeProvider := ActiveProvider{provider, comms}
		g.providers = append(g.providers, activeProvider)
	}

	logging.Info("Provider Manager using %d providers", len(g.providers))
}

// BlockProvider adds a host to disallowed list of providers
func (g *ProviderManager) BlockProvider(hostName string) {
	// TODO
}

// UnblockProvider add a host to allowed list of providers
func (g *ProviderManager) UnblockProvider(hostName string) {
	// TODO

}

// FindOffersStandardDiscovery finds offers using the standard discovery mechanism.
// func (g *ProviderManager) FindOffersStandardDiscovery(contentID *cid.ContentID) ([]cidoffer.CidGroupOffer, error) {
// 	if len(g.providers) == 0 {
// 		return nil, fmt.Errorf("No providers available")
// 	}

// 	var aggregateOffers []cidoffer.CidGroupOffer
// 	for _, gw := range g.providers {
// 		// TODO need to do nonce management
// 		// TODO need to do requests to all providers in parallel, rather than serially
// 		offers, err := gw.comms.ProviderStdCIDDiscovery(contentID, 1)
// 		if err != nil {
// 			logging.Warn("ProviderStdDiscovery error. Provider: %s, Error: %s", gw.info.NodeID, err)
// 		}
// 		// TODO: probably should remove duplicate offers at this point
// 		aggregateOffers = append(aggregateOffers, offers...)
// 	}
// 	return aggregateOffers, nil
// }

// GetConnectedProviders returns the list of domain names of providers that the client
// is currently connected to.
func (g *ProviderManager) GetConnectedProviders() []string {
	urls := make([]string, len(g.providers))
	for i, provider := range g.providers {
		urls[i] = provider.comms.ApiURL
	}
	return urls
}

// Shutdown stops go routines and closes sockets. This should be called as part
// of the graceful library shutdown
func (g *ProviderManager) Shutdown() {
	// TODO
}
