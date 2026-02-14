package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/modelgate/modelgate/internal/app"
	"github.com/modelgate/modelgate/internal/config"
	relaymodel "github.com/modelgate/modelgate/internal/relay/model"
	relayimpl "github.com/modelgate/modelgate/internal/relay/service"
	"github.com/modelgate/modelgate/internal/runtime"
	"github.com/modelgate/modelgate/internal/server"
	systemmodel "github.com/modelgate/modelgate/internal/system/model"
	systemimpl "github.com/modelgate/modelgate/internal/system/service"
	"github.com/modelgate/modelgate/pkg/db"
)

var (
	appPath   string
	container do.Injector
	Version   string = "unknown"
	BuildTime string = "unknown"
	GitCommit string = "unknown"

	rootCmd = &cobra.Command{
		Use:   "modelgate",
		Short: "modelgate是一个基于Go语言的大模型接口转发服务",
		Long:  "modelgate是一个基于Go语言的大模型接口转发服务，支持多种大模型，包括OpenAI、Qwen、DeepSeek等。",
	}

	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "转发服务",
		Long:  "转发服务",
		Run: func(cmd *cobra.Command, args []string) {
			serverCfg := config.GetConfig().ApiServer
			sv, err := server.NewServer(&server.Profile{
				Name:      serverCfg.Name,
				Addr:      fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port),
				Container: container,
			})
			if err != nil {
				log.Fatalln(err)
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			go func() {
				if err := sv.Start(); err != nil {
					log.Fatalln(err)
				}
			}()

			<-quit
			log.Info("signal received")
			sv.Shutdown()
		},
	}

	adminCmd = &cobra.Command{
		Use:   "admin",
		Short: "管理后台服务",
		Long:  "管理后台服务",
		Run: func(cmd *cobra.Command, args []string) {
			serverCfg := config.GetConfig().AdminServer
			sv, err := server.NewH2CServer(&server.Profile{
				Name:      serverCfg.Name,
				Addr:      fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port),
				Container: container,
			})
			if err != nil {
				log.Fatalln(err)
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			go func() {
				if err := sv.Start(); err != nil {
					log.Fatalln(err)
				}
			}()

			<-quit
			log.Info("signal received")
			sv.Shutdown()
		},
	}

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "自动管理数据库表结构",
		Long:  "自动管理数据库表结构",
		Run: func(cmd *cobra.Command, args []string) {
			dbConn := do.MustInvoke[*gorm.DB](container)

			// 执行 GORM AutoMigrate
			dbConn.AutoMigrate(
				// Relay
				&relaymodel.Account{},
				&relaymodel.Provider{},
				&relaymodel.ProviderApiKey{},
				&relaymodel.ModelPricing{},
				&relaymodel.Model{},
				&relaymodel.Ledger{},
				&relaymodel.Request{},
				&relaymodel.RequestAttempt{},
				&relaymodel.AccountApiKey{},
				&relaymodel.RelayHourlyUsage{},
				&relaymodel.RelayUsage{},

				// System
				&systemmodel.Role{},
				&systemmodel.User{},
				&systemmodel.RefreshToken{},
				&systemmodel.Menu{},
				&systemmodel.Permission{},
			)

			dbConn.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					DoUpdates: clause.AssignmentColumns([]string{"id"}),
				}).Create(&relaymodel.RelayUsage{Model: db.Model{ID: int64(1)}})
		},
	}

	allCmd = &cobra.Command{
		Use:   "all",
		Short: "同时启动api、admin服务",
		Long:  "同时启动api、admin服务",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var svs []*server.Server

			g, ctx := errgroup.WithContext(ctx)
			{
				serverCfg := config.GetConfig().AdminServer
				svc, err := server.NewH2CServer(&server.Profile{
					Name:      serverCfg.Name,
					Addr:      fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port),
					Container: container,
				})
				if err != nil {
					log.Fatalln(err)
				}
				svs = append(svs, svc)
				g.Go(func() error {
					return svc.Start()
				})
			}
			{
				serverCfg := config.GetConfig().ApiServer
				svc, err := server.NewServer(&server.Profile{
					Name:      serverCfg.Name,
					Addr:      fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port),
					Container: container,
				})
				if err != nil {
					log.Fatalln(err)
				}
				svs = append(svs, svc)
				g.Go(func() error {
					return svc.Start()
				})
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

			select {
			case <-quit:
				log.Info("signal received")
			case <-ctx.Done():
				log.Error("server start failed:", ctx.Err())
			}

			for _, sv := range svs {
				sv.Shutdown()
			}
			if err := g.Wait(); err != nil {
				log.Fatalln(err)
			}
		},
	}
)

func main() {
	config.Version = Version
	config.BuildTime = BuildTime
	config.GitCommit = GitCommit
	container = do.New()

	log.SetFormatter(&log.JSONFormatter{})

	// 设置固定时区, Docker默认时区是UTC
	if l, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		time.Local = l
	}

	cobra.OnInitialize(func() {
		config.Init(container, appPath)
		config.CheckReady(container)

		systemimpl.Init(container)
		relayimpl.Init(container)
		runtime.Init(container)
		app.Init(container)
	})

	rootCmd.PersistentFlags().StringVarP(&appPath, "app_path", "p", "./", "app path")
	rootCmd.AddCommand(apiCmd, adminCmd, migrateCmd, allCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Exit(0)
}
