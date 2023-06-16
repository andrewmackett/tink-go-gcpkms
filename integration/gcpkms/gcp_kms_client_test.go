// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

package gcpkms_test

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/core/registry"
	"github.com/tink-crypto/tink-go-gcpkms/integration/gcpkms"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

func Example() {
	const keyURI = "gcp-kms://......"
	ctx := context.Background()
	gcpclient, err := gcpkms.NewClientWithOptions(ctx, keyURI, option.WithCredentialsFile("/mysecurestorage/credentials.json"))
	if err != nil {
		log.Fatal(err)
	}
	registry.RegisterKMSClient(gcpclient)

	dek := aead.AES128CTRHMACSHA256KeyTemplate()
	template, err := aead.CreateKMSEnvelopeAEADKeyTemplate(keyURI, dek)
	if err != nil {
		log.Fatal(err)
	}
	handle, err := keyset.NewHandle(template)
	if err != nil {
		log.Fatal(err)
	}
	a, err := aead.New(handle)
	if err != nil {
		log.Fatal(err)
	}

	ct, err := a.Encrypt([]byte("this data needs to be encrypted"), []byte("this data needs to be authenticated, but not encrypted"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.Decrypt(ct, []byte("this data needs to be authenticated, but not encrypted"))
	if err != nil {
		log.Fatal(err)
	}
}
