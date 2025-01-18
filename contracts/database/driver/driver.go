package driver

import (
	"gorm.io/gorm"

	"github.com/goravel/framework/contracts/database"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/testing"
)

type Driver interface {
	Config() database.Config
	Docker() (testing.DatabaseDriver, error)
	Gorm() (*gorm.DB, error)
	Grammar() schema.Grammar
	Processor() schema.Processor
	Schema(orm.Orm) schema.DriverSchema
}
