package parser

type ImportData struct {
	GTIN          string `bson:"gtin" json:"gtin"`
	BrandName     string `bson:"brand_name" json:"brand_name"`
	ItemDesc      string `bson:"item_desc" json:"item_desc"`
	TraderName    string `bson:"trader_name" json:"trader_name"`
	CityOfOrigin  string `bson:"city_of_origin" json:"city_of_origin"`
	CityLocation  string `bson:"city_location" json:"city_location"`
	CountryCode   string `bson:"country_code" json:"country_code"`
	ProductQty    string `bson:"product_qty" json:"product_qty"`
	DateOfArrival string `bson:"date_of_arrival" json:"date_of_arrival"`
	HealthCertNum string `bson:"health_cert_num" json:"health_cert_num"`
	UENNo         string `bson:"uen_no" json:"uen_no"`
	FarmCode      string `bson:"farm_code" json:"farm_code"`
	CCPNo         string `bson:"ccp_no" json:"ccp_no"`
	InwardMode    string `bson:"inward_mode" json:"inward_mode"`
}

func (txData *ImportData) PopulateWithMap(record map[string]string) {
	txData.GTIN = record["gtin"]
	txData.BrandName = record["brand_name"]
	txData.ItemDesc = record["item_desc"]
	txData.TraderName = record["trader_name"]
	txData.CityOfOrigin = record["city_of_origin"]
	txData.CityLocation = record["city_location"]
	txData.CountryCode = record["country_code"]
	txData.ProductQty = record["product_qty"]
	txData.DateOfArrival = record["date_of_arrival"]
	txData.HealthCertNum = record["health_cert_num"]
	txData.UENNo = record["uen_no"]
	txData.FarmCode = record["farm_code"]
	txData.CCPNo = record["ccp_no"]
	txData.InwardMode = record["inward_mode"]
}

func (txData *ImportData) IsValid() bool {
	if txData.GTIN != "" {
		return true
	}
	return false
}

func (txData *ImportData) GetId() string {
	return txData.BrandName
}
