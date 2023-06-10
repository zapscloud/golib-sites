package mongodb_repository

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-sites/sites_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnquiryMonogoDBao -Enquiry Repository
type EnquiryMongoDBDao struct {
	Client     utils.Map
	businessId string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *EnquiryMongoDBDao) InitializeDao(Client utils.Map, businessId string) {
	log.Println("Initialaze Enquiry Mongodb DAO")
	p.Client = Client
	p.businessId = businessId
}

// Create  - Creat Collection
func (t *EnquiryMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("EnquiryMongoDBDao Create  - Begin ", indata)
	// sites Enquiry
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.Client, sites_common.Dbenquirys)
	if err != nil {
		log.Println("Error in insert", err)
		return utils.Map{}, err
	}
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	insertResult1, err := collection.InsertOne(ctx, indata)
	if err != nil {
		log.Println("Error in insert", err)
		return utils.Map{}, err

	}
	log.Println("Inserted a single document :", insertResult1.InsertedID)
	log.Println("EnquiryMongoDBDao Create - End ", indata[sites_common.FLD_ENQUIRY_ID])

	return t.Get(indata[sites_common.FLD_ENQUIRY_ID].(string))
}

// Get - Get By code
func (t *EnquiryMongoDBDao) Get(enquiryId string) (utils.Map, error) {

	// Get a Single Document
	var result utils.Map

	log.Println("EnquiryMongoDBDao Get  -  Begin ", enquiryId)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.Client, sites_common.Dbenquirys)
	log.Println("Get :: Got collection")

	filter := bson.D{{Key: sites_common.FLD_ENQUIRY_ID, Value: enquiryId}, {}}

	filter = append(filter,
		bson.E{Key: sites_common.FLD_BUSINESS_ID, Value: t.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Get::Got filter ", filter)
	singleResult := collection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		log.Println("Get :: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decord", err)
		return result, err

	}

	// Remove fileds from result
	result = db_common.AmendFldsForGet(result)

	log.Printf("Business EnquiryMongoDBDao    Get - End Found a single document: %+v\n", result)
	return result, nil
}

// List - List all Collections
func (t *EnquiryMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var result []utils.Map

	log.Println("EnquiryMongoDBDao  Begin  - Find All Collection Dao ", sites_common.Dbenquirys)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.Client, sites_common.Dbenquirys)
	if err != nil {
		return nil, err
	}

	log.Println("Get collection - find All Collection Dao ", filter, len(filter), sort, len(sort))

	opts := options.Find()

	filterdoc := bson.D{}
	if len(filter) > 0 {
		// filters,_ := strconv.Unquote (string(filter))
		err = bson.UnmarshalExtJSON([]byte(filter), true, &filterdoc)
		if err != nil {
			log.Println("unmarshal Ext JSON error ", err)
			log.Println(filterdoc)
		}
	}

	if len(sort) > 0 {
		var sortdoc interface{}
		err = bson.UnmarshalExtJSON([]byte(sort), true, &sortdoc)
		if err != nil {
			log.Println("sort unmarshal error ", sort)
		} else {
			opts.SetSort(sortdoc)
		}
	}

	if skip > 0 {
		log.Println(filterdoc)
		opts.SetSkip(skip)
	}

	if limit > 0 {
		log.Println(filterdoc)
		opts.SetLimit(limit)
	}
	filterdoc = append(filterdoc,
		bson.E{Key: sites_common.FLD_BUSINESS_ID, Value: t.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Parameter values ", filterdoc, opts)
	cursor, err := collection.Find(ctx, filterdoc, opts)
	if err != nil {
		return nil, err
	}

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	log.Println("EnquiryMongoDBDao  End - Find All Collection Dao", result)

	listdata := []utils.Map{}
	for idx, value := range result {
		log.Println("Item ", idx)
		// Remove fields from result
		value = db_common.AmendFldsForGet(value)
		listdata = append(listdata, value)
	}

	log.Println("Parameter values ", filterdoc)
	filtercount, err := collection.CountDocuments(ctx, filterdoc)
	if err != nil {
		return nil, err
	}

	basefilterdoc := bson.D{
		{Key: sites_common.FLD_BUSINESS_ID, Value: t.businessId},
		{Key: db_common.FLD_IS_DELETED, Value: false}}
	totalcount, err := collection.CountDocuments(ctx, basefilterdoc)
	if err != nil {
		return nil, err
	}

	response := utils.Map{
		db_common.LIST_SUMMARY: utils.Map{
			db_common.LIST_TOTALSIZE:    totalcount,
			db_common.LIST_FILTEREDSIZE: filtercount,
			db_common.LIST_RESULTSIZE:   len(listdata),
		},
		db_common.LIST_RESULT: listdata,
	}

	return response, nil
}

// Update - Update Collection
func (t *EnquiryMongoDBDao) Update(enquiryId string, indata utils.Map) (utils.Map, error) {

	log.Println("EnquiryMongoDBDao  Update - Begin")

	//Sites Enquiry
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.Client, sites_common.Dbenquirys)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)
	log.Printf("Update - Values %v", indata)

	filterEnquiry := bson.D{{Key: sites_common.FLD_ENQUIRY_ID, Value: enquiryId}}
	updateResult1, err := collection.UpdateOne(ctx, filterEnquiry, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult1.ModifiedCount)

	log.Println("Update - End")
	return t.Get(enquiryId)
}

// Delete - Delete Collection
func (t *EnquiryMongoDBDao) Delete(enquiryId string) (int64, error) {

	log.Println("EnquiryMongoDBDao  Delete - Begin ", enquiryId)

	// Sites Enquiry
	collection, ctx, err := mongo_utils.GetMongoDbCollection(t.Client, sites_common.Dbenquirys)
	if err != nil {
		return 0, err
	}
	optsEnquiry := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filterEnquiry := bson.D{{Key: sites_common.FLD_ENQUIRY_ID, Value: enquiryId}}
	resEnquiry, err := collection.DeleteOne(ctx, filterEnquiry, optsEnquiry)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("EnquiryMongoDBDao  Delete - End deleted %v documents\n", resEnquiry.DeletedCount)
	return resEnquiry.DeletedCount, nil
}

// Find - Find by Filter
func (p *EnquiryMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("EnquiryMongoDBDao  Find - Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.Client, sites_common.Dbenquirys)
	log.Println("EnquiryMongoDBDao Find:: Got Collection ", err)

	bfilter := bson.D{}
	err = bson.UnmarshalExtJSON([]byte(filter), true, &bfilter)
	if err != nil {
		fmt.Println("Error on filter Unmarshal", err)
	}
	bfilter = append(bfilter,
		bson.E{Key: sites_common.FLD_BUSINESS_ID, Value: p.businessId},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Find:: Got filter ", bfilter)
	singleResult := collection.FindOne(ctx, bfilter)
	if singleResult.Err() != nil {
		log.Println("Find:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("EnquiryMongoDBDao  Find -  End Found a single document: \n", err)
	return result, nil
}
