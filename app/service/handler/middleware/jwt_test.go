package middleware

import (
	"app/global"
	"testing"

	"github.com/spf13/viper"
)

func TestGetToken(t *testing.T) {
	// Keep the test independent of app startup and an expired hard-coded token.
	previous := global.Global
	config := viper.New()
	config.Set("env.signkey", "unit-test-signing-key")
	config.Set("env.tokenex", 60)
	global.Global = &global.RunTime{Viper: config}
	t.Cleanup(func() { global.Global = previous })
	token, err := SetToken("test-phone", "test-user")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := GetTokenClaims(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UUID != "test-user" || claims.Phone != "test-phone" {
		t.Fatalf("unexpected token claims: %+v", claims)
	}
}
