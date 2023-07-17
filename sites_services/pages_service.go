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

type PagesService interface {
	// Create -Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Get -Find By Code
	Get(pageId string) (utils.Map, error)
	// list -List All records
	List(fliter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update update the Service
	Update(pageId string, indata utils.Map) (utils.Map, error)
	// Delete -Delete Srevice
	Delete(pageId string, delete_permanent bool) error
	// Find -from the item
	Find(filter string) (utils.Map, error)

	EndService()
}

type pagesBaseService struct {
	db_utils.DatabaseService
	daopage     sites_repository.PagesDao
	daoBusiness platform_repository.BusinessDao
	child       PagesService
	businessId  string
}

// NewBannerService - Construct Banner
func NewPagesService(props utils.Map) (PagesService, error) {
	funcode := sites_common.GetServiceModuleCode() + "M" + "01"

	p := pagesBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BannerService ")
	// Verify whether the business id data passed
	businessid, err := utils.GetMemberDataStr(props, sites_common.FLD_BUSINESS_ID)
	if err != nil {
		return p.errorReturn(err)
	}

	// Assign the BusinessId
	p.businessId = businessid
	p.initializeService()

	_, err = p.daoBusiness.Get(businessid)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given business_id is not exist"}
		return p.errorReturn(err)
	}

	p.child = &p

	return &p, err
}

// EndLoyaltyCardService - Close all the service
func (p *pagesBaseService) EndService() {
	log.Printf("Endservice")
	p.CloseDatabaseService()
}
func (p *pagesBaseService) initializeService() {
	log.Println("PagesService::GetBusinessDao ")
	p.daopage = sites_repository.NewPagesDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// Create -Create Service
func (p *pagesBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("PagesService::Create - Begin")
	var pageId string

	dataval, dataok := indata[sites_common.FLD_PAGE_ID]
	if dataok {
		pageId = strings.ToLower(dataval.(string))
	} else {
		pageId = utils.GenerateUniqueId("pag")
		log.Println("Unique Page Id", pageId)

	}
	indata[sites_common.FLD_BUSINESS_ID] = p.businessId
	indata[sites_common.FLD_PAGE_ID] = pageId

	data, err := p.daopage.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("PagesService::Creat - End")
	return data, nil
}

// Get - Find By Code
func (p *pagesBaseService) Get(pageId string) (utils.Map, error) {
	log.Printf("pagesBaseService::Get Begin %v", pageId)

	data, err := p.daopage.Get(pageId)

	log.Println("pagesBaseService::Get::End", data, err)
	return data, err
}

// list -List All records
func (p *pagesBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("pagesBaseService::FindAll - Begin")

	listdata, err := p.daopage.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("pagesBaseService::FindAll  -End")
	return listdata, nil

}

// Update - Update Service
func (p *pagesBaseService) Update(pageId string, indata utils.Map) (utils.Map, error) {

	log.Println("PagesService::Update -Begin")

	data, err := p.daopage.Update(pageId, indata)

	log.Println("PagesService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *pagesBaseService) Delete(pageId string, delete_permanent bool) error {

	log.Println("pagesBaseService ::Delete - Begin", pageId)

	if delete_permanent {
		result, err := p.daopage.Delete(pageId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(pageId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}
	log.Printf("PagesService :: Delete - End")
	return nil
}

func (p *pagesBaseService) Find(filter string) (utils.Map, error) {

	fmt.Println("PagesService ::FindByCode - Begin ", filter)

	data, err := p.daopage.Find(filter)
	log.Println("PagesService::FindByCode - End ", data, err)
	return data, err
}

func (p *pagesBaseService) errorReturn(err error) (PagesService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
