package app

import (
	"context"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"sync"
	"time"
)

type ModeServiceImpl struct {
	currentMode     string
	mu              sync.RWMutex
	exchangeService domain.ExchangeClient
	dataProcessor   domain.DataProcessingService
	cache           domain.RedisClient
	repository      domain.Repository
	config          *config.Config
	cancelFunc      context.CancelFunc
	messagesChan    <-chan domain.PriceTick
}

func NewModeService(
	exchangeService domain.ExchangeClient,
	dataProcessor domain.DataProcessingService,
	cache domain.RedisClient,
	repository domain.Repository,
	config *config.Config,
) domain.SystemService {
	return &ModeServiceImpl{
		exchangeService: exchangeService,
		dataProcessor:   dataProcessor,
		cache:           cache,
		repository:      repository,
		config:          config,
		currentMode:     "none",
	}
}

func (m *ModeServiceImpl) SwitchToTestMode(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentMode == "test" {
		return nil
	}

	slog.Info("Starting switch to test mode", "currentMode", m.currentMode)

	m.stopCurrentMode()

	modeCtx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	m.currentMode = "test"
	
	messages := m.exchangeService.StartTestMode(modeCtx)
	m.messagesChan = messages
	m.dataProcessor.StartWorkers(modeCtx, messages)

	slog.Info("Successfully switched to test mode")
	return nil
}

func (m *ModeServiceImpl) SwitchToLiveMode(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentMode == "live" {
		return nil
	}

	slog.Info("Starting switch to live mode", "currentMode", m.currentMode)

	m.stopCurrentMode()

	modeCtx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	m.currentMode = "live"

	messages := m.exchangeService.StartLiveMode(modeCtx)
	m.messagesChan = messages
	m.dataProcessor.StartWorkers(modeCtx, messages)

	slog.Info("Successfully switched to live mode")
	return nil
}

func (m *ModeServiceImpl) stopCurrentMode() {
	if m.currentMode == "none" {
		return
	}

	slog.Info("Stopping current mode", "mode", m.currentMode)

	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}

	m.dataProcessor.StopWorkers()

	slog.Info("Current mode stopped", "mode", m.currentMode)
	m.currentMode = "none"
}

func (m *ModeServiceImpl) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	slog.Info("SystemService shutdown requested", "mode", m.currentMode)
	m.stopCurrentMode()
	return nil
}

func (m *ModeServiceImpl) GetCurrentMode(ctx context.Context) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentMode, nil
}

func (m *ModeServiceImpl) IsLiveMode() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentMode == "live"
}

func (m *ModeServiceImpl) GetSystemHealth(ctx context.Context) (domain.SystemHealth, error) {
	health := domain.SystemHealth{
		Status:     "healthy",
		Service:    "marketflow",
		Message:    "System is operational",
		Redis:      m.CheckRedisHealth(ctx),
		PostgreSQL: m.CheckPostgresHealth(ctx),
		Exchanges:  m.CheckExchangeHealth(ctx),
		Mode:       m.currentMode,
		Timestamp:  time.Now(),
	}

	if !health.Redis || !health.PostgreSQL {
		health.Status = "degraded"
		health.Message = "Some components are unavailable"
	}

	return health, nil
}

func (m *ModeServiceImpl) CheckRedisHealth(ctx context.Context) bool {
	if m.cache == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	price := m.cache.GetLatestPriceBySymbol(ctx, "BTCUSDT")
	return price > 0
}

func (m *ModeServiceImpl) CheckPostgresHealth(ctx context.Context) bool {
	if m.repository == nil {
		return false
	}

	price := m.repository.GetLowestPriceBySymbol(ctx, "BTCUSDT")
	return price > 0
}

func (m *ModeServiceImpl) CheckExchangeHealth(ctx context.Context) map[string]bool {
	health := make(map[string]bool)

	if !m.IsLiveMode() {
		for _, exchange := range m.config.Exchanges {
			health[exchange.Name] = true
		}
		return health
	}

	for _, exchange := range m.config.Exchanges {
		health[exchange.Name] = true
	}

	return health
}
