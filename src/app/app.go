package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-quantus-service/engine/route"
	_ "go-quantus-service/src/config"
	"go-quantus-service/src/pkg"
	"log"
	"net"
	"net/http"
	"os"
)

func ParseCIDRs(cidrs []string) ([]*net.IPNet, error) {
	var ipNets []*net.IPNet
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR %q: %w", cidr, err)
		}
		ipNets = append(ipNets, ipnet)
	}
	return ipNets, nil
}

func IPWhitelistGinMiddleware(allowedNets []*net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipStr := c.ClientIP()
		clientIP := net.ParseIP(ipStr)
		if clientIP == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid IP",
			})
			return
		}

		for _, ipNet := range allowedNets {
			if ipNet.Contains(clientIP) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: IP not allowed " + ipStr,
		})
	}
}

type Headers struct {
	Key   string
	Value string
}

func GinHeaderMiddleware(headers ...Headers) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, h := range headers {
			if h.Value != "" {
				c.Writer.Header().Set(h.Key, h.Value)
			}
		}
		c.Next()
	}
}

func Start() {
	if viper.GetString("SERVICE_MODE") == "development" {
		gin.SetMode(gin.DebugMode)
		log.Println("Running in development mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Logging
	router.Use(gin.LoggerWithWriter(os.Stdout))

	// IP Whitelist Middleware
	ipNets, err := ParseCIDRs(pkg.WhitelistCIDRs)
	if err != nil {
		log.Printf("Failed to parse CIDRs: %v\n", err)
	}
	router.Use(IPWhitelistGinMiddleware(ipNets))

	// Default Headers Middleware
	defaultHeaders := []Headers{
		{Key: pkg.H_LANG, Value: "en"},
		{Key: pkg.H_CURRENCY, Value: pkg.CURRENCY_USD},
		{Key: pkg.H_USERID, Value: ""},
		{Key: pkg.H_VISIBILITY, Value: "0"},
		{Key: pkg.H_XTIMEZONE, Value: viper.GetString("TIMEZONE")},
	}
	router.Use(GinHeaderMiddleware(defaultHeaders...))

	// CORS
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "X-Client", "X-License-ID", "X-License", "X-Device-ID", "X-Session-ID", "token", "X-Api-Key"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// Register Routes
	routeDefine(router)
	// Start Server

	port := ":" + viper.GetString("SERVICE_PORT")
	fmt.Println("Running server on", port)
	if err := router.Run(port); err != nil {
		panic(err)
	}
}

func routeDefine(router *gin.Engine) {
	login, _ := InitializeStartupController()
	var intialController = route.InitialController{
		UserController: login,
	}
	intialController.RegisterGinRoutes(router)
}
