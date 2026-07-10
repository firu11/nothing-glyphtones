package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestAuthWriteAndReadCookie(t *testing.T) {
	t.Parallel()

	auth := NewAuth("super-secret", false)
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	if err := auth.WriteAuthCookie(c, 42); err != nil {
		t.Fatalf("WriteAuthCookie() error = %v", err)
	}

	cookies := res.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected auth cookie to be written")
	}
	req.AddCookie(cookies[0])

	gotID := auth.GetIDFromCookie(c)
	if gotID != 42 {
		t.Fatalf("GetIDFromCookie() = %d, want 42", gotID)
	}
}

func TestAuthValidateTokenRoundTrip(t *testing.T) {
	t.Parallel()

	auth := NewAuth("super-secret", false)
	token, err := auth.generateToken(42)
	if err != nil {
		t.Fatalf("generateToken() error = %v", err)
	}

	valid, gotID, err := auth.validateToken(token)
	if err != nil {
		t.Fatalf("validateToken() error = %v", err)
	}
	if !valid {
		t.Fatal("validateToken() returned invalid token")
	}
	if gotID != 42 {
		t.Fatalf("validateToken() id = %d, want 42", gotID)
	}
}

func TestAuthRejectsTamperedCookie(t *testing.T) {
	t.Parallel()

	auth := NewAuth("super-secret", false)
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(newTestCookie(CookieName, "definitely-not-a-jwt"))
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	gotID := auth.GetIDFromCookie(c)
	if gotID != 0 {
		t.Fatalf("GetIDFromCookie() = %d, want 0 for invalid cookie", gotID)
	}
}

func newTestCookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value, Path: "/"}
}
