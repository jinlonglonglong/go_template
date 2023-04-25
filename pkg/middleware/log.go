package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io/ioutil"
	"time"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 请求 Body
		body, err := c.GetRawData()
		if err != nil {
			panic(err.Error())
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 关键点

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求 IP
		clientIP := c.ClientIP()

		// 日志格式
		glog.Infof("| %3d | %13v | %15s | %s | %s | %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			string(body),
		)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
