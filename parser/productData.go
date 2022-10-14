package parser

type Retailer struct {
	RetailerName string `bson:"retailer" json:"retailer"`
	Quantity     int    `bson:"qty" json:"qty"`
}

type ProductData struct {
	GTIN          string     `bson:"gtin" json:"gtin"`
	GLN           string     `bson:"gln" json:"gln"`
	Name          string     `bson:"name" json:"name"`
	CompanyName   string     `bson:"company_name" json:"company_name"`
	StreetAddress string     `bson:"street_address" json:"street_address"`
	City          string     `bson:"city" json:"city"`
	CountryCode   string     `bson:"country_code" json:"country_code"`
	PostalCode    string     `bson:"postal_code" json:"postal_code"`
	Location      string     `bson:"location" json:"location"`
	Image         string     `bson:"img" json:"img"`
	Ingredients   []string   `bson:"ingredients" json:"ingredients"`
	Retailers     []Retailer `bson:"retailer_qty" json:"retailer_qty"`
}

func (txData *ProductData) PopulateWithMap(record map[string]string) {
	txData.GTIN = record["gtin"]
	txData.GLN = record["gln"]
	txData.Name = record["name"]
	txData.CompanyName = record["company_name"]
	txData.StreetAddress = record["street_address"]
	txData.City = record["city"]
	txData.CountryCode = record["country_code"]
	txData.PostalCode = record["postal_code"]
	txData.Location = record["location"]
	txData.Image = record["img"]
	//....
}

func (txData *ProductData) IsValid() bool {
	if txData.GTIN != "" {
		return true
	}
	return false
}

func (txData *ProductData) GetId() string {
	return txData.GLN
}
