package app

import (
	"context"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"strings"
	"sync"
	"time"
)

type ModeServiceImpl struct {
	appCtx          context.Context 	
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
	appCtx context.Context, 
	exchangeService domain.ExchangeClient,
	dataProcessor domain.DataProcessingService,
	cache domain.RedisClient,
	repository domain.Repository,
	config *config.Config,
) domain.SystemService {
	return &ModeServiceImpl{
		appCtx: appCtx,
		exchangeService: exchangeService,
		dataProcessor:   dataProcessor,
		cache:           cache,
		repository:      repository,
		config:          config,
		currentMode:     "none",
	}
}

func (m *ModeServiceImpl) SwitchToTestMode() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentMode == "test" {
		return nil
	}

	slog.Info("Starting switch to test mode", "currentMode", m.currentMode)

	m.stopCurrentMode()

	modeCtx, cancel := context.WithCancel(m.appCtx)
	m.cancelFunc = cancel

	m.currentMode = "test"
	
	messages := m.exchangeService.StartTestMode(modeCtx)
	m.messagesChan = messages
	m.dataProcessor.StartWorkers(modeCtx, messages)

	slog.Info("Successfully switched to test mode")
	return nil
}

func (m *ModeServiceImpl) SwitchToLiveMode() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentMode == "live" {
		return nil
	}

	slog.Info("Starting switch to live mode", "currentMode", m.currentMode)

	m.stopCurrentMode()

	modeCtx, cancel := context.WithCancel(m.appCtx)
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

	issues := make([]string, 0, 4)

	if !health.Redis {
		issues = append(issues, "Redis")
	}
	if !health.PostgreSQL {
		issues = append(issues, "PostgreSQL")
	}

	downEx := make([]string, 0)
	for name, ok := range health.Exchanges {
		if !ok {
			downEx = append(downEx, name)
		}
	}
	if len(downEx) > 0 {
		issues = append(issues, "Exchanges: "+strings.Join(downEx, ", "))
	}

	if len(issues) > 0 {
		health.Status = "degraded"
		health.Message = "Some components are unavailable: " + strings.Join(issues, "; ")
	}

	return health, nil
}


func (m *ModeServiceImpl) CheckRedisHealth(ctx context.Context) bool {
	if m.cache == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	price := m.cache.GetLatestPriceBySymbol(ctx, "BTCUSDT")
	return price > 0
}

func (m *ModeServiceImpl) CheckPostgresHealth(ctx context.Context) bool {
	if m.repository == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	err := m.repository.Ping(ctx)
	if err != nil {
		return false
	}
	return true
}

func (m *ModeServiceImpl) CheckExchangeHealth(ctx context.Context) map[string]bool {
    names := make([]string, 0, len(m.config.Exchanges))
    for _, ex := range m.config.Exchanges {
        names = append(names, ex.Name)
    }

    return m.dataProcessor.ExchangesHealth(500*time.Millisecond, names)
}
