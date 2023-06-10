package sites_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sites/sites_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// MediaGallery DAO repositery
type MediaGalleryDao interface {
	// IntializeDao
	InitializeDao(client utils.Map, business string)
	// Create - Creat Collection
	Create(indata utils.Map) (utils.Map, error)
	// Get - Get By Code
	Get(mediagallertId string) (utils.Map, error)
	// List - List All Collection
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update - Update Collection
	Update(mediagalleryId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(mediagalleryId string) (int64, error)
	// Find - Find By Collection
	Find(filter string) (utils.Map, error)
}

// NewMediaGalleryDao  - Contruct Business Banner Dao
func NewMediaGalleryDao(Client utils.Map, business_id string) MediaGalleryDao {
	var daoMediaGallery MediaGalleryDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigind with correct after db Sercive created
	dbType, _ := db_common.GetDatabaseType(Client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoMediaGallery = &mongodb_repository.MediaGalleryMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
	// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoMediaGallery != nil {
		// Initialize the Dao
		daoMediaGallery.InitializeDao(Client, business_id)
	}

	return daoMediaGallery

}
