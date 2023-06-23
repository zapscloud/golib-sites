package sites_services

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform/platform_repository"
	"github.com/zapscloud/golib-sites/sites_common"
	"github.com/zapscloud/golib-sites/sites_repository"
	"github.com/zapscloud/golib-utils/utils"
)

type MediaGalleryService interface {
	// Create -Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Get -Find By Code
	Get(mediagalleryId string) (utils.Map, error)
	// list -List All records
	List(fliter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update update the Service
	Update(mediagalleryId string, indata utils.Map) (utils.Map, error)
	// Delete -Delete Srevice
	Delete(mediagalleryId string, delete_permanent bool) error
	// Find -from the item
	Find(filter string) (utils.Map, error)

	EndService()
}

type MediaGalleryBaseService struct {
	db_utils.DatabaseService
	daomediagallery sites_repository.MediaGalleryDao
	daoBusiness     platform_repository.BusinessDao
	child           MediaGalleryService
	businessId      string
}

// NewMediaGalleryService - Construct MediaGallery
func NewMediaGalleryService(props utils.Map) (MediaGalleryService, error) {
	funcode := sites_common.GetServiceModuleCode() + "M" + "01"

	p := MediaGalleryBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MediaGalleryService ")
	// Verify whether the business id data passed
	businessid, err := utils.GetMemberDataStr(props, sites_common.FLD_BUSINESS_ID)
	if err != nil {
		return nil, err
	}

	// Assign the BusinessId
	p.businessId = businessid
	p.initializeService()

	_, err = p.daoBusiness.Get(businessid)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given business_id is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the service
func (p *MediaGalleryBaseService) EndService() {
	log.Printf("Endservice")
	p.CloseDatabaseService()
}
func (p *MediaGalleryBaseService) initializeService() {
	log.Println("MediaGalleryService::GetBusinessDao ")
	p.daomediagallery = sites_repository.NewMediaGalleryDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// Create -Create Service
func (p *MediaGalleryBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("EnquiryService::Cerate - Beging")
	var mediagalleryId string

	dataval, dataok := indata[sites_common.FLD_MEDIA_GALLERY_ID]
	if dataok {
		mediagalleryId = strings.ToLower(dataval.(string))
	} else {
		mediagalleryId = utils.GenerateUniqueId("medgal")
		log.Println("Unique Media Gallery Id", mediagalleryId)

	}
	indata[sites_common.FLD_BUSINESS_ID] = p.businessId
	indata[sites_common.FLD_MEDIA_GALLERY_ID] = mediagalleryId

	data, err := p.daomediagallery.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("MediaGallery Service::Create - End")
	return data, nil
}

// Get - Find By Code
func (p *MediaGalleryBaseService) Get(mediagalleryId string) (utils.Map, error) {
	log.Printf("MediaGalleryBaseService::Get Begin %v", mediagalleryId)

	data, err := p.daomediagallery.Get(mediagalleryId)

	log.Println("MediaGalleryBaseService::Get::End", data, err)
	return data, err
}

// list -List All records
func (p *MediaGalleryBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("MediaGalleryBaseService::FindAll - Begin")

	listdata, err := p.daomediagallery.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("MediaGalleryBaseService::FindAll  -End")
	return listdata, nil

}

// Update - Update Service
func (p *MediaGalleryBaseService) Update(mediagalleryId string, indata utils.Map) (utils.Map, error) {

	log.Println("MediaGallery Service::Update - Begin")

	data, err := p.daomediagallery.Update(mediagalleryId, indata)

	log.Println("MediaGallery Service::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *MediaGalleryBaseService) Delete(mediagalleryId string, delete_permanent bool) error {

	log.Println("MediaGallery Service ::Delete - Begin", mediagalleryId)

	if delete_permanent {
		result, err := p.daomediagallery.Delete(mediagalleryId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(mediagalleryId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}
	log.Printf("MediaGallery Service :: Delete - End")
	return nil
}

func (p *MediaGalleryBaseService) Find(filter string) (utils.Map, error) {

	fmt.Println("MediaGallery Service::FindByCode - Begin ", filter)

	data, err := p.daomediagallery.Find(filter)
	log.Println("MediaGallery Service::FindByCode - End ", data, err)
	return data, err
}
