package model

import (
	"fmt"
	"localcache/src/config"
	"log"
)

type Tenant struct {
	Id                 int     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Name               string  `json:"name" gorm:"column:name"`
	Organization       string  `json:"organization" gorm:"column:organization"`
	Desc               string  `json:"desc" gorm:"column:desc"`
	Balance            float64 `json:"balance" gorm:"column:balance"`
	Status             int     `json:"status" gorm:"column:status"`
	ClusterName        string  `json:"clusterName" gorm:"column:clusterName"`
	CreateTime         string  `json:"create_time" gorm:"column:create_time"`
	ModifyTime         string  `json:"modify_time" gorm:"column:modify_time;default:CURRENT_TIMESTAMP;NOT NULL"`
	PreSales           string  `json:"preSales" gorm:"column:preSales"`
	GroupName          string  `json:"groupName" gorm:"column:group_name"`
	NewOldStatus       int     `json:"newOldStatus" gorm:"column:new_old_status;default:1"`
	Salesman           string  `json:"salesman" gorm:"column:salesman"`
	AlertBalance       int     `json:"alertBalance" gorm:"column:alertBalance"`
	FullName           string  `json:"fullName" gorm:"column:fullName"`
	BillReceiver       string  `json:"billReceiver" gorm:"column:billReceiver"`
	PayType            string  `json:"payType" gorm:"column:payType;default:other;NOT NULL"`
	OrgFrom            string  `json:"orgFrom" gorm:"column:orgFrom;default:shumei;NOT NULL"`
	POrg               string  `json:"pOrg" gorm:"column:pOrg"`
	CustomerLevel      string  `json:"customerLevel" gorm:"column:customerLevel"`
	CustomerSuccess    string  `json:"customerSuccess" gorm:"column:customer_success"`
	CustomerIndustry   string  `json:"customerIndustry" gorm:"column:customerIndustry"`
	CustomerApp        string  `json:"customerApp" gorm:"column:customer_app"`
	CustomerIndustryCN string  `json:"customerIndustryCN" gorm:"-"`
	TechnicalSupport   string  `json:"technicalSupport" gorm:"column:technicalSupport"`
	IncomeLevel        string  `json:"incomeLevel" gorm:"column:incomeLevel"`
	Ext                string  `json:"ext"`
	IndustryFirst      string  `json:"industry_first" gorm:"column:industry_first"`
	IndustrySecond     string  `json:"industry_second" gorm:"column:industry_second"`
	MergeOrg           string  `json:"merge_org" gorm:"column:mergeOrg"`
	Source             string  `json:"source" gorm:"column:source"`
	CurrencyCode       string  `json:"currencyCode" gorm:"-"`
	TenantView
}

type TenantView struct {
	Account        string   `json:"account" gorm:"-"`
	FormerNameList []string `json:"formerNameList" gorm:"-"`
	AccessKey      string   `json:"accessKey" gorm:"-"`
	SecretKey      string   `json:"secretKey" gorm:"-"`
	Services       string   `json:"services" gorm:"-"`
	Tel            string   `json:"tel" gorm:"-"`
	Email          string   `json:"email" gorm:"-"`
	StatusCn       string   `json:"statusCn" gorm:"-"`
	CurrencyName   string   `json:"currencyName" gorm:"-"`
	CurrencyUnit   string   `json:"currencyUnit" gorm:"-"`
}

func (t *Tenant) TableName() string {
	return "saas_tenant"
}

//QueryTenantNameMap organization -> name
func (t *Tenant) QueryTenantNameMap() (tenantList []Tenant, err error) {
	db := config.Ds
	err = db.
		Table(t.TableName()).
		Select("organization,name").
		Find(&tenantList).Error
	return
}

func GetMap() (res map[string]interface{}, err error) {
	t := Tenant{}
	tt, err := t.QueryTenantNameMap()
	if err != nil {
		fmt.Println(err.Error())
	}
	res = map[string]interface{}{}
	for _, v := range tt {
		res[v.Organization] = v.Name
	}
	return
}

type App struct {
	Id           int    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Organization string `json:"organization" gorm:"column:organization"`
	AppId        string `json:"appId" gorm:"column:appId"`
	AppName      string `json:"appName" gorm:"column:appName"`
	CreateTime   string `json:"createTime" gorm:"column:createTime"`
	UpdateTime   string `json:"updateTime" gorm:"column:updateTime"`
	SourceAppId  string `json:"sourceAppId" gorm:"column:sourceAppId"`
	MailFlag     bool   `json:"mailFlag" gorm:"column:mailFlag"`
	User         string `json:"user" gorm:"column:user"`
}

//返回表名
func (app *App) TableName() string {
	return "sentry_rule_engine_app"
}

//应用下拉列表
func (app *App) List(where map[string]interface{}) (totalCount int64, AppData []App, err error) {
	db := config.Ds
	if err = db.
		Table(app.TableName()).
		Select("id,organization,appId,appName,sourceAppId").
		Where(where).
		Find(&AppData).Error; err != nil {
		return totalCount, AppData, err
	}
	return totalCount, AppData, nil
}

func GetAppMap() (map[string]interface{}, error) {
	a := App{}
	_, list, err := a.List(map[string]interface{}{})
	if err != nil {
		log.Println(err.Error())
	}
	res := map[string]interface{}{}
	for _, v := range list {
		res[v.Organization+v.AppId] = v.AppName
	}
	return res, nil
}


