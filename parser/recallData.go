package parser

type ExpiryDate struct {
	StartDate string `bson:"start_date" json:"start_date"`
	EndDate   string `bson:"end_date" json:"end_date"`
}

type Company struct {
	Name     string `bson:"name" json:"name"`
	Location string `bson:"location" json:"location"`
}

type Info struct {
	InformationSource string `bson:"information_src" json:"information_src"`
	URL               string `bson:"url" json:"url"`
	Summary           string `bson:"summary" json:"summary"`
	RecallDate        string `bson:"recall_date" json:"recall_date"`
	DateEnd           string `bson:"date_end" json:"date_end"`
}

type RecallData struct {
	GTIN        string       `bson:"gtin" json:"gtin"`
	ProductName string       `bson:"product_name" json:"product_name"`
	ExpiryDates []ExpiryDate `bson:"expiry_date" json:"expiry_date"`
	RiskLevel   int          `bson:"risk_level" json:"risk_level"`
	Company     Company      `bson:"company" json:"company"`
	Information []Info       `bson:"info" json:"info"`
}

func (txData *RecallData) PopulateWithMap(record map[string]string) {
	txData.GTIN = record["gtin"]
	// ...
}

func (txData *RecallData) IsValid() bool {
	if txData.GTIN != "" {
		return true
	}
	return false
}

func (txData *RecallData) GetId() string {
	return txData.Company.Name
}
