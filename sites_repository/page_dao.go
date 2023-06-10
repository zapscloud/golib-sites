package sites_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sites/sites_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// PageDao - Card DAO Repository
type PageDao interface {
	// InitilaizeDao
	InitializeDao(client utils.Map, businessId string)
	// Create - Create collection
	Create(indata utils.Map) (utils.Map, error)
	// Get - Get by code
	Get(pageId string) (utils.Map, error)
	// List - List all collection
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update - Updata Collection
	Update(pageId string, indata utils.Map) (utils.Map, error)
	// Delete -Delete Collection
	Delete(pageId string) (int64, error)
	// Findn - Find by filter
	Find(filter string) (utils.Map, error)
}

// NewPageDao  - Contruct Business Banner Dao
func NewPageDao(Client utils.Map, business_id string) PageDao {
	var daoPage PageDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigind with correct after db Sercive created
	dbType, _ := db_common.GetDatabaseType(Client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoPage = &mongodb_repository.PageMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
	// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoPage != nil {
		// Initialize the Dao
		daoPage.InitializeDao(Client, business_id)
	}

	return daoPage

}
