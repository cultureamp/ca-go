package jwt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDecoder(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	require.NoError(t, err)
	validJwks := string(b)

	testCases := []struct {
		desc           string
		jwks           string
		defaultKid     string
		expectedErrMsg string
	}{
		{
			desc:           "Success 1: valid jwks and found kid",
			jwks:           validJwks,
			defaultKid:     "web-gateway",
			expectedErrMsg: "",
		},
		{
			desc:           "Error 1: missing key",
			jwks:           "",
			defaultKid:     "",
			expectedErrMsg: "missing jwks",
		},
		{
			desc:           "Error 2: JWKS json bad",
			jwks:           "{\"bad\": \"jwks-json\" }",
			defaultKid:     "missing-kid",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
		{
			desc:           "Error 3: JWKS json invalid",
			jwks:           "invalid JSON",
			defaultKid:     "missing-kid",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
		{
			desc:           "Error 4: missing default kid",
			jwks:           validJwks,
			defaultKid:     "missing-kid",
			expectedErrMsg: "missing default key in JWKS",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoderWithDefaultKid(tC.jwks, tC.defaultKid)
			if tC.expectedErrMsg != "" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
				assert.Nil(t, decoder)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, decoder)
			}
		})
	}
}

func TestDecoderDecodeAllClaims(t *testing.T) {
	// useful to create RS256 test tokens https://jwt.io/

	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	require.NoError(t, err)

	testCases := []struct {
		desc            string
		token           string
		defaultKid      string
		expectedErrMsg  string
		accountId       string
		realUserId      string
		effectiveUserId string
		year            int
	}{
		{
			desc:            "Success 1: valid token with default kid",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMDg5NTI4NjEsIm5iZiI6MTU3NzgwMDg2MSwiaWF0IjoxNTc3ODAwODYxfQ.wV6z_kUjsKUebT7wcyTLqelUGhGfNA_82I_3lSsZOutegU4ct_652tp8enACcpgv2ZRmhbmCIC5w7PQn4rGLPx4Bdffjvf_4HvCqtvig0JV-0lmCpbNaadK93kiYNteZYFjLokLRKEHDt-uOoQbiWhc8DQBn7KbebLBRqp28HF28WL4-WVPDFsQ6H6pWL1RsXiuGyY8pI1y5b02t8-mte7CzrVx6uBgHPvfGgzWiTw4WpauxOxXUWTBIfK34OmPLb5sJQjrhM9RysE76j9703ptVfygTpokCcit-v_K3XlzQWw0T9sVfOu35mOS-NtXPJLB4-PK__gR60nANB-nNMsFf2Z1_ok44GAKE3an5Bi7cEaM-S5ZSbkq0rm6gbxEVZT5yjmJMNeecSgKc1dt_TAK5VMg5SJKzxXJ1DhLvKIB3rVNLyNfZVJT3mQW5NIBMjmfZad69_cu2TtS2_b8jOso7C7Vc3V-rB7MWYLsS47SDA46HFcAJvq7vsUHWM7POhoZEdSyN0cnpw0pWEOnhtpguJIw5XtrLQX00h2FytTAZEBnBUBYU-4AMQUBjK_FeEv5zDXSiXtR-iMs1YM0Qryhcw3Cx2EInkF8qt5AGNp6mYMwsFLJiv0RMO0CxxZ-08uZsng3Ekv8ewquKsXR5ZjSvtHO0vYSQEUcYvRUjELE",
			defaultKid:      "web-gateway",
			expectedErrMsg:  "",
			accountId:       "abc123",
			realUserId:      "xyz234",
			effectiveUserId: "xyz345",
			year:            2040,
		},
		{
			desc:            "Success 2: valid token with different default kid",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMDg5NTI4NjEsIm5iZiI6MTU3NzgwMDg2MSwiaWF0IjoxNTc3ODAwODYxfQ.wV6z_kUjsKUebT7wcyTLqelUGhGfNA_82I_3lSsZOutegU4ct_652tp8enACcpgv2ZRmhbmCIC5w7PQn4rGLPx4Bdffjvf_4HvCqtvig0JV-0lmCpbNaadK93kiYNteZYFjLokLRKEHDt-uOoQbiWhc8DQBn7KbebLBRqp28HF28WL4-WVPDFsQ6H6pWL1RsXiuGyY8pI1y5b02t8-mte7CzrVx6uBgHPvfGgzWiTw4WpauxOxXUWTBIfK34OmPLb5sJQjrhM9RysE76j9703ptVfygTpokCcit-v_K3XlzQWw0T9sVfOu35mOS-NtXPJLB4-PK__gR60nANB-nNMsFf2Z1_ok44GAKE3an5Bi7cEaM-S5ZSbkq0rm6gbxEVZT5yjmJMNeecSgKc1dt_TAK5VMg5SJKzxXJ1DhLvKIB3rVNLyNfZVJT3mQW5NIBMjmfZad69_cu2TtS2_b8jOso7C7Vc3V-rB7MWYLsS47SDA46HFcAJvq7vsUHWM7POhoZEdSyN0cnpw0pWEOnhtpguJIw5XtrLQX00h2FytTAZEBnBUBYU-4AMQUBjK_FeEv5zDXSiXtR-iMs1YM0Qryhcw3Cx2EInkF8qt5AGNp6mYMwsFLJiv0RMO0CxxZ-08uZsng3Ekv8ewquKsXR5ZjSvtHO0vYSQEUcYvRUjELE",
			defaultKid:      "test-other",
			expectedErrMsg:  "",
			accountId:       "abc123",
			realUserId:      "xyz234",
			effectiveUserId: "xyz345",
			year:            2040,
		},
		{
			desc:            "Error 1: invalid token with default kid",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6ld2F50.eyJhY2g2MSwiaWF0IjoxNTc3ODAwODYxfQ.wV6z_kUjsKUebT7RUjELE",
			defaultKid:      "test-other",
			expectedErrMsg:  "token is malformed",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
		{
			desc:            "Error 2: non RSA token with default kid",
			token:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			defaultKid:      "test-other",
			expectedErrMsg:  "unexpected signing method",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
		{
			desc:            "Error 3: RSA token with wrong kid",
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ",
			defaultKid:      "test-other",
			expectedErrMsg:  "token signature is invalid",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
		{
			desc:            "Error 4: RSA token with missing kid",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6Im1pc3Npbmcta2lkIiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMDg5NTI4NjEsIm5iZiI6MTU3NzgwMDg2MSwiaWF0IjoxNTc3ODAwODYxfQ.gpIoU05GPVaCdv1X9HTMH7Wotrlhwa_O-RpTTnN0H9cjYBvYICO25hUs4kw0ipV1F3QMo-k7a4F9FkmH2RL-gMp93wfFNpldaIASUqvrHZ3jBTDEKpk1VDD7YWDp8eqYkLetVSN0nui16U1HSOf2_Aw1vKqZUqOQGNYtjgNavPfSaBposa_ujfsA-On4guYwf2QDrgIuuOrJjoTys6mxaV5LBklKFJq_F5DrGTo_dBr9CSBUV0eIW6YLMpVuBBAJeLTAk3hFbxqz1Lcy68dW6_1bfhj9znuizgL5DMjXFNoKO9EXct_0Q8xa7RckT3Eu72ZPHWImUl1L9KepIM3N61SRdM7EOaLcT3G39a1Bjk1NLeNLA6S-m-dKAPDdkaGjF4yTTtsEmZS-W2gINuUyKs2Z_IbagDpmd1vAFMyl7MMs6Ul4Z9Rheuu_eHd4XKjgoLwd85hymqcKED70XuLSgIak0vtVXWixRDEAYvyeHlrFsZc4Vo4X_q1SyxIbwj1laDC0Y8Y7O_LnGrOH7D5DmVbCXbm1vJUJ7U98K8JRrv5y7ZETofvyLqSE4FY2aAiE7tc9BuTDOjIcEiUPLxEUSB6c08L53FmJCehdB80OR1Fv24G_vothsTnRtsBpiXVTQPpbS5Po_CNl-sFQBKViRsE6IxQt81HhoEJ1Ui4A0r8",
			defaultKid:      "test-other",
			expectedErrMsg:  "token signature is invalid",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoderWithDefaultKid(string(b), tC.defaultKid)
			assert.Nil(t, err)
			assert.NotNil(t, decoder)

			claim, err := decoder.Decode(tC.token)

			if tC.expectedErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
			}

			assert.Equal(t, tC.accountId, claim.AccountId)
			assert.Equal(t, tC.realUserId, claim.RealUserId)
			assert.Equal(t, tC.effectiveUserId, claim.EffectiveUserId)
			assert.Equal(t, tC.year, claim.ExpiresAt.Year())
		})
	}
}
