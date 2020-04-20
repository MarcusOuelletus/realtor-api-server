package templates

// Property - Object for each individual property. This whole object is returned to the frontend.
type Property struct {
	Address         interface{} `bson:"Addr"`
	Community       interface{} `bson:"Community"`
	Municipality    interface{} `bson:"Municipality"`
	Area            interface{} `bson:"Area"`
	MLS             interface{} `bson:"Ml_num"`
	Price           interface{} `bson:"Lp_dol"`
	SoldPrice       interface{} `bson:"Sp_dol"`
	Baths           interface{} `bson:"Bath_tot"`
	Beds            interface{} `bson:"Br"`
	BedsPlus        interface{} `bson:"Br_plus"`
	Realtor         interface{} `bson:"Rltr"`
	Description     interface{} `bson:"Ad_text"`
	TourURL         interface{} `bson:"Tour_url"`
	SaleType        interface{} `bson:"S_r"`
	BuildingType    interface{} `bson:"Style"`
	PropertyType    interface{} `bson:"Type_own1_out"`
	GarageType      interface{} `bson:"Gar_type"`
	Sqft            interface{} `bson:"Sqft"`
	FireplaceCount  interface{} `bson:"Fpl_num"`
	Pool            interface{} `bson:"Pool"`
	CrossStreet     interface{} `bson:"Cross_st"`
	Unit            interface{} `bson:"Apt_num"`
	AirConditioning interface{} `bson:"A_c"`
	FuelType        interface{} `bson:"Fuel"`
	Rooms           interface{} `bson:"Rms"`
	RoomsPlus       interface{} `bson:"Rooms_plus"`
	CentralVac      interface{} `bson:"Central_vac"`
	YearBuilt       interface{} `bson:"Yr_built"`
	Taxes           interface{} `bson:"Taxes"`
	PostalCode      interface{} `bson:"Zip"`
}

// SearchRequest - Object representing the incoming JSON from a user's property search.
type SearchRequest struct {
	Type         string `json:"type"`
	Area         string `json:"area"`
	MinPrice     string `json:"minprice"`
	MaxPrice     string `json:"maxprice"`
	Beds         string `json:"beds"`
	Baths        string `json:"baths"`
	BuildingType string `json:"btype"`
	SaleType     string `json:"stype"`
	GarageType   string `json:"garage"`
	YearBuilt    string `json:"year"`
	Address      string `json:"address"`
	MLS          string `json:"mls"`
	Page         int    `json:"page"`
	All          bool   `json:"all"`
	Sold         bool   `json:"sold"`
	UserID       string `json:"userID"`
	Unit         string `json:"unit"`
}

// fieldsStruct - Understandable version of jibberish TREB column names
type fieldsStruct struct {
	Address        string
	Area           string
	Community      string
	Municipality   string
	MLS            string
	Price          string
	Beds           string
	BedsPlus       string
	Baths          string
	Status         string
	BuildingType   string
	Unit           string
	SaleType       string
	GarageType     string
	YearBuilt      string
	Timestamp      string
	MarcPix        string
	Brokerage      string
	PostalCode     string
	SoldPrice      string
	Expiry         string
	DisplayAddress string
}

// Fields - Easy way to reference TREB column names
var Fields = fieldsStruct{
	Address:        "Addr",
	Area:           "Area",
	Community:      "Community",
	Municipality:   "Municipality",
	MLS:            "Ml_num",
	Price:          "Lp_dol",
	Beds:           "Br",
	BedsPlus:       "Br_plus",
	Baths:          "Bath_tot",
	Status:         "Status",
	BuildingType:   "Style",
	Unit:           "Apt_num",
	SaleType:       "S_r",
	GarageType:     "Gar_type",
	YearBuilt:      "Yr_built",
	Timestamp:      "Timestamp_sql",
	MarcPix:        "Marc_pix",
	Brokerage:      "Rltr",
	PostalCode:     "Zip",
	SoldPrice:      "Sp_dol",
	Expiry:         "Xd",
	DisplayAddress: "Disp_addr",
}

// SearchReturn - All the database field names being returned when querying a property.
var SearchReturn = &[]string{
	"Addr",
	"Community",
	"Municipality",
	"Area",
	"Ml_num",
	"Lp_dol",
	"Sp_dol",
	"Bath_tot",
	"Br",
	"Br_plus",
	"Ad_text",
	"Tour_url",
	"Rltr",
	"S_r",
	"Type_own1_out",
	"Gar_type",
	"Sqft",
	"Style",
	"Fpl_num",
	"Pool",
	"Cross_st",
	"Apt_num",
	"A_c",
	"Fuel",
	"Rms",
	"Rooms_plus",
	"Central_vac",
	"Taxes",
	"Yr_built",
	"Zip",
}

// InvalidToken - this will cause the browser to show a message to refresh the page
const InvalidToken = "invalid user token"

type APIResponse struct {
	Token    string      `json:"token"`
	Error    string      `json:"error"`
	Code     int         `json:"code"`
	Response interface{} `json:"response"`
}

/*

Response Error Codes
--------------------

100s: Query Error

100: no results

----------------

800s: Permission Error

800: user accessing sold listings without signing up
801: realtor hasn't purchased access to sold listings

----------------

*/
