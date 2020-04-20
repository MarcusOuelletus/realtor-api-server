package tokens

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/MarcusOuelletus/rets-server/database"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

// key must be length 16, 24 or 32
var key = []byte("MY_TOKEN")

type claims struct {
	BrokerToken string `json:"BrokerToken"`
	jwt.StandardClaims
}

type tokenClass struct {
	claim *claims
}

func new() *tokenClass {
	return &tokenClass{
		claim: &claims{},
	}
}

// getClientTokenFromRequest - returns the Authorization header which should be browser token
func (t *tokenClass) getClientTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")

	if token == "" {
		return "", errors.New("no token specified")
	}

	return token, nil
}

func (t *tokenClass) getJWTFromBrowserToken(browserToken string) (string, error) {
	encryptedToken := t.decodeEncryptedJWT(browserToken)
	jwtToken, err := t.decryptJWT(encryptedToken)

	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t *tokenClass) storeToken(encryptedClientToken string, ip string) error {
	err := database.Insert(&database.InsertQuery{
		Collection: "jwt",
		Data: map[string]interface{}{
			"Token":     encryptedClientToken,
			"ExpiresAt": t.claim.ExpiresAt,
			"IP":        ip,
		},
	})

	if err != nil {
		glog.Errorf("error inserting token: %s\n", err.Error())
		return err
	}

	return nil
}

func (t *tokenClass) deleteExpiredTokens() error {
	return database.Remove(&database.InsertQuery{
		Collection: "jwt",
		Data: map[string]interface{}{
			"ExpiresAt": bson.M{"$lt": time.Now().Unix()},
		},
	})
}

func (t *tokenClass) generateJWT(brokerToken string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes expiry time
	t.claim = &claims{
		BrokerToken: brokerToken,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.claim)

	// Create the JWT string
	tokenString, err := token.SignedString(key)

	if err != nil {
		glog.Errorln("failed to sign jwt")
		glog.Errorln(err.Error())
		return "", errors.New("failed to signed jwt")
	}

	return tokenString, nil
}

func (t *tokenClass) encodeEncryptedJWT(encryptedToken []byte) string {
	encodedToken := base64.URLEncoding.EncodeToString(encryptedToken)

	return encodedToken
}

func (t *tokenClass) decodeEncryptedJWT(encoded string) []byte {
	data, _ := base64.URLEncoding.DecodeString(encoded)

	return data
}

func (t *tokenClass) validateBrowserToken(browserToken string) error {
	jwtToken, err := t.getJWTFromBrowserToken(browserToken)

	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(jwtToken, t.claim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		// glog.Errorf("JWT: error validating token - %s\n", err.Error())
		return err
	}

	if !token.Valid {
		return errors.New("JWT: token is invalid")
	}

	return nil
}

func (t *tokenClass) encryptJWT(jwtToken string) ([]byte, error) {
	gcm, err := buildGCM()

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		glog.Errorf("JWT: error reading nonce - %s\n")
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(jwtToken), nil)

	return ciphertext, nil
}

func (t *tokenClass) decryptJWT(decodedJWT []byte) (string, error) {
	gcm, err := buildGCM()

	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	nonce, ciphertext := decodedJWT[:nonceSize], decodedJWT[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		glog.Errorf("JWT: error calling Open on gcm - %s\n", err.Error())
		return "", nil
	}

	return string(plaintext), nil
}

func buildGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		glog.Errorf("JWT: error in NewCipher - %s\n", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		glog.Errorf("JWT: error creating NewGCM - %s\n", err.Error())
		return nil, err
	}

	return gcm, nil
}
