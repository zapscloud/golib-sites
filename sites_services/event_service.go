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

type EventService interface {
	// Create -Create Service
	Create(indata utils.Map) (utils.Map, error)
	// Get -Find By Code
	Get(eventId string) (utils.Map, error)
	// list -List All records
	List(fliter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Update update the Service
	Update(eventId string, indata utils.Map) (utils.Map, error)
	// Delete -Delete Srevice
	Delete(eventId string, delete_permanent bool) error
	// Find -from the item
	Find(filter string) (utils.Map, error)

	EndService()
}

type eventBaseService struct {
	db_utils.DatabaseService
	daoevent    sites_repository.EventDao
	daoBusiness platform_repository.BusinessDao
	child       EventService
	businessId  string
}

// NewBannerService - Construct Banner
func NewEventService(props utils.Map) (EventService, error) {
	funcode := sites_common.GetServiceModuleCode() + "M" + "01"

	p := eventBaseService{}
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
func (p *eventBaseService) EndService() {
	log.Printf("endservice")
	p.CloseDatabaseService()
}
func (p *eventBaseService) initializeService() {
	log.Println("EventService::GetBusinessDao ")
	p.daoevent = sites_repository.NewEventDao(p.GetClient(), p.businessId)
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
}

// Create -Create Service
func (p *eventBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("EventService::Cerate - Beging")
	var eventId string

	dataval, dataok := indata[sites_common.FLD_EVENT_ID]
	if dataok {
		eventId = strings.ToLower(dataval.(string))
	} else {
		eventId = utils.GenerateUniqueId("evn")
		log.Println("Unique event Id", eventId)

	}
	indata[sites_common.FLD_BUSINESS_ID] = p.businessId
	indata[sites_common.FLD_EVENT_ID] = eventId

	data, err := p.daoevent.Create(indata)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("EventService::Creat - End")
	return data, nil
}

// Get - Find By Code
func (p *eventBaseService) Get(eventId string) (utils.Map, error) {
	log.Printf("eventBaseService::Get Begin %v", eventId)

	data, err := p.daoevent.Get(eventId)

	log.Println("eventBaseService::Get::End", data, err)
	return data, err
}

// list -List All records
func (p *eventBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("eventBaseService::FindAll - Begin")

	listdata, err := p.daoevent.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("eventBaseService::FindAll  -End")
	return listdata, nil

}

// Update - Update Service
func (p *eventBaseService) Update(eventId string, indata utils.Map) (utils.Map, error) {

	log.Println("EventService::Update -Begin")

	data, err := p.daoevent.Update(eventId, indata)

	log.Println("EventService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *eventBaseService) Delete(eventId string, delete_permanent bool) error {

	log.Println("eventService ::Delete - Begin", eventId)

	if delete_permanent {
		result, err := p.daoevent.Delete(eventId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(eventId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}
	log.Printf("EventService :: Delete - End")
	return nil
}

// Find - Find Service
func (p *eventBaseService) Find(filter string) (utils.Map, error) {

	fmt.Println("EventService ::FindByCode:: Begin ", filter)

	data, err := p.daoevent.Find(filter)
	log.Println("EventService::FindByCode :: End ", data, err)
	return data, err
}

func (p *eventBaseService) errorReturn(err error) (EventService, error) {
	// Close the Database Connection
	p.CloseDatabaseService()
	return nil, err
}
