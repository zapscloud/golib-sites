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

type EnquiryService interface {
	// Create -Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Get -Find By Code
	Get(enquiryId string) (utils.Map, error)
	// list -List All records
	List(fliter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update update the Service
	Update(enquiryId string, indata utils.Map) (utils.Map, error)
	// Delete -Delete Srevice
	Delete(enquiryId string, delete_permanent bool) error
	// Find -from the item
	Find(filter string) (utils.Map, error)

	EndService()
}

type enquiryBaseService struct {
	db_utils.DatabaseService
	daoenquiry  sites_repository.EnquiryDao
	daoBusiness platform_repository.BusinessDao
	child       EnquiryService
	businessId  string
}

// NewEnquiryService - Construct Enquiry
func NewEnquiryService(props utils.Map) (EnquiryService, error) {
	funcode := sites_common.GetServiceModuleCode() + "M" + "01"

	p := enquiryBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("EnquiryService ")
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
func (p *enquiryBaseService) EndService() {
	log.Printf("endservice")
	p.CloseDatabaseService()
}
func (p *enquiryBaseService) initializeService() {
	log.Println("enquiryService::GetBusinessDao ")
	p.daoenquiry = sites_repository.NewEnquiryDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// Create -Create Service
func (p *enquiryBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("EnquiryService::Cerate - Begin")
	var enquiryId string

	dataval, dataok := indata[sites_common.FLD_ENQUIRY_ID]
	if dataok {
		enquiryId = strings.ToLower(dataval.(string))
	} else {
		enquiryId = utils.GenerateUniqueId("evn")
		log.Println("Unique enquiry Id", enquiryId)

	}
	indata[sites_common.FLD_BUSINESS_ID] = p.businessId
	indata[sites_common.FLD_ENQUIRY_ID] = enquiryId

	data, err := p.daoenquiry.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("EnquiryService  Creat - End")
	return data, nil
}

// Get - Find By Code
func (p *enquiryBaseService) Get(enquiryId string) (utils.Map, error) {
	log.Printf("EnquiryBaseService  Get - Begin %v", enquiryId)

	data, err := p.daoenquiry.Get(enquiryId)

	log.Println("EnquiryBaseService  Get - End", data, err)
	return data, err
}

// list -List All records
func (p *enquiryBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("EnquiryBaseService  FindAll - Begin")

	listdata, err := p.daoenquiry.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("EnquiryBaseService  FindAll - End")
	return listdata, nil

}

// Update - Update Service
func (p *enquiryBaseService) Update(enquiryId string, indata utils.Map) (utils.Map, error) {

	log.Println("EnquiryService  Update - Begin")

	data, err := p.daoenquiry.Update(enquiryId, indata)

	log.Println("EnquiryService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *enquiryBaseService) Delete(enquiryId string, delete_permanent bool) error {

	log.Println("EnquiryService  Delete - Begin", enquiryId)

	if delete_permanent {
		result, err := p.daoenquiry.Delete(enquiryId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(enquiryId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}
	log.Printf("EnquiryService  Delete - End")
	return nil
}

// Find  - Find Service
func (p *enquiryBaseService) Find(filter string) (utils.Map, error) {

	fmt.Println("EnquiryService  Find By Code - Begin ", filter)

	data, err := p.daoenquiry.Find(filter)
	log.Println("EnquiryService  Find By Code -  End ", data, err)
	return data, err
}
