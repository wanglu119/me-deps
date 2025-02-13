package orm

import (
	"errors"
	"reflect"

	"github.com/wanglu119/me-deps/db/common"
	"xorm.io/xorm"
)

type Common struct {
}

func (c *Common) InsertOne(bean interface{}) error {
	_, err := x.InsertOne(bean)
	return err
}

func (c *Common) FindOne(bean interface{}) error {
	ok, err := x.Get(bean)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Not Found")
	}
	return nil
}

func (c *Common) Get(beans ...interface{}) (bool, error) {
	return x.Get(beans...)
}

func (c *Common) Find(beans interface{}, condiBean interface{}) error {
	var err error
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		err = x.Find(beans)
	} else {
		err = x.Find(beans, condiBean)
	}

	if err != nil {
		return err
	}
	return nil
}

func (c *Common) FindByPage(beans interface{}, condiBean interface{}, page *common.PageInfo) (int64, error) {
	var err error
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		err = x.Limit(page.Limit, page.Skip).Find(beans)
	} else {
		err = x.Limit(page.Limit, page.Skip).Find(beans, condiBean)
	}
	if err != nil {
		return 0, err
	}

	count, err := x.Count(condiBean)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Common) UpdateOne(bean interface{}, condiBean interface{}) error {
	_, err := x.Update(bean, condiBean)
	return err
}

func (c *Common) DeleteOne(bean interface{}) error {
	_, err := x.Delete(bean)
	return err
}

func (c *Common) DeleteMany(bean interface{}) error {
	return c.DeleteOne(bean)
}

func (c *Common) Count(bean interface{}) (int64, error) {
	return x.Count(bean)
}

func (c *Common) Sum(bean interface{}, colName string) (float64, error) {
	return x.Sum(bean, colName)
}

func (cm *Common) FindWithSort(beans interface{}, condiBean interface{}, sort *common.SortInfo) error {
	var err error
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		if sort == nil || reflect.ValueOf(sort).IsNil() {
			err = x.Find(beans)
		} else {
			if sort.SortType == common.SortDesc {
				err = x.Desc(sort.FieldName).Find(beans)
			} else {
				err = x.Asc(sort.FieldName).Find(beans)
			}
		}

	} else {
		if sort == nil || reflect.ValueOf(sort).IsNil() {
			err = x.Find(beans, condiBean)
		} else {
			if sort.SortType == common.SortDesc {
				err = x.Desc(sort.FieldName).Find(beans, condiBean)
			} else {
				err = x.Asc(sort.FieldName).Find(beans, condiBean)
			}
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (cm *Common) FindByPageWithSort(beans interface{}, condiBean interface{}, page *common.PageInfo, sort *common.SortInfo) (int64, error) {
	var err error
	if page == nil || reflect.ValueOf(page).IsNil() {
		page = &common.PageInfo{
			Limit: 10,
			Skip:  0,
		}
	}

	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		if sort == nil || reflect.ValueOf(sort).IsNil() {
			err = x.Limit(page.Limit, page.Skip).Find(beans)
		} else {
			if sort.SortType == common.SortDesc {
				err = x.Limit(page.Limit, page.Skip).Desc(sort.FieldName).Find(beans)
			} else {
				err = x.Limit(page.Limit, page.Skip).Asc(sort.FieldName).Find(beans)
			}
		}
	} else {
		if sort == nil || reflect.ValueOf(sort).IsNil() {
			err = x.Limit(page.Limit, page.Skip).Find(beans, condiBean)
		} else {
			if sort.SortType == common.SortDesc {
				err = x.Limit(page.Limit, page.Skip).Desc(sort.FieldName).Find(beans, condiBean)
			} else {
				err = x.Limit(page.Limit, page.Skip).Asc(sort.FieldName).Find(beans, condiBean)
			}
		}
	}
	if err != nil {
		return 0, err
	}

	count, err := x.Count(condiBean)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (cm *Common) SessionFindWithSort(s *xorm.Session, beans interface{}, condiBean interface{}, sort *common.SortInfo) error {
	var err error
	if condiBean == nil || reflect.ValueOf(condiBean).IsNil() {
		if sort == nil || reflect.ValueOf(sort).IsNil() {
			err = s.Find(beans)
		} else {
			if sort.SortType == common.SortDesc {
				err = s.Desc(sort.FieldName).Find(beans)
			} else {
				err = s.Asc(sort.FieldName).Find(beans)
			}
		}

	} else {
		if sort == nil || reflect.ValueOf(sort).IsNil() {
			err = s.Find(beans, condiBean)
		} else {
			if sort.SortType == common.SortDesc {
				err = s.Desc(sort.FieldName).Find(beans, condiBean)
			} else {
				err = s.Asc(sort.FieldName).Find(beans, condiBean)
			}
		}
	}
	if err != nil {
		return err
	}

	return nil
}
