package jwt

import (
	"testing"
	"time"
)

func TestCreateAndParse(t *testing.T) {
	p := NewProvider(Config{Secret: "test-secret", ValiditySeconds: 3600})
	token, err := p.CreateToken("admin", "abc123", false)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := p.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Login != "admin" || claims.AuthToken != "abc123" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestValidateWrongSecret(t *testing.T) {
	p1 := NewProvider(Config{Secret: "a", ValiditySeconds: 60})
	p2 := NewProvider(Config{Secret: "b", ValiditySeconds: 60})
	token, _ := p1.CreateToken("u", "t", false)
	if p2.ValidateToken(token) {
		t.Fatal("expected invalid with wrong secret")
	}
}

func TestExpiredToken(t *testing.T) {
	p := NewProvider(Config{Secret: "s", ValiditySeconds: 1})
	p.validity = -time.Hour
	token, err := p.CreateToken("u", "t", false)
	if err != nil {
		t.Fatal(err)
	}
	if p.ValidateToken(token) {
		t.Fatal("expected expired token to fail")
	}
}
