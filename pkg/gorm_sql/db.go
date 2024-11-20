package gormsql

import (
	"context"
	"database/sql"
	"fmt"
	"golang_sql/pkg/utils"
	"log"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	mssql "github.com/microsoft/go-mssqldb"

	//mssql "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type GormSqlConfig struct {
	Server   string `mapstructure:"host"`
	Port     int    `mapstructure:port`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	Password string `mapstructure:"password"`
}

type Gorm struct {
	DB     *gorm.DB
	config *GormSqlConfig
}

func NewGorm(config *GormSqlConfig) (*gorm.DB, error) {

	var dataSourceName string

	if config.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	err := createDB(config)

	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", *&config.User, *&config.Password, *&config.Server, *&config.DBName)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 1
	var gormDb *gorm.DB

	err = backoff.Retry(func() error {

		var err error
		gormDb, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

		if err != nil {
			return errors.Errorf("failed to connect sqlserver: %v and connection information: %s", err, dataSourceName)
		}

		return nil

	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	return gormDb, err
}

func createDB(config *GormSqlConfig) error {

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", *&config.User, *&config.Password, *&config.Server, *&config.DBName)

	connector, err := mssql.NewConnector(dsn)
	if err != nil {
		log.Println(err)

	}
	sqlDb := sql.OpenDB(connector)

	var exists int
	rows, err := sqlDb.Query(fmt.Sprintf("SELECT 1 FROM  sys.databases WHERE name='%s'", config.DBName))
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = sqlDb.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
	if err != nil {
		return err
	}

	defer sqlDb.Close()

	return nil
}
func Migrate(gorm *gorm.DB, types ...interface{}) error {

	for _, t := range types {
		err := gorm.AutoMigrate(t)
		if err != nil {
			return err
		}
	}
	return nil
}
func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, db *gorm.DB) (*utils.ListResult[T], error) {

	var items []T
	var totalRows int64
	db.Model(items).Count(&totalRows)

	// generate where query
	query := db.Offset(listQuery.GetOffset()).Limit(listQuery.GetLimit()).Order(listQuery.GetOrderBy())

	if listQuery.Filters != nil {
		for _, filter := range listQuery.Filters {
			column := filter.Field
			action := filter.Comparison
			value := filter.Value

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				query = query.Where(whereQuery, value)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				query = query.Where(whereQuery, "%"+value+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(value, ",")
				query = query.Where(whereQuery, queryArray)
				break

			}
		}
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, errors.Wrap(err, "error in finding products.")
	}

	return utils.NewListResult[T](items, listQuery.GetSize(), listQuery.GetPage(), totalRows), nil
}
