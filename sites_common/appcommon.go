package sites_common

import (
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-platform/platform_common"
)

// DB Module tables
const (
	// Database   Prefix
	DbPrefix = db_common.DB_COLLECTION_PREFIX
	// Collection Names
	Dbevents        = DbPrefix + "sites_event"
	Dbenquirys      = DbPrefix + "sites_enquiry"
	Dbpages         = DbPrefix + "sites_page"
	Dbmedia_gallery = DbPrefix + "sites_media_gallery"
)

// Sites Module table fields
const (
	FLD_BUSINESS_ID      = platform_common.FLD_BUSINESS_ID
	FLD_EVENT_ID         = "eventId"
	FLD_ENQUIRY_ID       = "enquiryId"
	FLD_PAGE_ID          = "pageId"
	FLD_MEDIA_GALLERY_ID = "mediagalleryId"
)

func GetServiceModuleCode() string {
	return "SITES"
}
