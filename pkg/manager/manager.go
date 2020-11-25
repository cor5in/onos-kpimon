// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package manager

import (
	"github.com/onosproject/onos-kpimon/pkg/southbound/ricapie2"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-ric-sdk-go/pkg/gnmi"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	CAPath      string
	KeyPath     string
	CertPath    string
	E2tEndpoint string
	GRPCPort    int
}

// NewManager creates a new manager
func NewManager(config Config) *Manager {
	log.Info("Creating Manager")

	e2tSession, err := ricapie2.NewSession(config.E2tEndpoint)

	if err != nil {
		log.Error(err)
		return nil
	}

	return &Manager{
		Config:     config,
		E2tSession: e2tSession,
	}
}

// Manager is a manager for the KPIMON service
type Manager struct {
	Config     Config
	E2tSession *ricapie2.E2Session
}

// Run starts the manager and the associated services
func (m *Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {

	err := m.startNorthboundServer()
	if err != nil {
		return err
	}

	go m.E2tSession.Run()
	return nil
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}

func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(northbound.NewServerCfg(
		m.Config.CAPath,
		m.Config.KeyPath,
		m.Config.CertPath,
		int16(m.Config.GRPCPort),
		true,
		northbound.SecurityConfig{}))

	// TODO add services including gnmi service
	s.AddService(gnmi.NewService(gnmi.ModelInfo{
		ModelType: gnmi.RIC,
		Version:   "1.0.0",
	}))

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}