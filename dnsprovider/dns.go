// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googledns "google.golang.org/api/dns/v1"
)

// TODO: extend Certificate kube object to have field 'provider'

type Interface interface {
	ProviderName() string
	CreateDNSRecord(fqdn, value string, ttl int) error
	DeleteDNSRecord(fqdn string) error
}

func DNSChallengeRecord(domain, token, jwkThumbprint string) (string, string, int) {
	fqdn := fmt.Sprintf("_acme-challenge.%s.", domain)
	keyAuthorization := fmt.Sprintf("%s.%s", token, jwkThumbprint)
	keyAuthorizationShaBytes := sha256.Sum256([]byte(keyAuthorization))
	value := base64.URLEncoding.EncodeToString(keyAuthorizationShaBytes[:sha256.Size])
	value = strings.TrimRight(value, "=")
	ttl := 30
	return fqdn, value, ttl
}
