package db_gorm

import (
	"github.com/haozzzzzzzz/go-rapid-development/v2/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"runtime/debug"
	"time"
)

var (
	InstanceId          string
	ExecTimesCounterVec *prometheus.CounterVec // 执行语句统计次数
	SpentTimeSummaryVec *prometheus.SummaryVec // 耗时
)

// InitMetrics 初始化metrics对象
func InitMetrics(
	namespace string,
	instanceId string,
) (err error) {
	InstanceId = instanceId

	ExecTimesCounterVec, err = metrics.NewCounterVec(
		namespace,
		"gorm",
		"exec_times",
		"执行次数",
		[]string{"instance", "type"},
	)
	if err != nil {
		logrus.Errorf("new metrics counter vec failed. error: %s", err)
		return
	}

	SpentTimeSummaryVec, err = metrics.NewSummaryVec(
		namespace,
		"gorm",
		"spent_time",
		"耗时统计",
		[]string{"instance"},
	)
	if err != nil {
		logrus.Errorf("new metrics summary vec failed. error: %s", err)
		return
	}

	return
}

type DbMetrics struct {
	Statement *gorm.Statement
	StartTime time.Time
	EndTime   time.Time
}

func DbMetricsBefore(db *gorm.DB) {
	defer func() {
		if iRec := recover(); iRec != nil {
			logrus.Errorf("DbMetricsBefore panic. error: %s", debug.Stack())
		}
	}()

	db.InstanceSet("db_metrics", &DbMetrics{
		StartTime: time.Now(),
		Statement: db.Statement,
	})

	ExecTimesCounterVec.WithLabelValues(InstanceId, "total").Inc()
}

func DbMetricsAfter(db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{
		"func": "DbMetricsAfter",
		"sql":  db.Statement.SQL.String(),
		"vars": db.Statement.Vars,
	})

	defer func() {
		if iRec := recover(); iRec != nil {
			logrus.Errorf("DbMetricsAfter panic. error: %s", debug.Stack())
		}
	}()

	iDbMetrics, ok := db.InstanceGet("db_metrics")
	if !ok {
		logger.Errorf("db metrics after get obj not ok. %p", db.Statement)
		return
	}

	dbMetrics, ok := iDbMetrics.(*DbMetrics)
	if !ok {
		logger.Errorf("db metrics after get obj not ok. %p", db.Statement)
		return
	}

	if dbMetrics.Statement != db.Statement {
		logger.Errorf("db metrics statement not match")
		return
	}

	dbErr := db.Error
	if dbErr != nil && dbErr != gorm.ErrRecordNotFound {
		ExecTimesCounterVec.WithLabelValues(InstanceId, "total").Inc()
	}

	dbMetrics.EndTime = time.Now()
	duration := dbMetrics.EndTime.Sub(dbMetrics.StartTime)
	milSeconds := duration.Milliseconds()
	SpentTimeSummaryVec.WithLabelValues(InstanceId).Observe(float64(milSeconds))

	if milSeconds > 1000 {
		ExecTimesCounterVec.WithLabelValues(InstanceId, "slow").Inc()
	}
}

func SetDbMetrics(
	dbClient *gorm.DB,
) (err error) {
	err = dbClient.Callback().Create().Before("*").Register("create_metrics_before", DbMetricsBefore)
	if err != nil {
		return
	}

	err = dbClient.Callback().Create().After("*").Register("create_metrics_after", DbMetricsAfter)
	if err != nil {
		return
	}

	err = dbClient.Callback().Query().Before("*").Register("query_metrics_before", DbMetricsBefore)
	if err != nil {
		return
	}

	err = dbClient.Callback().Query().After("*").Register("query_metrics_after", DbMetricsAfter)
	if err != nil {
		return
	}

	err = dbClient.Callback().Update().Before("*").Register("update_metrics_before", DbMetricsBefore)
	if err != nil {
		return
	}

	err = dbClient.Callback().Update().After("*").Register("update_metrics_after", DbMetricsAfter)
	if err != nil {
		return
	}

	err = dbClient.Callback().Delete().Before("*").Register("delete_metrics_before", DbMetricsBefore)
	if err != nil {
		return
	}

	err = dbClient.Callback().Delete().After("*").Register("delete_metrics_after", DbMetricsAfter)
	if err != nil {
		return
	}

	err = dbClient.Callback().Row().Before("*").Register("row_metrics_before", DbMetricsBefore)
	if err != nil {
		return
	}

	err = dbClient.Callback().Row().After("*").Register("row_metrics_after", DbMetricsAfter)
	if err != nil {
		return
	}

	err = dbClient.Callback().Raw().Before("*").Register("raw_metrics_before", DbMetricsBefore)
	if err != nil {
		return
	}

	err = dbClient.Callback().Raw().After("*").Register("raw_metrics_after", DbMetricsAfter)
	if err != nil {
		return
	}
	return
}
