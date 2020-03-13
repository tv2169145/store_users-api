package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/tv2169145/store_oauth-go/oauth/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-Caller-Id"

	paramAccessToken = "access_token"
)

var (
	oauthRestClient = resty.New()

)

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(request.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func GetClientId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return clientId
}

func AuthenticateRequest(request *http.Request) *errors.RestErr {
	if request == nil {
		return nil
	}
	cleanRequest(request)
	// http://api.store.com/resource?access_token=abc123
	accessTokenId := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	if accessTokenId == "" {
		return nil
	}
	token, err := getAccessToken(accessTokenId)
	if err != nil {
		if err.Status == http.StatusNotFound {
			return nil
		}
		return err
	}
	request.Header.Add(headerXClientId, fmt.Sprintf("%v", token.ClientId))
	request.Header.Add(headerXCallerId, fmt.Sprintf("%v", token.UserId))
	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}
	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(accessTokenId string) (*accessToken, *errors.RestErr) {
	oauthRestClient.SetHostURL("http://localhost:8080")
	oauthRestClient.SetTimeout(1 * time.Minute)
	response, err := oauthRestClient.R().Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenId))
	if err != nil {
		return nil, errors.NewInternalServerError("invalid restclient response when trying to get access token")
	}
	if response.StatusCode() > 299 {
		var restErr errors.RestErr
		if err := json.Unmarshal(response.Body(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("invalid response body when unmarshal response to restErr")
		}
		return nil, &restErr
	}
	var token accessToken
	if err := json.Unmarshal(response.Body(), &token); err != nil {
		return nil, errors.NewInternalServerError("invalid response body when unmarshal response to token")
	}
	return &token, nil
}
