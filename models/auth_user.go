package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type AuthUser struct {
	Id       int    `orm:"column(id);auto"`
	Name     string `orm:"column(name);size(50)" description:"姓名"`
	Userid   string `orm:"column(userid);size(50)" description:"企业微信id"`
	Gender   int8   `orm:"column(gender)" description:"0男，1女"`
	Avatar   string `orm:"column(avatar);size(255);null" description:"头像"`
	Username string `orm:"column(username);size(50);null" description:"用户名"`
	Password string `orm:"column(password);size(100);null" description:"密码"`
}

func (t *AuthUser) TableName() string {
	return "auth_user"
}

func init() {
	orm.RegisterModel(new(AuthUser))
}

// AddAuthUser insert a new AuthUser into database and returns
// last inserted Id on success.
func AddAuthUser(m *AuthUser) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAuthUserById retrieves AuthUser by Id. Returns error if
// Id doesn't exist
func GetAuthUserById(id int) (v *AuthUser, err error) {
	o := orm.NewOrm()
	v = &AuthUser{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAuthUser retrieves all AuthUser matches certain condition. Returns empty list if
// no records exist
func GetAllAuthUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AuthUser))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []AuthUser
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateAuthUser updates AuthUser by Id and returns error if
// the record to be updated doesn't exist
func UpdateAuthUserById(m *AuthUser) (err error) {
	o := orm.NewOrm()
	v := AuthUser{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAuthUser deletes AuthUser by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAuthUser(id int) (err error) {
	o := orm.NewOrm()
	v := AuthUser{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AuthUser{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
