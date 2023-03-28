package service

import (
	"net/http"

	"github.com/keploy/go-sdk/integrations/khttpclient"
	"github.com/keploy/go-sdk/integrations/kmongo"
	"github.com/keploy/go-sdk/keploy"
	"go.keploy.io/server/config"
	mockPlatform "go.keploy.io/server/pkg/platform/fs"
	"go.keploy.io/server/pkg/platform/mgo"
	"go.keploy.io/server/pkg/platform/telemetry"
	"go.keploy.io/server/pkg/service/browserMock"
	"go.keploy.io/server/pkg/service/mock"
	"go.keploy.io/server/pkg/service/regression"
	"go.keploy.io/server/pkg/service/testCase"
	"go.uber.org/zap"
)

type KServices struct {
	TestcaseSrv     testCase.Service
	RegressionSrv   regression.Service
	MockSrv         mock.Service
	BrowserMockSrv  browserMock.Service
	TelemetrySrv    telemetry.Service
	TelemetryClient http.Client
}

func NewServices(ver string, conf *config.Config, logger *zap.Logger) *KServices {
	cl, err := mgo.New(conf.MongoURI)
	if err != nil {
		logger.Fatal("failed to create mgo db client", zap.Error(err))
		return nil
	}
	db := cl.Database(conf.DB)

	tdb := mgo.NewTestCase(kmongo.NewCollection(db.Collection(conf.TestCaseTable)), logger)
	rdb := mgo.NewRun(kmongo.NewCollection(db.Collection(conf.TestRunTable)), kmongo.NewCollection(db.Collection(conf.TestTable)), logger)
	mdb := mgo.NewBrowserMockDB(kmongo.NewCollection(db.Collection("test-browser-mocks")), logger)

	mockFS := mockPlatform.NewMockExportFS(keploy.GetMode() == keploy.MODE_TEST)
	testReportFS := mockPlatform.NewTestReportFS(keploy.GetMode() == keploy.MODE_TEST)
	teleFS := mockPlatform.NewTeleFS()

	browserMockSrv := browserMock.NewBrMockService(mdb, logger)
	enabled := conf.EnableTelemetry
	analyticsConfig := telemetry.NewTelemetry(mgo.NewTelemetryDB(db, conf.TelemetryTable, enabled, logger), enabled, keploy.GetMode() == keploy.MODE_OFF, conf.EnableTestExport, teleFS, logger, ver)

	client := http.Client{
		Transport: khttpclient.NewInterceptor(http.DefaultTransport),
	}

	tcSvc := testCase.New(tdb, logger, conf.EnableDeDup, analyticsConfig, client, conf.EnableTestExport, mockFS)
	regSrv := regression.New(tdb, rdb, testReportFS, analyticsConfig, client, logger, conf.EnableTestExport, mockFS)
	mockSrv := mock.NewMockService(mockFS, logger)

	return &KServices{
		TestcaseSrv:     tcSvc,
		RegressionSrv:   regSrv,
		MockSrv:         mockSrv,
		BrowserMockSrv:  browserMockSrv,
		TelemetrySrv:    analyticsConfig,
		TelemetryClient: client,
	}

}
