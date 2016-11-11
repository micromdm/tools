package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_signRequest_HTTP(t *testing.T) {
	tests := []struct {
		name    string
		fields  *signRequest
		wantErr bool
	}{
		{
			name:   "default",
			fields: newSignRequest("groob@acme.co", []byte("fakecsr"), []byte("fakecert")),
		},
	}
	for _, tt := range tests {
		sign := &signRequest{
			CSR:     tt.fields.CSR,
			Email:   tt.fields.Email,
			Key:     tt.fields.Key,
			Encrypt: tt.fields.Encrypt,
		}
		got, err := sign.HTTPRequest()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. signRequest.HTTP() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}

		if have, want := got.Header.Get("Content-Type"), "application/json"; have != want {
			t.Errorf("have %q, want %q\n", have, want)
		}
		if have, want := got.Method, "POST"; have != want {
			t.Errorf("have %q, want %q\n", have, want)
		}

		var have signRequest
		if err := json.NewDecoder(got.Body).Decode(&have); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(have, *sign) {
			t.Errorf("have %#v, want %#v", have, *sign)
		}
	}
}
