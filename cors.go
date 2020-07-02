package golden

import (
	"github.com/gin-contrib/cors"
	"time"
)

// https://github.com/gin-contrib/cors/blob/master/cors.go
type CORSConfig struct {
	AllowAllOrigins        bool
	AllowOrigins           []string
	AllowOriginFunc        func(origin string) bool
	AllowMethods           []string
	AllowHeaders           []string
	AllowCredentials       bool
	ExposeHeaders          []string
	MaxAge                 time.Duration
	AllowWildcard          bool
	AllowBrowserExtensions bool
	AllowWebSockets        bool
	AllowFiles             bool
}

func (golden *Golden) cors(corsConfig *CORSConfig) {
	golden.router.Use(cors.New(cors.Config{
		AllowAllOrigins:        corsConfig.AllowAllOrigins,
		AllowOrigins:           corsConfig.AllowOrigins,
		AllowOriginFunc:        corsConfig.AllowOriginFunc,
		AllowMethods:           corsConfig.AllowMethods,
		AllowHeaders:           corsConfig.AllowHeaders,
		AllowCredentials:       corsConfig.AllowCredentials,
		ExposeHeaders:          corsConfig.ExposeHeaders,
		MaxAge:                 corsConfig.MaxAge,
		AllowWildcard:          corsConfig.AllowWildcard,
		AllowBrowserExtensions: corsConfig.AllowBrowserExtensions,
		AllowWebSockets:        corsConfig.AllowWebSockets,
		AllowFiles:             corsConfig.AllowFiles,
	}))
}
