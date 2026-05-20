package application

import sharedcrypto "github.com/gis-mdm/server-backend-go/internal/shared/crypto"

func checkRequestSignature(signature, value string) bool {
	return sharedcrypto.CheckRequestSignature(signature, value)
}

func signSyncResponse(secret string, payload any) string {
	return sharedcrypto.SignSyncResponse(secret, payload)
}
