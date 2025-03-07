package verifier

import (
	"testing"
	"yip/src/config"
)

func getVerifier(t *testing.T) Verifier {
	return NewVerifier(
		config.JWTTokenConfig{
			TokenExpirationInSec:        40,
			RefreshTokenExpirationInSec: 10000,
			Issuer:                      "issuer",
			CertificatePrivate:          "../../test_certs/app.rsa",
			CertificatePublic:           "../../test_certs/app.rsa.pub",
		})
}

func TestVerifierEnvironment(t *testing.T) {
	privKey := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
	pubKey := "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2Igp95bgFC564/aVc70X\n0Dh+oQsXSeLJ3dY8PgVsjmvAdTngT868O72S7Crr9G88QH7cZDaGGc475pY1kc6H\n8A93b5VKDBELeK5qT2tHaQv6i8k2hrqE8d4c/VkURDsIyzha/of27ZacOKuhyxZj\nO/wu81DkDIl+drRRE3s0ZPNfuTRwigSixIbM/YdrkW2R7q0x3D3CeTyQOplS8gS/\nM+WF/fwvmLuAKp18l6BXjBWdXm6tsvI5+pUT5J/FLns6PuCcMfUSTYNV6ckLGM6F\nkvEvB2PjlWh8sumC+uxNZD4n8894mzAAQ/1ZqCJ8dbaTdBob/rYCfc3oTr6+KAbP\nRVsWAVV+zT/fHnyzfRZOymCRKzOgj1CgsBukYfTubZtAQzs9/dR8EJsJ1w8ASyFQ\n4YHTCjpy5ivkqTow1mrWuvuQ9zU2JsddwiIXk4Qq1y4RwaVDY/EYR17r3i3Leulp\nY1fpzyqf4kvOBrttnnM2jL6+IJeGrI/rZ1t9VyWqpAP1pYE/HopnaufF1yWs7utz\n5BhL0Ono/IrEoRzYnsQDueH+rtLLLSpIns5hyE93XlYpzSM31ojJ6LwJhSJMly3Q\nmT/ctDwFuQY/kPhxfwlqmJCDdGMK/H2V69TrjSsVwe0KZDlYrn+B+8yyKlPvA+N0\n3fCBJTn3tgICEazl899ft/sCAwEAAQ==\n-----END PUBLIC KEY-----\n"
	NewVerifier(
		config.JWTTokenConfig{
			TokenExpirationInSec:        40,
			RefreshTokenExpirationInSec: 10000,
			Issuer:                      "issuer",
			CertificatePrivate:          privKey,
			CertificatePublic:           pubKey,
		})

}
