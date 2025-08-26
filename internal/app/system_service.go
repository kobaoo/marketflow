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
	wg              sync.WaitGroup
	messagesChan    <-chan domain.PriceTick // храним текущий канал
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
		currentMode:     "none", // начальное состояние
	}
}

func (m *ModeServiceImpl) SwitchToTestMode(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.currentMode == "test" {
		return nil
	}
	
	slog.Info("Starting switch to test mode", "currentMode", m.currentMode)
	
	// 1. Останавливаем текущий режим
	m.stopCurrentMode()
	
	// 2. Создаем НОВЫЙ контекст (не из переданного ctx, который может быть отменен)
	modeCtx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel
	
	// 3. Устанавливаем новый режим
	m.currentMode = "test"
	
	// 4. Запускаем новый режим
	messages := m.exchangeService.StartTestMode(modeCtx)
	m.messagesChan = messages
	m.startWorkers(modeCtx, messages)
	
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
	
	// 1. Останавливаем текущий режим
	m.stopCurrentMode()
	
	// 2. Создаем НОВЫЙ контекст
	modeCtx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel
	
	// 3. Устанавливаем новый режим
	m.currentMode = "live"
	
	// 4. Запускаем новый режим
	messages := m.exchangeService.StartLiveMode(modeCtx)
	m.messagesChan = messages
	m.startWorkers(modeCtx, messages)
	
	slog.Info("Successfully switched to live mode")
	return nil
}

func (m *ModeServiceImpl) stopCurrentMode() {
	if m.currentMode == "none" {
		return
	}
	
	slog.Info("Stopping current mode", "mode", m.currentMode)
	
	// 1. Отменяем контекст режима
	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}
	
	// 2. Останавливаем exchange service
	if exchangeClient, ok := m.exchangeService.(interface{ Stop() }); ok {
		exchangeClient.Stop()
	}
	
	// 3. Останавливаем data processor
	m.dataProcessor.StopWorkers()
	
	// 4. Ждем завершения
	m.waitForShutdown()
	
	slog.Info("Current mode stopped", "mode", m.currentMode)
	m.currentMode = "none"
}

func (m *ModeServiceImpl) waitForShutdown() {
	done := make(chan struct{})
	go func() {
		m.wg.Wait() // Ждем завершения воркеров ModeService
		close(done)
	}()
	
	select {
	case <-done:
		slog.Info("All workers stopped gracefully")
	case <-time.After(5 * time.Second):
		slog.Warn("Workers stop timeout, forcing shutdown")
	}
}

func (m *ModeServiceImpl) startWorkers(ctx context.Context, messages <-chan domain.PriceTick) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.dataProcessor.StartWorkers(ctx, messages)
	}()
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