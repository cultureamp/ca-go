package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	day  = 24 * time.Hour
	year = 365 * day // approx

	decoderAuthKey       string = "./testKeys/jwt.rs256.key.development.pub"
	decoderSecondAuthKey string = "./testKeys/jwt.rs256.key.development.extra_1.pub"
	decoderThirdAuthKey  string = "./testKeys/jwt.rs256.key.development.extra_2.pub"
	decoderAuthJwks      string = "./testKeys/development.jwks"
)

func TestNewDecoderSuccess(t *testing.T) {
	pubKeyBytes, err := os.ReadFile(filepath.Clean(decoderAuthKey))
	require.NoError(t, err)

	pubSecondKeyBytes, err := os.ReadFile(filepath.Clean(decoderSecondAuthKey))
	require.NoError(t, err)

	pubJwkKeyBytes, err := os.ReadFile(filepath.Clean(decoderAuthJwks))
	require.NoError(t, err)

	decoder, err := NewJwtDecoder(string(pubKeyBytes), string(pubSecondKeyBytes), string(pubJwkKeyBytes))
	assert.Nil(t, err)
	assert.NotNil(t, decoder)
}

func TestNewDecoderErrors(t *testing.T) {
	invalidPublicKey := "invalid-public-key"

	pubKeyBytes, err := os.ReadFile(filepath.Clean(decoderAuthKey))
	require.NoError(t, err)

	pubSecondKeyBytes, err := os.ReadFile(filepath.Clean(decoderSecondAuthKey))
	require.NoError(t, err)

	testCases := []struct {
		desc           string
		pubKey         string
		secondPubKey   string
		jwks           string
		expectedErrMsg string
	}{
		{
			desc:           "Error 1: missing key",
			pubKey:         "",
			secondPubKey:   "",
			jwks:           "",
			expectedErrMsg: "invalid key",
		},
		{
			desc:           "Error 2: bad key",
			pubKey:         invalidPublicKey,
			secondPubKey:   invalidPublicKey,
			jwks:           invalidPublicKey,
			expectedErrMsg: "invalid key",
		},
		{
			desc:           "Error 3: keys ok, JWKS json bad",
			pubKey:         string(pubKeyBytes),
			secondPubKey:   string(pubSecondKeyBytes),
			jwks:           "{\"bad\": \"jwks-json\" }",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
		{
			desc:           "Error 3: keys ok, JWKS json invalid",
			pubKey:         string(pubKeyBytes),
			secondPubKey:   string(pubSecondKeyBytes),
			jwks:           "invalid JSON",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoder(tC.pubKey, tC.secondPubKey, tC.jwks)
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tC.expectedErrMsg)
			assert.Nil(t, decoder)
		})
	}
}

func TestDecoderDecodeAllClaims(t *testing.T) {
	pubKeyBytes, err := os.ReadFile(filepath.Clean(decoderAuthKey))
	require.NoError(t, err)

	pubSecondKeyBytes, err := os.ReadFile(filepath.Clean(decoderSecondAuthKey))
	require.NoError(t, err)

	pubJwksBytes, err := os.ReadFile(filepath.Clean(decoderAuthJwks))
	require.NoError(t, err)

	decoder, err := NewJwtDecoder(string(pubKeyBytes), string(pubSecondKeyBytes), string(pubJwksBytes))
	assert.Nil(t, err)
	assert.NotNil(t, decoder)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo"
	claim, err := decoder.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", claim.AccountId)
	assert.Equal(t, "xyz234", claim.RealUserId)
	assert.Equal(t, "xyz345", claim.EffectiveUserId)
	assert.Equal(t, 2030, claim.Expiry.Year())
}

func TestDecoderConfigureClaims(t *testing.T) {
	pubBytes, err := os.ReadFile(filepath.Clean(decoderAuthKey))
	require.NoError(t, err)

	jwkBytes, err := os.ReadFile(filepath.Clean(decoderAuthJwks))
	require.NoError(t, err)

	testCases := []struct {
		desc            string
		tokenString     string
		accountID       string
		realUserID      string
		effectiveUserID string
	}{
		{
			desc:            "Success 1: default with extra machine key",
			tokenString:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo",
			accountID:       "abc123",
			realUserID:      "xyz234",
			effectiveUserID: "xyz345",
		},
		{
			desc:            "Success 2: extra machine kid key",
			tokenString:     "eyJhbGciOiJSUzI1NiIsImtpZCI6ImV4dHJhXzEiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiJkZWY0NTYiLCJlZmZlY3RpdmVVc2VySWQiOiJhc2RmMTIzNCIsInJlYWxVc2VySWQiOiJhc2RmMTIzNCIsImV4cCI6MjQ4MDIyMTg0MywiaWF0IjoxNjkxODIxODQzfQ.eLLMWZJBlqHlpa8GhBMDWvhJSbOfqObxEIpuyr6sidBLTZ92zTm4H9aAUxd4qZITYjvH12JT8JmpJWEWnPlE6j5TgHWCVpIL1BCIoVhQruOkw9RC10kJ5pwRE1VpLooLrqR0cxDAGnD1pTQugPRJnqGqmT71Yqt-0ZBa-T3jJxOUHhTKHQ_PFJcrCP37Htx_iijtg0c9ej_pdplf_0q9fln0lqWAWQieI7YMd3SS3FDdHAT9YGfiRJkFHve7YBTiUruQFVthW0V2Us8Shsx8wtlZVc0XLw7E68PwSHkyRkZWV45PPHDr7z6q_9_iTHKyCGpjAG458jyn038sAMkz1ni-wY6gCGK7IIwUJ3mfquoSkfomnGBVz89GlTjJYeyep-22uK8j-PjCKeeUK2laEjH5Lxmk5mBRAWRASUTSOJNkSPF8pZBvGVSeoIW7hF1VL5FgQRl4ZvCRMFq2E78SoawbGcwb-YBQXlsD1fgGru6Rza5IzYc9iuZNRXLFSRDpv1c3gJeJiZqTb0yX8CopHgfNdXd0krKgLjLiLytaRfhSUKgc26_NPWIw4SlivVu21veDCPXzNXJCflfxB-iINN-COH6LobmhuXpYbVt0FbEV61L3AfoP0Yu7EkumHV9WvwhFd7LVh3T9dDYvmkhMO8eY4D23xN2elb76dGXN68o",
			accountID:       "def456",
			realUserID:      "asdf1234",
			effectiveUserID: "asdf1234",
		},
		{
			desc:            "Success 3: extra machine2 kid key",
			tokenString:     "eyJhbGciOiJSUzI1NiIsImtpZCI6ImV4dHJhXzIiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiIxYSIsImVmZmVjdGl2ZVVzZXJJZCI6IjJiIiwicmVhbFVzZXJJZCI6IjNjIiwiZXhwIjoyNDgwMjIxOTUxLCJpYXQiOjE2OTE4MjE5NTF9.1ahd21wx2WNTRMHPbqjOmv6vEm52RUU8PfQbRdg_xuThCHgUYkVvxG5H7J35KPZlV62yJVO5-G3qNnhfLzoZOS2B_aXriY0YoXuEi22g6g6xwlVHaRhyl0iVDCEViu5c-UX2QKqgZdHhfnSmKPMpni4mbbgHdx4Ki5_nBroRrw_48vQqkJd15fcF5MiEiso8swyoF9qh2c0bbEKf67hQQeO23wDbHaMPd3rtrRHMZAfTPvrKh0tul1_jEP_yJJBIPxDHS71iI1Tp9ohzDmiuDkAEUtFfsWbePS5COJzGI77eU6-w11W6VjeQd2wUUSOqqB8R9qrvE4h8QTGsWS8UMECIseTM-ZhKCZW_dwzMmu9VIHQiscJZ3r2PQZ3q5vDm0NpvZulykXzH4y7uaG6YTG158Px7OVOh3ysIJSQGSfYU5Sk38k08mc-izcSKI98tFE1EDLnXifK7wNhEPwNKBr0pDXBWuS-lltttMzEBrDlCd3Xj6LDeNMNkMYpHkyM1Q_kL6XnIZNv-apjkTAvmZ2nVy2j3HLOYrunBWJUx8qXqGVUGNgmuPXosTvgBpLCNoBul5MK6O0U7epeUClqpVkRZ-MAptfxU6yD5oP8Ry5Ig2K1PRJ3Cxikeeb6IZTUFH3GLBfM5Kmw9GNjrkk-WJmEoa0fzFxyvaGCfgS64ASU",
			accountID:       "1a",
			realUserID:      "3c",
			effectiveUserID: "2b",
		},
		{
			desc:            "Success 4: default decode with machine3 kid key",
			tokenString:     "eyJhbGciOiJSUzI1NiIsImtpZCI6ImV4dHJhXzMiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoyNDgwMjIyMDQ0LCJpYXQiOjE2OTE4MjIwNDR9.R1xlASH7izLCu6CdaRWhZrUZ6IdZwjkpLP0KFyDS3adVmbTPItzdJLbfQGGiSeCFrWb6r7MpjoinAzl-4E4szj7TxCtCj4vlNvharituHd_CO-AZ09xM0Dx25Ie1PtMinRKelEqV0k2sQ1txu9K9H88S-8LIHyMoVoIiT6I4kHC6bTHN2BJAmNszcByHo-D7ckdf7RDb-oHnR_oSr-budclK2uWBiWVKvwtK4tLIh3TXUBYjAifRVzKnA666axwRY1IqYVkonDV93Al-5sHQOQR8UFfVQ65D2VeNz-UydAFT0zdr_ltn8S7joRoWZMfeegsWubfq5v8J53Z4L0rttu2wNjQpYWnwFfP4KpbOgPNxjPWNAjdEwCB2NaXi5Y6Eghrc4sPsa3bHCOA-S_4vvnH7vr4bRxy-AUW90Ukms_TTTo1KdIS21V_k5G9cWVYEoNBWbfsqrUyyZypZDGdWa06uapW1NaT3gxA0CFbmq59G77PpKLWFtBg0eIsqlq_wKT-iX7SBUzg_Wx6gvyZd9xr2_3i9-KZzXH7fttFjo9uIaBpWzPxaSWO_i9noVrH4qRCuaYKqpb0E9vi9GwtuTdRzJvHimTEpGz_2vgg2iYpLkkHy8EuAlRNzwowOkorocNka7G_rikZQLS7fWGvrZW9mHJGlA7r0j18SoHqgL3g",
			accountID:       "abc123",
			realUserID:      "xyz234",
			effectiveUserID: "xyz345",
		},
		{
			desc:            "Success 5: default with jwks machine key",
			tokenString:     "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo",
			accountID:       "abc123",
			realUserID:      "xyz234",
			effectiveUserID: "xyz345",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoder(string(pubBytes), "", string(jwkBytes))
			assert.Nil(t, err)
			claim, err := decoder.Decode(tC.tokenString)
			assert.Nil(t, err)
			assert.Equal(t, tC.accountID, claim.AccountId)
			assert.Equal(t, tC.realUserID, claim.RealUserId)
			assert.Equal(t, tC.effectiveUserID, claim.EffectiveUserId)
		})
	}
}
