package job

import (
	"github.com/robfig/cron/v3"

	"github.com/horlakz/wallet-sync.api/internal/config"
	"github.com/horlakz/wallet-sync.api/lib/database"
	core_repository "github.com/horlakz/wallet-sync.api/repository/core"
	"github.com/horlakz/wallet-sync.api/service"
)

type CronService struct {
	cron                  *cron.Cron
	logger                *config.Logger
	reconciliationService service.ReconciliationService
}

type CronServiceInterface interface {
	Start()
}

func NewCronService() CronServiceInterface {
	env := config.GetEnv()
	db := database.StartDatabaseClient(env)

	ledgerEntryRepo := core_repository.NewLedgerEntryRepository(db)
	transactionRepo := core_repository.NewTransactionRepository(db)
	accountRepo := core_repository.NewAccountRepository(db)
	reconciliationLogRepo := core_repository.NewReconciliationLogRepository(db)

	reconciliationService := service.NewReconciliationService(
		ledgerEntryRepo,
		transactionRepo,
		accountRepo,
		reconciliationLogRepo,
	)

	return &CronService{
		cron:                  cron.New(cron.WithSeconds()),
		logger:                config.NewLogger(),
		reconciliationService: reconciliationService,
	}
}

func (c *CronService) Start() {
	// Run every hour
	c.cron.AddFunc("@every 1h", func() {
		if err := c.reconciliationService.ReconcileTransactions(); err != nil {
			c.logger.Log().Errorf("Failed to reconcile transactions: %v", err)
		} else {
			c.logger.Log().Info("Transaction reconciliation completed successfully")
		}
	})

	c.logger.Log().Info("Cron service started")
	c.cron.Start()
}
