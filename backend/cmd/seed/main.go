package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/diogenes/costforensics/backend/config"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/system"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	tenantapp "github.com/diogenes/costforensics/backend/internal/application/tenant"
	"github.com/diogenes/costforensics/backend/internal/domain/anomaly"
	"github.com/diogenes/costforensics/backend/internal/domain/budget"
	"github.com/diogenes/costforensics/backend/internal/domain/cloudaccount"
	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/costreport"
	"github.com/diogenes/costforensics/backend/internal/domain/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/exchangerate"
	"github.com/diogenes/costforensics/backend/internal/domain/forensicevent"
	"github.com/diogenes/costforensics/backend/internal/domain/importjob"
	"github.com/diogenes/costforensics/backend/internal/domain/integration"
	"github.com/diogenes/costforensics/backend/internal/domain/ledger"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/payment"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/diogenes/costforensics/backend/internal/platform/database"
	"github.com/diogenes/costforensics/backend/internal/platform/logger"
	"github.com/google/uuid"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.Setup(cfg.Log)
	log.Info("Starting seed command")

	// Connect to system database
	systemDB, err := database.NewSystemDB(cfg.Database, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to system database")
	}

	// Run system migrations
	migrator := database.NewMigrator(log)
	if err := migrator.MigrateSystem(systemDB); err != nil {
		log.WithError(err).Fatal("Failed to run system migrations")
	}

	// Tenant DB manager + provisioner
	dbManager := database.NewTenantDBManager(cfg.Database, log)
	defer dbManager.Close()
	provisioner := database.NewProvisioner(systemDB, dbManager, migrator, log)

	// Create demo tenant
	tenantRepo := system.NewTenantRepo(systemDB)
	tenantService := tenantapp.NewService(tenantRepo, provisioner, log)

	ctx := context.Background()
	t, err := tenantService.Create(ctx, "Demo Company", "demo")
	if err != nil {
		log.WithError(err).Fatal("Failed to create demo tenant (may already exist)")
	}
	log.WithField("tenant_id", t.ID()).Info("Demo tenant created")

	// Connect to tenant DB
	tenantDB, err := dbManager.GetDB(t.DBName())
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to tenant database")
	}

	// Initialize all repos
	cloudAccountRepo := tenantrepo.NewCloudAccountRepo(tenantDB)
	costRecordRepo := tenantrepo.NewCostRecordRepo(tenantDB)
	budgetRepo := tenantrepo.NewBudgetRepo(tenantDB)
	anomalyRepo := tenantrepo.NewAnomalyRepo(tenantDB)
	costReportRepo := tenantrepo.NewCostReportRepo(tenantDB)
	productRepo := tenantrepo.NewProductRepo(tenantDB)
	sellerRepo := tenantrepo.NewSellerRepo(tenantDB)
	customerRepo := tenantrepo.NewCustomerRepo(tenantDB)
	orderRepo := tenantrepo.NewOrderRepo(tenantDB)
	ledgerRepo := tenantrepo.NewLedgerRepo(tenantDB)
	forensicRepo := tenantrepo.NewForensicEventRepo(tenantDB)
	promoRuleRepo := tenantrepo.NewPromotionRuleRepo(tenantDB)
	discountAppRepo := tenantrepo.NewDiscountApplicationRepo(tenantDB)
	paymentRepo := tenantrepo.NewPaymentRepo(tenantDB)
	exchangeRateRepo := tenantrepo.NewExchangeRateRepo(tenantDB)
	integrationRepo := tenantrepo.NewIntegrationRepo(tenantDB)
	importJobRepo := tenantrepo.NewImportJobRepo(tenantDB)

	// --- Reference dates ---
	now := time.Now().UTC()
	month1 := time.Date(now.Year(), now.Month()-2, 1, 0, 0, 0, 0, time.UTC)
	month2 := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	month3 := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	q1Start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	q1End := time.Date(now.Year(), 4, 1, 0, 0, 0, 0, time.UTC)

	// =====================================================================
	// 1. Cloud Accounts (3)
	// =====================================================================
	log.Info("Seeding cloud accounts...")
	awsAcct := must(cloudaccount.NewCloudAccount(cloudaccount.ProviderAWS, "AWS Production", "123456789012"))
	gcpAcct := must(cloudaccount.NewCloudAccount(cloudaccount.ProviderGCP, "GCP Staging", "my-staging-project"))
	azureAcct := must(cloudaccount.NewCloudAccount(cloudaccount.ProviderAzure, "Azure Dev", "sub-abc-123"))

	for _, a := range []*cloudaccount.CloudAccount{awsAcct, gcpAcct, azureAcct} {
		mustDo(cloudAccountRepo.Create(ctx, a))
	}

	// =====================================================================
	// 2. Cost Records (~60)
	// =====================================================================
	log.Info("Seeding cost records...")
	type costSpec struct {
		accountID uuid.UUID
		service   string
		amounts   [3]int64 // cents per month
	}
	costSpecs := []costSpec{
		{awsAcct.ID(), "EC2", [3]int64{1200000, 1350000, 1500000}},
		{awsAcct.ID(), "S3", [3]int64{80000, 95000, 110000}},
		{awsAcct.ID(), "Lambda", [3]int64{45000, 52000, 180000}},
		{awsAcct.ID(), "RDS", [3]int64{320000, 340000, 360000}},
		{awsAcct.ID(), "CloudFront", [3]int64{60000, 72000, 85000}},
		{awsAcct.ID(), "DynamoDB", [3]int64{28000, 31000, 35000}},
		{gcpAcct.ID(), "Compute Engine", [3]int64{800000, 850000, 920000}},
		{gcpAcct.ID(), "Cloud SQL", [3]int64{190000, 210000, 225000}},
		{gcpAcct.ID(), "BigQuery", [3]int64{150000, 170000, 420000}},
		{gcpAcct.ID(), "Cloud Storage", [3]int64{35000, 42000, 48000}},
		{gcpAcct.ID(), "Cloud Functions", [3]int64{22000, 25000, 28000}},
		{gcpAcct.ID(), "Cloud CDN", [3]int64{15000, 18000, 21000}},
		{azureAcct.ID(), "Virtual Machines", [3]int64{450000, 475000, 510000}},
		{azureAcct.ID(), "Blob Storage", [3]int64{42000, 48000, 55000}},
		{azureAcct.ID(), "Azure SQL", [3]int64{120000, 125000, 135000}},
		{azureAcct.ID(), "Functions", [3]int64{18000, 21000, 24000}},
		{azureAcct.ID(), "Azure CDN", [3]int64{12000, 14000, 16000}},
		{awsAcct.ID(), "EBS", [3]int64{55000, 58000, 62000}},
		{gcpAcct.ID(), "Cloud Spanner", [3]int64{250000, 265000, 280000}},
		{azureAcct.ID(), "Cosmos DB", [3]int64{95000, 102000, 110000}},
	}

	months := []time.Time{month1, month2, month3}
	for _, spec := range costSpecs {
		for mi, m := range months {
			cr := must(costrecord.NewCostRecord(
				spec.accountID, spec.service,
				shared.MustNewMoney(spec.amounts[mi], "USD"),
				m.AddDate(0, 0, 15), // mid-month
			))
			mustDo(costRecordRepo.Create(ctx, cr))
		}
	}

	// =====================================================================
	// 3. Budgets (3)
	// =====================================================================
	log.Info("Seeding budgets...")
	b1 := must(budget.NewBudget("Q1 AWS Budget", shared.MustNewMoney(5000000, "USD"), q1Start, q1End, 80))
	b2 := must(budget.NewBudget("Q1 GCP Budget", shared.MustNewMoney(3000000, "USD"), q1Start, q1End, 80))
	b3 := must(budget.NewBudget("Q1 Total Cloud Budget", shared.MustNewMoney(10000000, "USD"), q1Start, q1End, 80))

	for _, b := range []*budget.Budget{b1, b2, b3} {
		mustDo(budgetRepo.Create(ctx, b))
	}

	// Record some spending
	b1.RecordSpend(shared.MustNewMoney(3800000, "USD"))
	mustDo(budgetRepo.Update(ctx, b1))
	b2.RecordSpend(shared.MustNewMoney(2400000, "USD"))
	mustDo(budgetRepo.Update(ctx, b2))
	b3.RecordSpend(shared.MustNewMoney(7500000, "USD"))
	mustDo(budgetRepo.Update(ctx, b3))

	// =====================================================================
	// 4. Anomalies (5)
	// =====================================================================
	log.Info("Seeding anomalies...")
	anomalies := []*anomaly.CostAnomaly{
		must(anomaly.NewCostAnomaly(awsAcct.ID(), "EC2",
			shared.MustNewMoney(1350000, "USD"), shared.MustNewMoney(1500000, "USD"), 11.1)),
		must(anomaly.NewCostAnomaly(awsAcct.ID(), "Lambda",
			shared.MustNewMoney(52000, "USD"), shared.MustNewMoney(180000, "USD"), 246.2)),
		must(anomaly.NewCostAnomaly(gcpAcct.ID(), "BigQuery",
			shared.MustNewMoney(170000, "USD"), shared.MustNewMoney(420000, "USD"), 147.1)),
		must(anomaly.NewCostAnomaly(gcpAcct.ID(), "Compute Engine",
			shared.MustNewMoney(850000, "USD"), shared.MustNewMoney(920000, "USD"), 8.2)),
		must(anomaly.NewCostAnomaly(azureAcct.ID(), "Virtual Machines",
			shared.MustNewMoney(475000, "USD"), shared.MustNewMoney(510000, "USD"), 7.4)),
	}
	for _, a := range anomalies {
		mustDo(anomalyRepo.Create(ctx, a))
	}

	// =====================================================================
	// 5. Products (12)
	// =====================================================================
	log.Info("Seeding products...")
	type prodSpec struct {
		sku, ean, upc, name, desc, cat string
		cost, price                     int64
	}
	prodSpecs := []prodSpec{
		{"ELEC-LAPTOP-001", "5901234123457", "012345678905", "ProBook Laptop 15\"", "15-inch business laptop, 16GB RAM, 512GB SSD", "electronics", 65000, 119900},
		{"ELEC-PHONE-001", "5901234123464", "012345678912", "SmartX Phone Pro", "6.7-inch AMOLED, 128GB", "electronics", 35000, 79900},
		{"ELEC-TABLET-001", "5901234123471", "012345678929", "TabletAir 10\"", "10-inch tablet, 64GB", "electronics", 22000, 44900},
		{"ELEC-HEADPH-001", "5901234123488", "012345678936", "SoundMax Headphones", "Wireless ANC over-ear headphones", "electronics", 8000, 24900},
		{"ELEC-CABLE-001", "5901234123495", "012345678943", "USB-C Cable 2m", "Braided USB-C to USB-C cable", "accessories", 300, 1499},
		{"ELEC-CASE-001", "5901234123501", "012345678950", "Universal Phone Case", "Slim protective phone case", "accessories", 500, 1999},
		{"ELEC-CHRGR-001", "5901234123518", "012345678967", "65W GaN Charger", "Compact 65W USB-C charger", "accessories", 1500, 4499},
		{"ELEC-KEYBD-001", "5901234123525", "012345678974", "MechType Keyboard", "Mechanical keyboard, RGB, hot-swap", "peripherals", 4000, 8999},
		{"ELEC-MOUSE-001", "5901234123532", "012345678981", "ErgoClick Mouse", "Ergonomic wireless mouse", "peripherals", 1500, 4999},
		{"ELEC-MONTR-001", "5901234123549", "012345678998", "UltraView 27\" Monitor", "27-inch 4K IPS monitor", "peripherals", 25000, 54900},
		{"ELEC-WEBCM-001", "5901234123556", "012345679001", "ClearSight Webcam", "1080p webcam with autofocus", "peripherals", 2000, 5999},
		{"ELEC-SPKR-001", "5901234123563", "012345679018", "BassWave Speaker", "Portable Bluetooth speaker", "electronics", 3000, 7999},
	}

	products := make([]*product.Product, len(prodSpecs))
	for i, ps := range prodSpecs {
		cost := shared.MustNewMoney(ps.cost, "USD")
		price := shared.MustNewMoney(ps.price, "USD")
		p := must(product.NewProduct(ps.sku, ps.name, cost, price))
		p.UpdateDetails(ps.name, ps.desc, ps.cat)
		products[i] = p
		mustDo(productRepo.Create(ctx, p))
	}

	// =====================================================================
	// 6. Sellers (5)
	// =====================================================================
	log.Info("Seeding sellers...")
	type sellerSpec struct {
		extID, code, name, email string
		commission               float64
	}
	sellerSpecs := []sellerSpec{
		{"S-001", "TD100", "TechDirect", "sales@techdirect.com", 8.0},
		{"S-002", "GH200", "GadgetHub", "contact@gadgethub.com", 10.0},
		{"S-003", "DW300", "DigitalWave", "hello@digitalwave.com", 12.0},
		{"S-004", "EM400", "ElectroMart", "info@electromart.com", 5.0},
		{"S-005", "BS500", "ByteShop", "orders@byteshop.com", 15.0},
	}

	sellers := make([]*seller.Seller, len(sellerSpecs))
	for i, ss := range sellerSpecs {
		s := must(seller.NewSeller(ss.extID, ss.code, ss.name, ss.email, ss.commission))
		sellers[i] = s
		mustDo(sellerRepo.Create(ctx, s))
	}

	// =====================================================================
	// 7. Customers (10)
	// =====================================================================
	log.Info("Seeding customers...")
	type custSpec struct {
		extID, name, email, segment string
	}
	custSpecs := []custSpec{
		{"C-001", "Acme Corporation", "procurement@acme.com", "enterprise"},
		{"C-002", "GlobalTech Inc.", "buying@globaltech.com", "enterprise"},
		{"C-003", "MidStream Solutions", "orders@midstream.com", "mid-market"},
		{"C-004", "BrightPath LLC", "admin@brightpath.com", "mid-market"},
		{"C-005", "QuickStart Co.", "hello@quickstart.com", "smb"},
		{"C-006", "NexaWorks", "team@nexaworks.com", "smb"},
		{"C-007", "Pixel & Code", "info@pixelcode.com", "smb"},
		{"C-008", "Jane Smith", "jane.smith@email.com", "consumer"},
		{"C-009", "Carlos Rivera", "carlos.r@email.com", "consumer"},
		{"C-010", "Aiko Tanaka", "aiko.t@email.com", "consumer"},
	}

	customers := make([]*customer.Customer, len(custSpecs))
	for i, cs := range custSpecs {
		c := must(customer.NewCustomer(cs.extID, cs.name, cs.email, cs.segment))
		customers[i] = c
		mustDo(customerRepo.Create(ctx, c))
	}

	// =====================================================================
	// 8. Exchange Rates (6)
	// =====================================================================
	log.Info("Seeding exchange rates...")
	type fxSpec struct {
		base, quote string
		rate        float64
		date        time.Time
	}
	fxSpecs := []fxSpec{
		{"USD", "EUR", 0.92, month1},
		{"USD", "GBP", 0.79, month1},
		{"USD", "ARS", 950.0, month1},
		{"USD", "EUR", 0.91, month2},
		{"USD", "GBP", 0.78, month2},
		{"EUR", "GBP", 0.86, month2},
	}

	for _, fs := range fxSpecs {
		er := must(exchangerate.NewExchangeRate(fs.base, fs.quote, fs.rate, "manual", fs.date))
		mustDo(exchangeRateRepo.Create(ctx, er))
	}

	// =====================================================================
	// 9. Promotion Rules (4)
	// =====================================================================
	log.Info("Seeding promotion rules...")
	promoStart := q1Start
	promoEnd := q1End

	promo1 := must(discount.NewPromotionRule(
		"10% Off Orders Over $500", "percentage", 10.0,
		map[string]any{"min_subtotal_cents": 50000},
		"seller", 100.0, 0, promoStart, promoEnd,
	))
	promo2 := must(discount.NewPromotionRule(
		"$20 Flat Platform Discount", "fixed_amount", 2000,
		map[string]any{},
		"platform", 0.0, 500, promoStart, promoEnd,
	))
	promo3 := must(discount.NewPromotionRule(
		"15% Off Electronics", "percentage", 15.0,
		map[string]any{"category": "electronics"},
		"shared", 60.0, 0, promoStart, promoEnd,
	))
	promo4 := must(discount.NewPromotionRule(
		"Buy 2+ Get 10% Off", "percentage", 10.0,
		map[string]any{"min_quantity": 2},
		"seller", 100.0, 0, promoStart, promoEnd,
	))

	promoRules := []*discount.PromotionRule{promo1, promo2, promo3, promo4}
	for _, r := range promoRules {
		mustDo(promoRuleRepo.Create(ctx, r))
	}

	// =====================================================================
	// 10. Orders (20) with items
	// =====================================================================
	log.Info("Seeding orders...")

	type orderDef struct {
		extID      string
		sellerIdx  int
		custIdx    int
		date       time.Time
		itemIdxs   []int // product indexes
		itemQtys   []int
		promoIdx   int // -1 for no promo
		targetStatus string
	}

	orderDefs := []orderDef{
		{"ORD-001", 0, 0, month1.AddDate(0, 0, 3), []int{0, 4, 6}, []int{2, 5, 3}, 0, "delivered"},
		{"ORD-002", 1, 1, month1.AddDate(0, 0, 7), []int{1, 3}, []int{3, 2}, 2, "delivered"},
		{"ORD-003", 2, 2, month1.AddDate(0, 0, 12), []int{7, 8}, []int{1, 1}, -1, "delivered"},
		{"ORD-004", 3, 3, month1.AddDate(0, 0, 18), []int{2, 5, 10}, []int{1, 2, 1}, 1, "delivered"},
		{"ORD-005", 4, 4, month1.AddDate(0, 0, 22), []int{9, 11}, []int{1, 2}, 0, "delivered"},
		{"ORD-006", 0, 5, month1.AddDate(0, 0, 25), []int{4, 5, 6}, []int{10, 5, 2}, 3, "delivered"},
		{"ORD-007", 1, 6, month2.AddDate(0, 0, 2), []int{0}, []int{1}, -1, "delivered"},
		{"ORD-008", 2, 7, month2.AddDate(0, 0, 5), []int{1, 3, 4}, []int{1, 1, 3}, 2, "delivered"},
		{"ORD-009", 3, 8, month2.AddDate(0, 0, 9), []int{7, 8, 10}, []int{2, 2, 1}, 3, "shipped"},
		{"ORD-010", 4, 9, month2.AddDate(0, 0, 13), []int{2}, []int{3}, 0, "shipped"},
		{"ORD-011", 0, 0, month2.AddDate(0, 0, 17), []int{11, 4, 6}, []int{1, 4, 1}, 1, "shipped"},
		{"ORD-012", 1, 1, month2.AddDate(0, 0, 20), []int{9}, []int{1}, -1, "confirmed"},
		{"ORD-013", 2, 2, month2.AddDate(0, 0, 24), []int{0, 8}, []int{1, 1}, -1, "confirmed"},
		{"ORD-014", 3, 3, month3.AddDate(0, 0, 1), []int{1, 5, 6}, []int{2, 3, 2}, 2, "confirmed"},
		{"ORD-015", 4, 4, month3.AddDate(0, 0, 3), []int{3, 7}, []int{4, 1}, 3, "pending"},
		{"ORD-016", 0, 5, month3.AddDate(0, 0, 5), []int{10, 11}, []int{2, 1}, -1, "pending"},
		{"ORD-017", 1, 6, month3.AddDate(0, 0, 7), []int{2, 4, 5}, []int{1, 2, 2}, 1, "pending"},
		{"ORD-018", 2, 7, month1.AddDate(0, 0, 15), []int{0, 1}, []int{1, 1}, -1, "cancelled"},
		{"ORD-019", 3, 8, month2.AddDate(0, 0, 8), []int{9, 7}, []int{1, 1}, 0, "refunded"},
		{"ORD-020", 4, 9, month2.AddDate(0, 0, 15), []int{3, 11}, []int{2, 1}, -1, "delivered"},
	}

	orders := make([]*order.Order, len(orderDefs))
	for i, od := range orderDefs {
		o := must(order.NewOrder(od.extID, sellers[od.sellerIdx].ID(), customers[od.custIdx].ID(), "USD"))

		// Add items
		for j, pidx := range od.itemIdxs {
			p := products[pidx]
			mustDo(o.AddItem(p.ID(), p.SKU(), p.Name(), od.itemQtys[j], p.UnitPriceCents(), p.UnitCostCents()))
		}

		// Apply promotion discount if applicable
		if od.promoIdx >= 0 {
			rule := promoRules[od.promoIdx]
			discAmt := rule.CalculateDiscount(o.SubtotalCents())
			if !discAmt.IsZero() {
				o.ApplyDiscount(discAmt)

				// Create discount application
				da := discount.NewDiscountApplication(o.ID(), rule, discAmt)
				mustDo(discountAppRepo.Create(ctx, da))
				rule.IncrementUsage()
			}
		}

		// Set shipping and tax
		o.SetShipping(shared.MustNewMoney(999, "USD"))
		o.SetTax(o.SubtotalCents().Multiply(0.08)) // 8% tax

		// Transition to target status
		switch od.targetStatus {
		case "confirmed":
			mustDo(o.Confirm())
		case "shipped":
			mustDo(o.Confirm())
			mustDo(o.Ship())
		case "delivered":
			mustDo(o.Confirm())
			mustDo(o.Ship())
			mustDo(o.Deliver())
		case "cancelled":
			mustDo(o.Cancel())
		case "refunded":
			mustDo(o.Confirm())
			mustDo(o.Ship())
			mustDo(o.Deliver())
			mustDo(o.Refund())
		}

		orders[i] = o
		mustDo(orderRepo.Create(ctx, o))
	}

	// Update promo rules (usage counts)
	for _, r := range promoRules {
		mustDo(promoRuleRepo.Update(ctx, r))
	}

	// =====================================================================
	// 11. Ledger Entries (~80)
	// =====================================================================
	log.Info("Seeding ledger entries...")
	for _, o := range orders {
		if o.Status() == order.StatusCancelled {
			continue
		}

		oid := o.ID()

		// Revenue credit
		revenueEntry := must(ledger.NewLedgerEntry(
			&oid, "revenue", ledger.SideCredit,
			o.SubtotalCents(), "Order revenue",
			"order", o.ID(), o.OrderDate(),
		))
		mustDo(ledgerRepo.Create(ctx, revenueEntry))

		// COGS debit
		totalCOGS := shared.ZeroMoney("USD")
		for _, item := range o.Items() {
			totalCOGS, _ = totalCOGS.Add(item.COGSCents())
		}
		cogsEntry := must(ledger.NewLedgerEntry(
			&oid, "cogs", ledger.SideDebit,
			totalCOGS, "Cost of goods sold",
			"order", o.ID(), o.OrderDate(),
		))
		mustDo(ledgerRepo.Create(ctx, cogsEntry))

		// Discount debit (if any)
		if !o.DiscountCents().IsZero() {
			discEntry := must(ledger.NewLedgerEntry(
				&oid, "discounts", ledger.SideDebit,
				o.DiscountCents(), "Discount applied",
				"order", o.ID(), o.OrderDate(),
			))
			mustDo(ledgerRepo.Create(ctx, discEntry))
		}

		// Commission debit (calculate from seller)
		sellerIdx := findSellerIdx(sellers, o.SellerID())
		if sellerIdx >= 0 {
			commission := sellers[sellerIdx].CalculateCommission(o.SubtotalCents())
			if !commission.IsZero() {
				commEntry := must(ledger.NewLedgerEntry(
					&oid, "commission", ledger.SideDebit,
					commission, "Seller commission",
					"order", o.ID(), o.OrderDate(),
				))
				mustDo(ledgerRepo.Create(ctx, commEntry))
			}
		}
	}

	// =====================================================================
	// 12. Payments (~30)
	// =====================================================================
	log.Info("Seeding payments...")
	for i, o := range orders {
		if o.Status() == order.StatusCancelled {
			continue
		}

		// Inbound payment from customer
		inPay := must(payment.NewPayment(
			o.ID(), payment.DirectionInbound, customers[orderDefs[i].custIdx].ID(),
			o.TotalCents(), "credit_card", fmt.Sprintf("PAY-IN-%s", o.ExternalID()),
		))
		if o.Status() == order.StatusDelivered || o.Status() == order.StatusShipped || o.Status() == order.StatusRefunded {
			inPay.MarkProcessed()
		}
		if o.Status() == order.StatusRefunded {
			inPay.MarkRefunded()
		}
		mustDo(paymentRepo.Create(ctx, inPay))

		// Outbound payment to seller (only for delivered/refunded)
		if o.Status() == order.StatusDelivered || o.Status() == order.StatusRefunded {
			sellerIdx := findSellerIdx(sellers, o.SellerID())
			commission := sellers[sellerIdx].CalculateCommission(o.SubtotalCents())
			payout, _ := o.SubtotalCents().Subtract(commission)
			if !payout.IsZero() {
				outPay := must(payment.NewPayment(
					o.ID(), payment.DirectionOutbound, sellers[sellerIdx].ID(),
					payout, "bank_transfer", fmt.Sprintf("PAY-OUT-%s", o.ExternalID()),
				))
				outPay.MarkProcessed()
				mustDo(paymentRepo.Create(ctx, outPay))
			}
		}
	}

	// =====================================================================
	// 13. Forensic Events (~40)
	// =====================================================================
	log.Info("Seeding forensic events...")
	for _, o := range orders {
		// Order created event
		must(forensicevent.NewForensicEvent(
			"order", o.ID(), "order.created",
			o.OrderDate(),
			map[string]any{
				"external_id": o.ExternalID(),
				"total_cents": o.TotalCents().Amount(),
				"status":      string(o.Status()),
			},
			"seed",
		))

		fe1 := must(forensicevent.NewForensicEvent(
			"order", o.ID(), "order.created",
			o.OrderDate(),
			map[string]any{"external_id": o.ExternalID(), "item_count": len(o.Items())},
			"seed",
		))
		mustDo(forensicRepo.Create(ctx, fe1))

		if o.Status() == order.StatusDelivered || o.Status() == order.StatusShipped ||
			o.Status() == order.StatusConfirmed || o.Status() == order.StatusRefunded {
			fe2 := must(forensicevent.NewForensicEvent(
				"order", o.ID(), "order.confirmed",
				o.OrderDate().Add(2*time.Hour),
				map[string]any{"total_cents": o.TotalCents().Amount()},
				"seed",
			))
			mustDo(forensicRepo.Create(ctx, fe2))
		}
	}

	// Price change events for a few products
	for _, pidx := range []int{0, 1, 3} {
		p := products[pidx]
		fe := must(forensicevent.NewForensicEvent(
			"product", p.ID(), "product.price_changed",
			month2.AddDate(0, 0, 10),
			map[string]any{
				"sku":             p.SKU(),
				"old_price_cents": p.UnitPriceCents().Amount() - 1000,
				"new_price_cents": p.UnitPriceCents().Amount(),
			},
			"seed",
		))
		mustDo(forensicRepo.Create(ctx, fe))
	}

	// Discount application events
	for _, o := range orders {
		if !o.DiscountCents().IsZero() {
			fe := must(forensicevent.NewForensicEvent(
				"order", o.ID(), "discount.applied",
				o.OrderDate().Add(1*time.Hour),
				map[string]any{"discount_cents": o.DiscountCents().Amount()},
				"seed",
			))
			mustDo(forensicRepo.Create(ctx, fe))
		}
	}

	// =====================================================================
	// 14. Cost Reports (2)
	// =====================================================================
	log.Info("Seeding cost reports...")
	report1 := must(costreport.NewCostReport(
		"January Cloud Costs", costreport.ReportTypeMonthly,
		month1, month2, "USD",
	))
	report1.AddServiceCost("EC2", shared.MustNewMoney(1200000, "USD"))
	report1.AddServiceCost("S3", shared.MustNewMoney(80000, "USD"))
	report1.AddServiceCost("Lambda", shared.MustNewMoney(45000, "USD"))
	report1.AddServiceCost("Compute Engine", shared.MustNewMoney(800000, "USD"))
	report1.AddServiceCost("BigQuery", shared.MustNewMoney(150000, "USD"))
	report1.AddServiceCost("Virtual Machines", shared.MustNewMoney(450000, "USD"))
	mustDo(costReportRepo.Create(ctx, report1))

	report2 := must(costreport.NewCostReport(
		"Q1 Cloud Costs Summary", costreport.ReportTypeCustom,
		q1Start, q1End, "USD",
	))
	report2.AddServiceCost("EC2", shared.MustNewMoney(4050000, "USD"))
	report2.AddServiceCost("S3", shared.MustNewMoney(285000, "USD"))
	report2.AddServiceCost("Lambda", shared.MustNewMoney(277000, "USD"))
	report2.AddServiceCost("Compute Engine", shared.MustNewMoney(2570000, "USD"))
	report2.AddServiceCost("BigQuery", shared.MustNewMoney(740000, "USD"))
	report2.AddServiceCost("Virtual Machines", shared.MustNewMoney(1435000, "USD"))
	mustDo(costReportRepo.Create(ctx, report2))

	// =====================================================================
	// 15. Integrations (3)
	// =====================================================================
	log.Info("Seeding integrations...")
	shopify := must(integration.NewIntegration("shopify", "Shopify Production",
		map[string]string{"store_url": "demo-store.myshopify.com", "api_version": "2024-01"}))
	shopify.CompleteSync()
	mustDo(integrationRepo.Create(ctx, shopify))

	woo := must(integration.NewIntegration("woocommerce", "WooCommerce Store",
		map[string]string{"site_url": "shop.democompany.com", "consumer_key": "ck_***"}))
	mustDo(integrationRepo.Create(ctx, woo))

	medusa := must(integration.NewIntegration("medusajs", "MedusaJS Headless",
		map[string]string{"api_url": "https://medusa.democompany.com"}))
	medusa.Deactivate()
	mustDo(integrationRepo.Create(ctx, medusa))

	// =====================================================================
	// 16. Import Jobs (2)
	// =====================================================================
	log.Info("Seeding import jobs...")
	job1 := must(importjob.NewImportJob("January Products CSV", importjob.SourceCSVUpload, "products", "/uploads/products_jan.csv"))
	job1.Begin()
	job1.SetTotalRows(150)
	for i := 0; i < 150; i++ {
		job1.IncrementProcessed()
	}
	job1.Complete()
	mustDo(importJobRepo.Create(ctx, job1))

	job2 := must(importjob.NewImportJob("February Orders CSV", importjob.SourceCSVUpload, "orders", "/uploads/orders_feb.csv"))
	job2.Begin()
	job2.SetTotalRows(200)
	for i := 0; i < 87; i++ {
		job2.IncrementProcessed()
	}
	job2.Fail("Invalid date format in row 88: expected YYYY-MM-DD")
	mustDo(importJobRepo.Create(ctx, job2))

	// =====================================================================
	log.Info("Seed completed successfully!")
	log.WithField("tenant_slug", "demo").Info("Use tenant slug 'demo' to access demo data")
}

// must panics on error, returning the value.
func must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("seed error: %v", err))
	}
	return v
}

// mustDo panics on error for void operations.
func mustDo(err error) {
	if err != nil {
		panic(fmt.Sprintf("seed error: %v", err))
	}
}

func findSellerIdx(sellers []*seller.Seller, id uuid.UUID) int {
	for i, s := range sellers {
		if s.ID() == id {
			return i
		}
	}
	return -1
}
