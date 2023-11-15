package UpdateBedStatus

import (
	"github.com/project-flogo/core/data/coerce"
)

type Input struct {
	IP         string `md:"IP,required"`
	CustomerId string `md:"CustomerId,required"`
	Username   string `md:"Username,required"`
	Password   string `md:"Password,required"`
	MAC 	   string `md:"MAC,required"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["IP"])
	i.IP = strVal

	strVal, _ = coerce.ToString(values["CustomerId"])
	i.CustomerId = strVal

	strVal, _ = coerce.ToString(values["Username"])
	i.Username = strVal

	strVal, _ = coerce.ToString(values["Password"])
	i.Password = strVal

	strVal, _ = coerce.ToString(values["MAC"])
	i.MAC = strVal
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"IP":         i.IP,
		"CustomerId": i.CustomerId,
		"Username":   i.Username,
		"Password":   i.Password,
		"MAC": 	  i.MAC,
	}
}

type Output struct {
	Status bool `md:"Status"`
	BedStatus string `md:"BedStatus"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	boolStatus, _ := coerce.ToBool(values["Status"])
	o.Status = boolStatus
	bedStatus, _ := coerce.ToString(values["BedStatus"])
	o.BedStatus = bedStatus
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Status": o.Status,
		"BedStatus":o.BedStatus,
	}
}

type Response struct {
	ElapsedTimeInMillseconds float64 `json:"ElapsedTimeInMillseconds"`
	ErrorMessage             string  `json:"ErrorMessage"`
	SuccessMessage           string  `json:"SuccessMessage"`
	HasError                 bool    `json:"HasError"`
	ID                       int     `json:"Id"`
}

type Device struct {
	ItemID                   int       `json:"ItemId"`
}

type Staff struct {
	BedStatus                string `json:"BedStatus"`
	AssocItemID              int     `json:"AssocItemID"`
}