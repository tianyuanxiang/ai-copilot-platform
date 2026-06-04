package svc

import (
	"ai-copilot-platform/ai-rpc/internal/engine"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/config"
	"ai-copilot-platform/ai-rpc/internal/model"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	"go-zero-rpc/sys-rpc/pkg/orm"
	pkgsqlx "go-zero-rpc/sys-rpc/pkg/sqlx"

	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	_ "github.com/taosdata/driver-go/v3/taosWS"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	Orm *gorm.DB
	RDB *redis.Client
	ES  *elasticsearch.Client
	TD  *sql.DB

	AiDocumentModel               model.AiDocumentModel
	AiDocumentParentChunkModel    model.AiDocumentParentChunkModel
	AiDocumentChunkModel          model.AiDocumentChunkModel
	AiKnowledgeBaseModel          model.AiKnowledgeBaseModel
	AiKbDomainModel               model.AiKbDomainModel
	AiKbMemberModel               model.AiKbMemberModel
	AiToolCallLogModel            model.AiToolCallLogModel
	AiAlarmAnalysisModel          model.AiAlarmAnalysisModel
	AiHealthReportModel           model.AiHealthReportModel
	AiMaintenanceTicketDraftModel model.AiMaintenanceTicketDraftModel
	AiConversationModel           model.AiConversationModel
	AiMessageModel                model.AiMessageModel
	AiLlmCallLogModel             model.AiLlmCallLogModel

	PermRpc permclient.PermissionService
	// WindMetadataModel          model.WindMetadataModel
	WindFarmModel          model.WindFarmModel
	WindTowerModel         model.WindTowerModel
	WindDeviceModel        model.WindDeviceModel
	WindDeviceTypeModel    model.WindDeviceTypeModel
	WindStructureTypeModel model.WindStructureTypeModel
	WindCameraRecordModel  model.WindCameraRecordModel
	WindDeviceMetaModel    model.WindDeviceMetaModel
	TdengineModel          model.TdengineModel
	ThresholdResolver      model.ThresholdResolver

	EngineClient     *http.Client
	EngineCallClient *engine.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rawConn := sqlx.NewSqlConn("pgx", c.DB.DataSource)
	conn := pkgsqlx.NewTimeoutConn(rawConn, 30*time.Second)

	db := orm.NewPostgres(&orm.Config{
		DSN:         c.DB.DataSource,
		Active:      20,
		Idle:        10,
		IdleTimeout: 24 * time.Hour,
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:        c.BizRedis.Host,
		Password:    c.BizRedis.Pass,
		DB:          c.BizRedis.DB,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis连接失败: %v", err))
	}
	logx.Info("Redis连接初始化成功")

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: c.Elasticsearch.Addresses,
		Username:  c.Elasticsearch.Username,
		Password:  c.Elasticsearch.Password,
	})
	if err != nil {
		panic(fmt.Sprintf("init elasticsearch client failed: %v", err))
	}

	// 启动时做一次 Ping，尽早发现 ES 地址、账号、密码错误
	res, err := esClient.Info()
	if err != nil {
		panic(fmt.Sprintf("connect elasticsearch failed: %v", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		panic(fmt.Sprintf("connect elasticsearch failed, status: %s", res.Status()))
	}

	logx.Infof("elasticsearch connected: %s", res.Status())

	td := initTDengine(c)

	engineTimeout := time.Duration(c.Engine.TimeoutMs) * time.Millisecond
	if engineTimeout <= 0 {
		engineTimeout = 30 * time.Second
	}

	sysCli := zrpc.MustNewClient(c.SysRpc)
	permSvc := permclient.NewPermissionService(sysCli)

	return &ServiceContext{
		Config:                        c,
		Orm:                           db,
		ES:                            esClient,
		AiDocumentModel:               model.NewAiDocumentModel(conn, db),
		AiDocumentParentChunkModel:    model.NewAiDocumentParentChunkModel(conn, db),
		AiDocumentChunkModel:          model.NewAiDocumentChunkModel(db),
		AiKnowledgeBaseModel:          model.NewAiKnowledgeBaseModel(conn, db),
		AiKbDomainModel:               model.NewAiKbDomainModel(conn, db),
		AiKbMemberModel:               model.NewAiKbMemberModel(conn, db),
		AiToolCallLogModel:            model.NewAiToolCallLogModel(conn, db),
		AiAlarmAnalysisModel:          model.NewAiAlarmAnalysisModel(conn),
		AiHealthReportModel:           model.NewAiHealthReportModel(conn),
		AiMaintenanceTicketDraftModel: model.NewAiMaintenanceTicketDraftModel(conn),
		AiConversationModel:           model.NewAiConversationModel(conn, db),
		AiMessageModel:                model.NewAiMessageModel(conn, db),
		AiLlmCallLogModel:             model.NewAiLlmCallLogModel(conn, db),
		// WindMetadataModel:          model.NewWindMetadataModel(db),
		WindFarmModel:          model.NewWindFarmModel(conn, db),
		WindTowerModel:         model.NewWindTowerModel(conn, db),
		WindDeviceModel:        model.NewWindDeviceModel(conn, db),
		WindDeviceTypeModel:    model.NewWindDeviceTypeModel(conn, db),
		WindDeviceMetaModel:    model.NewWindDeviceMetaModel(conn, db),
		WindCameraRecordModel:  model.NewWindCameraRecordModel(conn),
		WindStructureTypeModel: model.NewWindStructureTypeModel(conn),
		TdengineModel:          model.NewTdengineModel(td),
		ThresholdResolver:      model.NewThresholdResolver(db),
		TD:                     td,
		RDB:                    rdb,
		EngineClient:           &http.Client{Timeout: engineTimeout},
		EngineCallClient:       engine.NewClient(c.Engine.BaseURL, &http.Client{Timeout: engineTimeout}),

		PermRpc: permSvc,
	}
}

func initTDengine(c config.Config) *sql.DB {
	link := strings.TrimSpace(c.TDengine.Link)
	if link == "" {
		logx.Info("TDengine link is empty; wind timeseries APIs will run in scaffold mode")
		return nil
	}

	db, err := sql.Open("taosWS", link)
	if err != nil {
		logx.Errorf("init TDengine connection failed: %v", err)
		return nil
	}
	if c.TDengine.MaxIdleConn > 0 {
		db.SetMaxIdleConns(c.TDengine.MaxIdleConn)
	}
	if c.TDengine.MaxOpenConn > 0 {
		db.SetMaxOpenConns(c.TDengine.MaxOpenConn)
	}
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	logx.Info("TDengine connection pool initialized")
	return db
}
