package service

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/system/model"
)

func (s *Service) DataMigrate(ctx context.Context) (err error) {
	if err = s.dataMigrationDao.AutoMigrate(); err != nil {
		return
	}
	migrations, err := s.dataMigrationDao.Find(ctx, &model.DataMigrationFilter{})
	if err != nil {
		return
	}
	executed := lo.Associate(migrations, func(m *model.DataMigration) (string, bool) {
		return m.Version, true
	})

	pendingMigrations, err := s.getPendingMigrations(executed)
	if err != nil {
		return
	}
	if len(pendingMigrations) == 0 {
		return
	}
	for _, migration := range pendingMigrations {
		log.Infof("executing migration: version: %s, description: %s, file_path: %s", migration.Version, migration.Description, migration.FilePath)
		content, err := os.ReadFile(migration.FilePath)
		if err != nil {
			return err
		}
		migration.RawSql = string(content)
		if err = s.dataMigrationDao.Create(ctx, migration); err != nil {
			return err
		}
	}
	return
}

// getPendingMigrations 获取待执行的迁移
func (s *Service) getPendingMigrations(executed map[string]bool) ([]*model.DataMigration, error) {
	migrationDir := config.GetPath("migration")
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var pending []*model.DataMigration
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		version, description := s.extractVersion(file.Name())
		if executed[version] {
			continue
		}
		pending = append(pending, &model.DataMigration{
			Version:     version,
			Description: description,
			FilePath:    filepath.Join(migrationDir, file.Name()),
		})
	}
	// 按版本号排序
	sort.Slice(pending, func(i, j int) bool {
		return semver.Compare(pending[i].Version, pending[j].Version) < 0
	})
	return pending, nil
}

// v1.0.0_init_schema.sql
func (s *Service) extractVersion(filename string) (string, string) {
	// 去掉 .sql 后缀，v 前缀
	name := strings.TrimSuffix(strings.TrimPrefix(filename, "v"), ".sql")
	version, description, _ := strings.Cut(name, "_")
	return version, description
}
