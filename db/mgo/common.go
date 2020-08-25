package mgo

import (
	"reflect"
	"fmt"
	"errors"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"github.com/wanglu119/me-deps/db/common"
)

// BeforeInsertProcessor executed before an object is initially persisted to the database
type BeforeInsertProcessor interface {
	BeforeInsert()
}

// BeforeUpdateProcessor executed before an object is updated
type BeforeUpdateProcessor interface {
	BeforeUpdate()
}

// ======================================================

type ModelI interface {
	CreateIndex()
}

type Common struct {
}

func (cm *Common) getColl(data interface{}) *mongo.Collection{
	val := reflect.Indirect(reflect.ValueOf(data))
	var elemType reflect.Type
	if val.Kind() == reflect.Slice {
		elemType = val.Type().Elem()
	} else {
		elemType = val.Type()
	}
	
	return GetCollection(elemType)
}

// return count and 
func (cm *Common) MgoFind(beans interface{}, filter interface{}, findOpts ...*options.FindOptions) (int64, error) {
	
	sliceValue := reflect.Indirect(reflect.ValueOf(beans))
	
	if sliceValue.Kind() != reflect.Slice {
		return 0, errors.New("needs a pointer to a slice")
	}
	
	elemType := sliceValue.Type().Elem()
	
	var isPointer bool
	if elemType.Kind() == reflect.Ptr {
		isPointer = true
		elemType = elemType.Elem()
	}
	if elemType.Kind() == reflect.Ptr {
		return 0, errors.New("pointer to pointer is not supported")
	}
	
	coll := GetCollection(elemType)
	
	cur, err := coll.Find(ctx, filter, findOpts...)
	if err != nil {
		log.Error(fmt.Sprintf("query %s error: %v", elemType.Name(), err))
		return 0, err
	}
	defer cur.Close(ctx)
	
	var containerValueSetFunc func(*reflect.Value) 

	containerValueSetFunc = func(newValue *reflect.Value) {
		if isPointer {
			sliceValue.Set(reflect.Append(sliceValue, newValue.Elem().Addr()))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, newValue.Elem()))
		}
	}
	
	for cur.Next(ctx) {
		val := reflect.New(elemType)
		bean := val.Interface()
		err := cur.Decode(bean)
		if err != nil {
			log.Error(fmt.Sprintf("decode %s error: %v", elemType.Name(), err))
			continue
		}
		
		containerValueSetFunc(&val)
	}
	
	count,err := coll.CountDocuments(ctx, filter)
	if err != nil {
		log.Error(fmt.Sprintf("find count error: %v", err))
		return 0, err
	}
	
	return count, nil
}


func (cm *Common) FindOne(bean interface{}) error {
	
	coll := cm.getColl(bean)
	r := coll.FindOne(ctx, bean)
	if r == nil {
		return errors.New("Not Found")
	}
	
	err := r.Decode(bean)
	if err != nil {
		return err
	}
	return nil
}

func (cm *Common) Find(beans interface{},condiBean interface{}) error {
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		condiBean = bson.D{{}}
	}
	_, err := cm.MgoFind(beans, condiBean)
	if err != nil {
		return err
	}
	return nil
}

func (cm *Common) FindByPage(beans interface{},condiBean interface{}, page *common.PageInfo) (int64, error) {
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		condiBean = bson.D{{}}
	}
	
	ops := options.Find()
	ops.SetLimit(int64(page.Limit))
	ops.SetSkip(int64(page.Skip))
	
	count, err := cm.MgoFind(beans, condiBean, ops)
	if err != nil {
		return 0,err
	}
	return count, nil
}

func (cm *Common) InsertOne(bean interface{}) error {
	
	coll := cm.getColl(bean)
	
	if processor, ok := interface{}(bean).(BeforeInsertProcessor); ok {
		processor.BeforeInsert()
	}
	
	_, err := coll.InsertOne(ctx, bean)
	
	if err != nil {
		return err
	}
	
	return nil
}


func (cm *Common) UpdateOne(bean interface{}, condiBean interface{}) error {
	
	coll := cm.getColl(bean)
	
	if processor, ok := interface{}(bean).(BeforeUpdateProcessor); ok {
		processor.BeforeUpdate()
	}
	
	pByte, err := bson.Marshal(bean)
	if err != nil {
		return err
	}
	
	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return err
	}
	
	  _, err = coll.UpdateOne(ctx, condiBean, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	
	return nil
}

func (cm *Common) DeleteOne(bean interface{}) error {
	coll := cm.getColl(bean)
	_, err := coll.DeleteOne(ctx, bean)
	
	return err
}

func (cm *Common) DeleteMany(bean interface{}) error {
	coll := cm.getColl(bean)
	_, err := coll.DeleteMany(ctx, bean)
	
	return err
}

func (cm *Common) Count(bean interface{}) (int64, error) {
	coll := cm.getColl(bean)
	
	return coll.CountDocuments(ctx, bean)
}

func (cm *Common) FindWithSort(beans interface{},condiBean interface{}, sort *common.SortInfo) error {
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		condiBean = bson.D{{}}
	}
	opt := options.Find()
	var s bson.M
	if sort.SortType == common.SortDesc {
		s = bson.M{sort.FieldName:-1}
	} else {
		s = bson.M{sort.FieldName:1}
	}
	opt.SetSort(s)
	_, err := cm.MgoFind(beans, condiBean,opt)
	if err != nil {
		return err
	}
	return nil
}

func (cm *Common) FindByPageWithSort(beans interface{},condiBean interface{}, page *common.PageInfo, sort *common.SortInfo) (int64,error) {
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		condiBean = bson.D{{}}
	}
	
	ops := options.Find()
	ops.SetLimit(int64(page.Limit))
	ops.SetSkip(int64(page.Skip))
	var s bson.M
	if sort.SortType == common.SortDesc {
		s = bson.M{sort.FieldName:-1}
	} else {
		s = bson.M{sort.FieldName:1}
	}
	ops.SetSort(s)
	
	count, err := cm.MgoFind(beans, condiBean, ops)
	if err != nil {
		return 0,err
	}
	return count, nil
}

// ================================================


