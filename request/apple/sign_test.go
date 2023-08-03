package apple

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestGetUniqueID(t *testing.T) {
	tests := []struct {
		name    string
		idToken string
		want    string
		wantErr bool
	}{
		{
			name:    "successful decode",
			idToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiY29tLmV4YW1wbGUuYXBwIiwiZXhwIjoxNTY4Mzk1Njc4LCJpYXQiOjE1NjgzOTUwNzgsInN1YiI6IjA4MjY0OS45MzM5MWQ4ZTExOTJmNTZiOGMxY2gzOWdzMmE0N2UyLjk3MzIiLCJhdF9oYXNoIjoickU3b3Brb1BSeVBseV9Pc2Rhc2RFQ1ZnIiwiYXV0aF90aW1lIjoxNTY4Mzk1MDc2fQ.PR3mMoVMdJo8EGPy6_aJ3sJGwAgcnnFjt9UCRXqWerI",
			want:    "082649.93391d8e1192f56b8c1ch39gs2a47e2.9732",
			wantErr: false,
		},
		{
			name:    "bad token",
			idToken: "badtoken",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUniqueId(tt.idToken)
			if !tt.wantErr {
				assert.NoError(t, err, "expected no error but received %s", err)
			}

			if tt.want != "" {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetClaims(t *testing.T) {
	tests := []struct {
		name      string
		idToken   string
		wantEmail string
		wantErr   bool
	}{
		{
			name:      "successful decode",
			idToken:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiY29tLmV4YW1wbGUuYXBwIiwiZXhwIjoxNTY4Mzk1Njc4LCJpYXQiOjE1NjgzOTUwNzgsInN1YiI6IjA4MjY0OS45MzM5MWQ4ZTExOTJmNTZiOGMxY2gzOWdzMmE0N2UyLjk3MzIiLCJhdF9oYXNoIjoickU3b3Brb1BSeVBseV9Pc2Rhc2RFQ1ZnIiwiZW1haWwiOiJmb29AYmFyLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjoidHJ1ZSIsImlzX3ByaXZhdGVfZW1haWwiOiJ0cnVlIiwiYXV0aF90aW1lIjoxNTY4Mzk1MDc2fQ.yPyUS_5k8RMvfowGylHqiCJqYwe-AOGtpBnjvqP4Na8",
			wantEmail: "foo@bar.com",
			wantErr:   false,
		},
		{
			name:      "bad token",
			idToken:   "badtoken",
			wantEmail: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClaims(tt.idToken)
			if !tt.wantErr {
				assert.NoError(t, err, "expected no error but received %s", err)
			}

			if tt.wantEmail != "" {
				assert.Equal(t, tt.wantEmail, cast.ToString(got["email"]))
			}
		})
	}
}

func TestGenerateClientSecret(t *testing.T) {
	testGoodKey := `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg+94fs23vSrhBIXNz
OdeRb7+FJkIsVrnTSf7eIYKdf4mgCgYIKoZIzj0DAQehRANCAATyBS3eRgOJ53OQ
LFhGSJw4aiqju7muVwoIWFxCcFJasRwyGcbs0C7vt3xKV/DRJvID4UljaI53wETq
RxlkNCeV
-----END PRIVATE KEY-----` // A revoked key that can be used for testing

	testWrongKey := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCjcGqTkOq0CR3rTx0ZSQSIdTrDrFAYl29611xN8aVgMQIWtDB/
lD0W5TpKPuU9iaiG/sSn/VYt6EzN7Sr332jj7cyl2WrrHI6ujRswNy4HojMuqtfa
b5FFDpRmCuvl35fge18OvoQTJELhhJ1EvJ5KUeZiuJ3u3YyMnxxXzLuKbQIDAQAB
AoGAPrNDz7TKtaLBvaIuMaMXgBopHyQd3jFKbT/tg2Fu5kYm3PrnmCoQfZYXFKCo
ZUFIS/G1FBVWWGpD/MQ9tbYZkKpwuH+t2rGndMnLXiTC296/s9uix7gsjnT4Naci
5N6EN9pVUBwQmGrYUTHFc58ThtelSiPARX7LSU2ibtJSv8ECQQDWBRrrAYmbCUN7
ra0DFT6SppaDtvvuKtb+mUeKbg0B8U4y4wCIK5GH8EyQSwUWcXnNBO05rlUPbifs
DLv/u82lAkEAw39sTJ0KmJJyaChqvqAJ8guulKlgucQJ0Et9ppZyet9iVwNKX/aW
9UlwGBMQdafQ36nd1QMEA8AbAw4D+hw/KQJBANJbHDUGQtk2hrSmZNoV5HXB9Uiq
7v4N71k5ER8XwgM5yVGs2tX8dMM3RhnBEtQXXs9LW1uJZSOQcv7JGXNnhN0CQBZe
nzrJAWxh3XtznHtBfsHWelyCYRIAj4rpCHCmaGUM6IjCVKFUawOYKp5mmAyObkUZ
f8ue87emJLEdynC1CLkCQHduNjP1hemAGWrd6v8BHhE3kKtcK6KHsPvJR5dOfzbd
HAqVePERhISfN6cwZt5p8B3/JUwSR8el66DF7Jm57BM=
-----END RSA PRIVATE KEY-----` // Wrong format - this is PKCS1

	tests := []struct {
		name       string
		signingKey string
		wantSecret bool
		wantErr    bool
	}{
		{
			name:       "bad key",
			signingKey: "bad_key",
			wantSecret: false,
			wantErr:    true,
		},
		{
			name:       "bad key wrong format",
			signingKey: testWrongKey,
			wantSecret: false,
			wantErr:    true,
		},
		{
			name:       "good key",
			signingKey: testGoodKey,
			wantSecret: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewAppleSign(tt.signingKey, "1234567890", "com.example.app", "0987654321")
			err := obj.GenerateClientSecret()
			got := obj.Secret
			if !tt.wantErr {
				assert.NoError(t, err, "expected no error but got %s", err)
			}
			if tt.wantSecret {
				assert.NotEmpty(t, got, "wanted a secret string returned but got none")

				decoded, err := GetClaims(got)
				assert.NoError(t, err, "error while decoding the secret")

				r := cast.ToString(decoded["iss"])
				b := r == ""
				assert.True(t, b, "invalid issuer")
				assert.Equal(t, "1234567890", r)

				r = cast.ToString(decoded["sub"])
				b = r == ""
				assert.True(t, b, "invalid subject")
				assert.Equal(t, "com.example.app", r)
			}
		})
	}
}
