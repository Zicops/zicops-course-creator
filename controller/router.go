package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/graph"
	"github.com/zicops/zicops-course-creator/graph/generated"
	"github.com/zicops/zicops-course-creator/handlers"
	"github.com/zicops/zicops-course-creator/lib/jwt"
)

// CCRouter ... the router for the controller
func CCRouter() (*gin.Engine, error) {
	restRouter := gin.Default()
	// configure cors as needed for FE/BE interactions: For now defaults
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	restRouter.MaxMultipartMemory = 200 << 20 // 200 MiB
	restRouter.Use(cors.New(configCors))
	restRouter.Use(func(c *gin.Context) {
		currentRequest := c.Request
		incomingToken := jwt.GetToken(currentRequest)
		claimsFromToken, _ := jwt.GetClaims(incomingToken)
		c.Set("zclaims", claimsFromToken)
	})
	restRouter.GET("/healthz", HealthCheckHandler)
	// create group for restRouter
	version1 := restRouter.Group("/api/v1")
	version1.POST("/query", graphqlHandler())
	version1.GET("/playql", playgroundHandler())
	// post zip file
	version1.POST("/uploadStaticZip", UploadStaticZipHandler)
	return restRouter, nil
}

func UploadStaticZipHandler(c *gin.Context) {
	ctxValue := c.Value("zclaims").(map[string]interface{})
	lspIdInt := ctxValue["tenant"]
	lspID := "d8685567-cdae-4ee0-a80e-c187848a760e"
	if lspIdInt != nil && lspIdInt.(string) != "" {
		lspID = lspIdInt.(string)
	}
	ctxValue["lsp_id"] = lspID
	// set ctxValue to request context
	request := c.Request
	ctx := context.WithValue(request.Context(), "zclaims", ctxValue)
	request = request.WithContext(ctx)
	c.Request = request
	res, err := handlers.UploadStaticZipHandler(c)
	if err != nil {
		ResponseError(c.Writer, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
func HealthCheckHandler(c *gin.Context) {
	log.Debugf("HealthCheckHandler Method --> %s", c.Request.Method)

	switch c.Request.Method {
	case http.MethodGet:
		GetHealthStatus(c.Writer)
	default:
		err := errors.New("method not supported")
		ResponseError(c.Writer, http.StatusBadRequest, err)
	}
}

//GetHealthStatus ...
func GetHealthStatus(w http.ResponseWriter) {
	healthStatus := "Super Dentist backend service is healthy"
	response, _ := json.Marshal(healthStatus)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(response); err != nil {
		log.Errorf("GetHealthStatus ... unable to write JSON response: %v", err)
	}
}

// ResponseError ... essentially a single point of sending some error to route back
func ResponseError(w http.ResponseWriter, httpStatusCode int, err error) {
	log.Errorf("Response error %s", err.Error())
	response, _ := json.Marshal(err)
	w.Header().Add("Status", strconv.Itoa(httpStatusCode)+" "+err.Error())
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(httpStatusCode)

	if _, err := w.Write(response); err != nil {
		log.Errorf("ResponseError ... unable to write JSON response: %v", err)
	}
}

func graphqlHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// NewExecutableSchema and Config are in the generated.go file
		// Resolver is in the resolver.go file
		srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
		srv.AddTransport(transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
		})
		srv.AddTransport(transport.Options{})
		srv.AddTransport(transport.GET{})
		srv.AddTransport(transport.POST{})
		cache := lru.New(1000)
		srv.SetQueryCache(cache)
		srv.Use(extension.Introspection{})
		srv.Use(extension.AutomaticPersistedQuery{
			Cache: lru.New(100),
		})
		var mb int64 = 1 << 20
		srv.AddTransport(transport.MultipartForm{
			MaxMemory:     250 * mb,
			MaxUploadSize: 250 * mb,
		})
		lspIdInt := c.Request.Header.Get("tenant")
		ctxValue := c.Value("zclaims").(map[string]interface{})
		lspID := "d8685567-cdae-4ee0-a80e-c187848a760e"
		if lspIdInt != "" {
			lspID = lspIdInt
		}
		ctxValue["lsp_id"] = lspID
		// set ctxValue to request context
		request := c.Request
		requestWithValue := request.WithContext(context.WithValue(request.Context(), "zclaims", ctxValue))
		srv.ServeHTTP(c.Writer, requestWithValue)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("CourseCreator", "/api/v1/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
