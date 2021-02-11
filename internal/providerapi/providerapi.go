package providerapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-provider-admin/internal/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
	"github.com/bitly/go-simplejson"
)

const (
	apiURLStart string = "http://" 
	apiURLEnd string = "/v1" 
//	apiURLEnd string = "/client/establishment" 
)


const (
	clientAPIProtocolVersion = 1
	clientAPIProtocolSupportedHi = 1
)
// Can't have constant slices so create this at runtime.
// Order the API versions from most desirable to least desirable.
var clientAPIProtocolSupported []int



// Comms holds the communications specific data
type Comms struct {
	ApiURL string
	providerPubKey *fcrcrypto.KeyPair
	providerPubKeyVer *fcrcrypto.KeyVersion
	settings *settings.ClientSettings
}

// NewProviderAPIComms creates a connection with a provider
func NewProviderAPIComms(providerInfo *register.ProviderRegister, settings *settings.ClientSettings) (*Comms, error){
	hostAndPort := providerInfo.NetworkInfo

	// Create the constant array.
	if (clientAPIProtocolSupported == nil) {
		clientAPIProtocolSupported = make([]int, 1)
		clientAPIProtocolSupported[0] = clientAPIProtocolSupportedHi
	}

	// Check that the host name is valid
	err := validateHostName(hostAndPort)
	// if (err != nil) {
	// 	logging.Error("Host name invalid: %s", err. Error())
	// 	return nil, err
	// }

	netComms := Comms{}
	netComms.ApiURL = apiURLStart + hostAndPort + apiURLEnd

	signingKeyStr := providerInfo.SigingKey
	if len(signingKeyStr) > 2 && signingKeyStr[:2] == "0x" {
		runes := []rune(signingKeyStr)
		signingKeyStr = string(runes[2:])
	}

	netComms.providerPubKey, err = fcrcrypto.DecodePublicKey(signingKeyStr)
	if err != nil {
		logging.Error("Unable to decode public key: %v", err)
		return nil, err
	}
	netComms.providerPubKeyVer = fcrcrypto.DecodeKeyVersion(1) // TODO providerInfo.ProviderRetrievalPublicKeyVersion
	netComms.settings = settings
	return &netComms, nil
}


// ProviderCall calls the Provider's REST API
func (c *Comms) providerCall(msg interface{}) (*simplejson.Json) {

	// Create HTTP request.
	mJSON, _ := json.Marshal(msg)
	logging.Info("JSON sent: %s", string(mJSON))
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", c.ApiURL, contentReader)
	req.Header.Set("Content-Type", "application/json")

	// Send request and receive response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.ErrorAndPanic("Client - Provider communications (%s): %s", c.ApiURL, err)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	logging.Info("response body: %s", string(data))

	js, err := simplejson.NewJson(data)
	if err != nil {
		logging.ErrorAndPanic("Error decoding JSON: %s", err.Error())
	}

	return js
}


func validateHostName(host string) error {
	if len(host) == 0 {
		return errors.New("Error: Host name empty")
	} 

	ra, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return errors.New("Error: DNS Look-up failed for host: " + host)
	}
	logging.Info("Resolved %s as %s\n", host, ra.String())
	return nil
}


func (c *Comms) verifyMessage(signature string, wholeMessage interface{}) bool {
	keyVersion, err := fcrcrypto.ExtractKeyVersionFromMessage(signature)
	if err != nil {
		logging.Warn("Error decodign signature: %+v", err)
		return false
	}
	if keyVersion.NotEquals(c.providerPubKeyVer) {
		// TODO need to allow for multiple key versions, and fetch correct key
		// TODO based on version.
		logging.Error("Unknown Key Version used by provider: %d", keyVersion.EncodeKeyVersion())
	}

	verified, err := fcrcrypto.VerifyMessage(c.providerPubKey, signature, wholeMessage)
	if err != nil {
		logging.Warn("Signature verification error: %+v", err)
		return false
	}
	return verified
}