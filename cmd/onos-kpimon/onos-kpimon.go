// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	"flag"

	"github.com/onosproject/onos-kpimon/pkg/manager"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("main")

func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	e2tEndpoint := flag.String("e2tEndpoint", "onos-e2t:5150", "E2T service endpoint")

	ready := make(chan bool)

	flag.Parse()

	_, err := certs.HandleCertPaths(*caPath, *keyPath, *certPath, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting onos-kpimon")
	cfg := manager.Config{
		CAPath:      *caPath,
		KeyPath:     *keyPath,
		CertPath:    *certPath,
		E2tEndpoint: *e2tEndpoint,
		GRPCPort:    5150,
	}

	mgr := manager.NewManager(cfg)
	mgr.Run()
	<-ready
}