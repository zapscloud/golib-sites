package sites_repository

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-sites/sites_repository/mongodb_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// Enquiry Crud DAO repositery
type EnquiryDao interface {
	// IntializeDao
	InitializeDao(client utils.Map, business string)
	// Create - Creat Collection
	Create(indata utils.Map) (utils.Map, error)
	// Get - Get By Code
	Get(enquiryId string) (utils.Map, error)
	// List - List All Collection
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update - Update Collection
	Update(enquiryId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(enquiryId string) (int64, error)
	// Find - Find By Collection
	Find(filter string) (utils.Map, error)
}

// NewEnquiryDao  - Contruct Business Enquiry Dao
func NewEnquiryDao(Client utils.Map, business_id string) EnquiryDao {
	var daoEnquiry EnquiryDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigind with correct after db Sercive created
	dbType, _ := db_common.GetDatabaseType(Client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoEnquiry = &mongodb_repository.EnquiryMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
	// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoEnquiry != nil {
		// Initialize the Dao
		daoEnquiry.InitializeDao(Client, business_id)
	}

	return daoEnquiry

}
