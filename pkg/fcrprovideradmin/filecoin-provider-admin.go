package fcrprovideradmin

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
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/control"
	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
)

// FilecoinRetrievalProviderAdminClient holds information about the interaction of
// the Filecoin Retrieval Provider Admin Client with Filecoin Retrieval Providers.
type FilecoinRetrievalProviderAdminClient struct {
	providerManager *control.ProviderManager
	// TODO have a list of provider objects of all the current providers being interacted with
}

var singleInstance *FilecoinRetrievalProviderAdminClient
var initialised = false

// InitFilecoinRetrievalProviderAdminClient initialise the Filecoin Retreival Client library
func InitFilecoinRetrievalProviderAdminClient(settings Settings) *FilecoinRetrievalProviderAdminClient {
	if initialised {
		log.ErrorAndPanic("Attempt to init Filecoin Retrieval Provider Admin Client a second time")
	}
	var c = FilecoinRetrievalProviderAdminClient{}
	c.startUp(settings)
	singleInstance = &c
	initialised = true
	return singleInstance

}

// GetFilecoinRetrievalProviderAdminClient creates a Filecoin Retrieval Provider Admin Client
func GetFilecoinRetrievalProviderAdminClient() *FilecoinRetrievalProviderAdminClient {
	if !initialised {
		log.ErrorAndPanic("Filecoin Retrieval Provider Admin Client not initialised")
	}

	return singleInstance
}

func (c *FilecoinRetrievalProviderAdminClient) startUp(conf Settings) {
	log.Info("Filecoin Retrieval Provider Admin Client started")
	clientSettings := conf.(*settings.ClientSettings)
	c.providerManager = control.NewProviderManager(*clientSettings)
}

// Shutdown releases all resources used by the library
func (c *FilecoinRetrievalProviderAdminClient) Shutdown() {
	log.Info("Filecoin Retrieval Provider Admin Client shutting down")
	c.providerManager.Shutdown()
}

// SendMessage send message to providers
func (c *FilecoinRetrievalProviderAdminClient) SendMessage(message *fcrmessages.FCRMessage) {
	log.Info("Filecoin Retrieval Provider Admin Client sending message")
	c.providerManager.SendMessage(message)
}
