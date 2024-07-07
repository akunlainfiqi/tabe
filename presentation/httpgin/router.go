package httpgin

import (
	"fmt"
	"saas-billing/app/commands"
	"saas-billing/config"
	midtransapi "saas-billing/infrastructures/api/midtrans"
	"saas-billing/infrastructures/pgsql"
	"saas-billing/infrastructures/pubsuber"
	"saas-billing/presentation/httpgin/controller"
	"saas-billing/presentation/httpgin/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	pgclient, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			config.DB_HOST, config.DB_USER, config.DB_PASS, config.DB_NAME, config.DB_PORT, config.DB_SSL),
	}))
	if err != nil {
		panic(err)
	}

	iampgclient, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			config.DB_HOST, config.DB_USER, config.DB_PASS, "iam", config.DB_PORT, config.DB_SSL),
	}))
	if err != nil {
		panic(err)
	}

	s, c := initMidtrans()
	midtransService := midtransapi.NewMidtrans(c, s)
	publisherService := pubsuber.NewPublisher()

	v1 := r.Group("/v1")

	appQuery := pgsql.NewAppQuery(pgclient)
	iamOrganizationQuery := pgsql.NewIamUserOrganizationQuery(iampgclient)
	billQueries := pgsql.NewBillQuery(pgclient)
	productQueries := pgsql.NewProductQuery(pgclient)
	tenantQueries := pgsql.NewTenantQuery(pgclient)
	transactQueries := pgsql.NewTransactionQuery(pgclient)

	appsRepository := pgsql.NewAppsRepository(pgclient)
	billsRepository := pgsql.NewBillsRepository(pgclient)
	productRepository := pgsql.NewProductRepository(pgclient)
	priceRepository := pgsql.NewPriceRepository(pgclient)
	tenantRepository := pgsql.NewTenantRepository(pgclient)
	organizationRepository := pgsql.NewOrganizationRepository(pgclient)
	transactionRepository := pgsql.NewTransactionRepository(pgclient)

	iamOrganizationRepository := pgsql.NewIamOrganizationRepository(iampgclient)

	createBillCommand := commands.NewCreateBillsCommand(billsRepository, tenantRepository, organizationRepository, transactionRepository, priceRepository)
	createOrganizationCommand := commands.NewCreateOrganizationCommand(organizationRepository, iamOrganizationRepository)
	createProductCommand := commands.NewCreateProductCommand(productRepository, appsRepository, priceRepository)
	createTenantCommand := commands.NewCreateTenantOnboardingCommand(tenantRepository, organizationRepository, priceRepository, billsRepository, iamOrganizationRepository, midtransService)
	checkPaymentCommand := commands.NewCheckPayment(transactionRepository, billsRepository, organizationRepository, tenantRepository, priceRepository, midtransService, publisherService)
	extendTenantCommand := commands.NewExtendTenantCommand(tenantRepository, organizationRepository, priceRepository, billsRepository, midtransService)
	upgradeTenantCommand := commands.NewTenantUpgradeCommand(tenantRepository, priceRepository, billsRepository, organizationRepository, midtransService)
	downgrandeTenantCommand := commands.NewTenantDowngradeCommand(tenantRepository, priceRepository, organizationRepository, billsRepository, midtransService)

	expireBillCommand := commands.NewExpireBillsCommand(billsRepository)
	payBillCommand := commands.NewPayBillsCommand(billsRepository, tenantRepository, transactionRepository)

	appControlerr := controller.NewAppController(appQuery)
	productController := controller.NewProductController(productQueries, *createProductCommand)
	billsControlerr := controller.NewBillController(*expireBillCommand, *payBillCommand, *createBillCommand, *checkPaymentCommand, iamOrganizationQuery, billQueries)
	organizationController := controller.NewOrganizationController(*createOrganizationCommand)
	tenantController := controller.NewTenantController(*createTenantCommand, *extendTenantCommand, *upgradeTenantCommand, *downgrandeTenantCommand, tenantQueries)
	transactionController := controller.NewTransactionController(transactQueries)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	jwt := v1.Group("/jwt")
	jwt.Use(middleware.JWTMiddleware())

	jwt.GET("/ping", func(c *gin.Context) {
		user_id, ok := c.Get("user_id")
		if !ok {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": user_id,
		})
	})

	jwt.POST("/tenants", tenantController.CreateTenant)
	jwt.POST("/tenants/extend", tenantController.ExtendTenant)
	jwt.GET("/organizations/:org_id/bills", billsControlerr.GetOrganizationBills)
	jwt.GET("/organizations/:org_id/bills/:bill_id", billsControlerr.GetBillDetail)
	jwt.GET("/organizations/:org_id/bills/:bill_id/transaction", transactionController.GetByBillsID)
	jwt.GET("/organizations/:org_id/tenants", tenantController.GetByOrgID)
	jwt.POST("/organizations/:org_id/tenants", tenantController.ChangeTenantTier)
	jwt.GET("/organizations/:org_id/transactions", transactionController.GetByOrgID)

	v1.GET("/apps", appControlerr.GetAll)
	v1.GET("/products", productController.GetAll)
	v1.GET("/products/:app_id", productController.GetByAppID)
	v1.POST("/products", productController.Create)

	v1.POST("/bills/expire", billsControlerr.InternalExpire)
	v1.POST("/bills/pay", billsControlerr.InternalPay)
	v1.GET("/bills/checkall", billsControlerr.InternalCheckPayment)
	v1.POST("/bills/callback", billsControlerr.PaymentCallback)

	v1.POST("/organizations", organizationController.Create)

	return r
}

func initMidtrans() (*snap.Client, *coreapi.Client) {
	midtransKey := config.MIDTRANS_SERVER_KEY

	var s snap.Client
	var c coreapi.Client

	s.New(midtransKey, midtrans.Sandbox)
	c.New(midtransKey, midtrans.Sandbox)

	return &s, &c
}
