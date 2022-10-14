package parser

type ExpiryDate struct {
	StartDate string `bson:"start_date" json:"start_date"`
	EndDate   string `bson:"end_date" json:"end_date"`
}
type Product struct {
	GTIN        string       `bson:"gtin" json:"gtin"`
	Name        string       `bson:"name" json:"name"`
	ExpiryDates []ExpiryDate `bson:"expiry_date" json:"expiry_date"`
}

type RecallData struct {
	InformationSource string    `bson:"information_src" json:"information_src"`
	URL               string    `bson:"url" json:"url"`
	Summary           string    `bson:"summary" json:"summary"`
	RecallDate        string    `bson:"recall_date" json:"recall_date"`
	DateEnd           string    `bson:"date_end" json:"date_end"`
	Products          []Product `bson:"product" json:"product"`
	Cause             []string  `bson:"cause" json:"cause"`
	RiskLevel         int       `bson:"risk_level" json:"risk_level"`
	CompanyName       string    `bson:"company_name" json:"company_name"`
	Location          string    `bson:"location" json:"location"`
	Image             []string  `bson:"img" json:"img"`
}

func (txData *RecallData) PopulateWithMap(record map[string]string) {
	txData.InformationSource = record["information_src"]
	txData.URL = record["url"]
	//....
}

func (txData *RecallData) IsValid() bool {
	if txData.InformationSource != "" {
		return true
	}
	return false
}

func (txData *RecallData) GetId() string {
	return txData.CompanyName
}
