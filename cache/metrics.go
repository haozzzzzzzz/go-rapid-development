package cache

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// cache消耗时间分布
	REDIS_SPENT_TIME_SUMMARY prometheus.Summary

	// redis执行次数统计
	REDIS_EXEC_TIMES_COUNTER_VEC   *prometheus.CounterVec
	REDIS_EXEC_TIMES_COUNTER_SLOW  prometheus.Counter // 慢执行
	REDIS_EXEC_TIMES_COUNTER_TOTAL prometheus.Counter // 总执行次数
	REDIS_EXEC_TIMES_COUNTER_ERROR prometheus.Counter // 出错次数
)

func initMetrics() {

}
