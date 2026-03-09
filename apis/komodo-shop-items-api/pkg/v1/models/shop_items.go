package models

// StockCode represents inventory availability status
type StockCode string

const (
	StockInStock              StockCode = "IS"
	StockOutOfStock           StockCode = "OS"
	StockLimitedSupply        StockCode = "LS"
	StockPreOrder             StockCode = "PO"
	StockSoldOut              StockCode = "SO"
	StockBackorder            StockCode = "BO"
	StockDiscontinued         StockCode = "DC"
	StockTemporarilyUnavail   StockCode = "TU"
)

type Product struct {
	ID                   string              `json:"id"`
	Slug                 string              `json:"slug"`
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	Brand                string              `json:"brand,omitempty"`
	Manufacturer         string              `json:"manufacturer,omitempty"`
	Status               string              `json:"status"`
	Currency             string              `json:"currency,omitempty"`
	Price                float64             `json:"price,omitempty"`
	CompareAtPrice       float64             `json:"compareAtPrice,omitempty"`
	Cost                 float64             `json:"cost,omitempty"`
	TaxCode              string              `json:"taxCode,omitempty"`
	TrackInventory       bool                `json:"trackInventory"`
	MinOrderQuantity     int                 `json:"minOrderQuantity,omitempty"`
	MaxOrderQuantity     int                 `json:"maxOrderQuantity,omitempty"`
	CustomizationOptions []CustomizationOption `json:"customizationOptions,omitempty"`
	AddOns               []AddOn             `json:"addOns,omitempty"`
	RelatedProductIds    []string            `json:"relatedProductIds,omitempty"`
	Variants             []Variant           `json:"variants"`
	Specs                map[string]any      `json:"specs,omitempty"`
	Meta                 *ProductMeta        `json:"meta,omitempty"`
	SEO                  *SEO                `json:"seo,omitempty"`
	CreatedAt            string              `json:"createdAt,omitempty"`
	UpdatedAt            string              `json:"updatedAt,omitempty"`
}

type ProductMeta struct {
	Tags       []string       `json:"tags,omitempty"`
	Categories []string       `json:"categories,omitempty"`
	IsPopular  bool           `json:"isPopular,omitempty"`
	IsFeatured bool           `json:"isFeatured,omitempty"`
	IsNew      bool           `json:"isNew,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
}

type SEO struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

type Variant struct {
	ID                string              `json:"id"`
	SKU               string              `json:"sku,omitempty"`
	UPC               string              `json:"upc,omitempty"`
	GTIN              string              `json:"gtin,omitempty"`
	EAN               string              `json:"ean,omitempty"`
	Model             string              `json:"model,omitempty"`
	Name              string              `json:"name"`
	Description       string              `json:"description,omitempty"`
	Price             float64             `json:"price"`
	CompareAtPrice    float64             `json:"compareAtPrice,omitempty"`
	Cost              float64             `json:"cost,omitempty"`
	TaxCode           string              `json:"taxCode,omitempty"`
	StockQty          int                 `json:"stockQty,omitempty"`
	StockCode         StockCode           `json:"stockCode,omitempty"`
	Images            []ProductImage      `json:"images,omitempty"`
	OptionCombination map[string]string   `json:"optionCombination,omitempty"`
	Weight            float64             `json:"weight,omitempty"`
	WeightUnit        string              `json:"weightUnit,omitempty"`
	Dimensions        *Dimensions         `json:"dimensions,omitempty"`
	RequiresShipping  bool                `json:"requiresShipping,omitempty"`
	ShippingClass     string              `json:"shippingClass,omitempty"`
	IsDefault         bool                `json:"isDefault,omitempty"`
}

type Dimensions struct {
	Length float64 `json:"length,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
	Unit   string  `json:"unit,omitempty"`
}

type ProductImage struct {
	URL        string      `json:"url"`
	Alt        string      `json:"alt,omitempty"`
	IsPrimary  bool        `json:"isPrimary,omitempty"`
	VariantIds []string    `json:"variantIds,omitempty"`
	Type       string      `json:"type,omitempty"`
	Spin360    *Spin360    `json:"spin360,omitempty"`
	Model3D    *Model3D    `json:"model3d,omitempty"`
}

type Spin360 struct {
	Frames     []string `json:"frames"`
	FrameCount int      `json:"frameCount"`
	StartFrame int      `json:"startFrame,omitempty"`
}

type Model3D struct {
	ModelURL     string   `json:"modelUrl"`
	Format       string   `json:"format"`
	TextureURLs  []string `json:"textureUrls,omitempty"`
	ThumbnailURL string   `json:"thumbnailUrl,omitempty"`
}

type CustomizationOption struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Type         string               `json:"type"`
	Required     bool                 `json:"required"`
	DisplayOrder int                  `json:"displayOrder"`
	Values       []CustomizationValue `json:"values"`
}

type CustomizationValue struct {
	ID             string    `json:"id"`
	Label          string    `json:"label"`
	Value          string    `json:"value"`
	PriceModifier  float64   `json:"priceModifier,omitempty"`
	HexColor       string    `json:"hexColor,omitempty"`
	ImageURL       string    `json:"imageUrl,omitempty"`
	StockCode      StockCode `json:"stockCode,omitempty"`
	StockQty       int       `json:"stockQty,omitempty"`
	IsDefault      bool      `json:"isDefault,omitempty"`
	Disabled       bool      `json:"disabled,omitempty"`
	DisabledReason string    `json:"disabledReason,omitempty"`
}

type AddOn struct {
	ID               string          `json:"id"`
	SKU              string          `json:"sku,omitempty"`
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	Manufacturer     string          `json:"manufacturer,omitempty"`
	Price            float64         `json:"price"`
	CompareAtPrice   float64         `json:"compareAtPrice,omitempty"`
	ImageURL         string          `json:"imageUrl,omitempty"`
	StockCode        StockCode       `json:"stockCode,omitempty"`
	StockQty         int             `json:"stockQty,omitempty"`
	Weight           float64         `json:"weight,omitempty"`
	RequiresShipping bool            `json:"requiresShipping,omitempty"`
	MaxQuantity      int             `json:"maxQuantity,omitempty"`
	IsRecommended    bool            `json:"isRecommended,omitempty"`
	CompatibleWith   *Compatibility  `json:"compatibleWith,omitempty"`
}

type Compatibility struct {
	OptionIds  []string `json:"optionIds,omitempty"`
	VariantIds []string `json:"variantIds,omitempty"`
}

type Service struct {
	ID                 string              `json:"id"`
	Slug               string              `json:"slug"`
	SKU                string              `json:"sku,omitempty"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Category           string              `json:"category"`
	Status             string              `json:"status"`
	Price              float64             `json:"price"`
	CompareAtPrice     float64             `json:"compareAtPrice,omitempty"`
	Cost               float64             `json:"cost,omitempty"`
	Currency           string              `json:"currency,omitempty"`
	TaxCode            string              `json:"taxCode,omitempty"`
	Duration           *ServiceDuration    `json:"duration,omitempty"`
	Images             []ServiceImage      `json:"images,omitempty"`
	ServiceOptions     []ServiceOption     `json:"serviceOptions,omitempty"`
	LocationTypes      []string            `json:"locationTypes"`
	Availability       *ServiceAvailability `json:"availability,omitempty"`
	Requirements       []string            `json:"requirements,omitempty"`
	IncludedItems      []string            `json:"includedItems,omitempty"`
	RelatedServiceIds  []string            `json:"relatedServiceIds,omitempty"`
	RelatedProductIds  []string            `json:"relatedProductIds,omitempty"`
	Meta               *ProductMeta        `json:"meta,omitempty"`
	SEO                *SEO                `json:"seo,omitempty"`
	CreatedAt          string              `json:"createdAt,omitempty"`
	UpdatedAt          string              `json:"updatedAt,omitempty"`
}

type ServiceImage struct {
	URL       string `json:"url"`
	Alt       string `json:"alt,omitempty"`
	IsPrimary bool   `json:"isPrimary,omitempty"`
	Type      string `json:"type,omitempty"`
}

type ServiceDuration struct {
	Estimated int    `json:"estimated"`
	Unit      string `json:"unit"`
	Min       int    `json:"min,omitempty"`
	Max       int    `json:"max,omitempty"`
}

type ServiceOption struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Type             string   `json:"type"`
	PriceModifier    float64  `json:"priceModifier"`
	DurationModifier float64  `json:"durationModifier,omitempty"`
	MaxBookings      int      `json:"maxBookings,omitempty"`
	IsDefault        bool     `json:"isDefault,omitempty"`
	RequiresProducts []string `json:"requiresProducts,omitempty"`
	CompatibleWith   []string `json:"compatibleWith,omitempty"`
}

type ServiceAvailability struct {
	DaysOfWeek        []int              `json:"daysOfWeek,omitempty"`
	TimeSlots         []TimeSlot         `json:"timeSlots,omitempty"`
	BlackoutDates     []string           `json:"blackoutDates,omitempty"`
	LeadTimeDays      int                `json:"leadTimeDays,omitempty"`
	MaxBookingsPerDay int                `json:"maxBookingsPerDay,omitempty"`
	ServiceAreaZips   []string           `json:"serviceAreaZipCodes,omitempty"`
	ServiceAreaRadius *ServiceAreaRadius `json:"serviceAreaRadius,omitempty"`
}

type TimeSlot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ServiceAreaRadius struct {
	Miles     float64 `json:"miles"`
	CenterZip string  `json:"centerZip"`
}

// API request/response models

type InventoryItem struct {
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	StockQty  int       `json:"stockQty"`
	StockCode StockCode `json:"stockCode"`
	Price     float64   `json:"price"`
	UpdatedAt string    `json:"updatedAt,omitempty"`
}

type InventoryResponse struct {
	Items []InventoryItem `json:"items"`
	Total int             `json:"total"`
}

type SuggestionRequest struct {
	UserID     string   `json:"userId,omitempty"`
	SKUs       []string `json:"skus,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

type SuggestionResponse struct {
	Suggestions []Product `json:"suggestions"`
	Total       int       `json:"total"`
}
