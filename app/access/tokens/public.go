package tokens

import (
	"net/http"

	"github.com/MarcusOuelletus/rets-server/helpers"
)

type TokenData struct {
	BrokerToken string
}

// Parse - takes a request object, returns a struct containing the JWT data, returns an error if the token's expired
func Parse(r *http.Request) (*TokenData, error) {
	t := &tokenClass{
		claim: &claims{},
	}

	browserToken, err := t.getClientTokenFromRequest(r)

	if err != nil {
		return nil, err
	}

	if err := t.validateBrowserToken(browserToken); err != nil {
		return nil, err
	}

	tokenData := &TokenData{
		BrokerToken: t.claim.BrokerToken,
	}

	return tokenData, nil
}

// CreateClientToken - returns a base64 encoded jwt token encrypted with AES-GCM
func CreateClientToken(brokerToken string, remoteAddr string) (string, error) {
	var t = new()

	ip := helpers.ParseIP(remoteAddr)

	clientJWT, err := t.generateJWT(brokerToken)

	if err != nil {
		return "", err
	}

	err = t.storeToken(clientJWT, ip)

	if err != nil {
		return "", err
	}

	encryptedToken, err := t.encryptJWT(clientJWT)

	if err != nil {
		return "", err
	}

	encodedToken := t.encodeEncryptedJWT(encryptedToken)

	go t.deleteExpiredTokens()

	return encodedToken, nil
}
