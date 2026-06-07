package model

import (
	"time"

	"github.com/warerastats/models/models/enums"
)

// Union interfaces. Defined here so hand-written models (e.g. ActivityConnection)
// can reference them; gqlgen autobinds the GraphQL unions to these.
type (
	Activity     interface{ IsActivity() }
	SearchResult interface{ IsSearchResult() }
	Entity       interface{ IsEntity() }
)

// These hand-written models carry hidden foreign-key IDs and raw maps that are
// not part of the GraphQL schema. Edge fields (country, seller, …) and derived
// fields (wealth, skills, …) are resolved on demand by field resolvers that
// read these hidden fields off the parent. gqlgen binds to these via `autobind`
// in gqlgen.yml and generates everything else (unions, enums, Query model).

// ---------- shared leaf types ----------

// FloatEntry is one entry of a string->float map exposed as a typed pair list.
type FloatEntry struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

// IntEntry is one entry of a string->int map exposed as a typed pair list.
type IntEntry struct {
	Key   string `json:"key"`
	Value int32  `json:"value"`
}

// Taxes is a country's tax rates.
type Taxes struct {
	Income   float64 `json:"income"`
	Market   float64 `json:"market"`
	SelfWork float64 `json:"selfWork"`
}

// Ethics is a party's four-axis ethics vector.
type Ethics struct {
	Militarism    int32 `json:"militarism"`
	Isolationism  int32 `json:"isolationism"`
	Imperialism   int32 `json:"imperialism"`
	Industrialism int32 `json:"industrialism"`
}

// SkillSet is a user's 14-axis skill snapshot.
type SkillSet struct {
	Energy           int32 `json:"energy"`
	Health           int32 `json:"health"`
	Hunger           int32 `json:"hunger"`
	Attack           int32 `json:"attack"`
	Companies        int32 `json:"companies"`
	Entrepreneurship int32 `json:"entrepreneurship"`
	Production       int32 `json:"production"`
	CriticalChance   int32 `json:"criticalChance"`
	CriticalDamages  int32 `json:"criticalDamages"`
	Armor            int32 `json:"armor"`
	Precision        int32 `json:"precision"`
	Dodge            int32 `json:"dodge"`
	LootChance       int32 `json:"lootChance"`
	Management       int32 `json:"management"`
}

// ---------- trackers ----------

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Level        int32     `json:"level"`
	AvatarURL    string    `json:"avatarUrl"`
	MilitaryRank int32     `json:"militaryRank"`
	LastDate     time.Time `json:"lastDate"`
	LastSeen     time.Time `json:"lastSeen"`

	CountryID string
	CompanyID *string
	PartyID   *string
	MuID      *string
	WealthMap map[string]float64
	SkillsMap map[string]int
}

func (User) IsSearchResult() {}
func (User) IsEntity()       {}

type Country struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Code  string  `json:"code"`
	Money float64 `json:"money"`
	Taxes *Taxes  `json:"taxes"`

	Specialisation *string
	RulingPartyID  *string
}

func (Country) IsSearchResult() {}
func (Country) IsEntity()       {}

type Region struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	IsCapital         bool    `json:"isCapital"`
	IsLinkedToCapital bool    `json:"isLinkedToCapital"`
	Resistance        float64 `json:"resistance"`
	MaxResistance     float64 `json:"maxResistance"`

	CountryID        string
	InitialCountryID string
	NeighborIDs      []string
}

type Party struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatarUrl"`
	Ethics      *Ethics   `json:"ethics"`
	LastUpdated time.Time `json:"lastUpdated"`

	CountryID     string
	RegionID      string
	LeaderUserID  string
	MemberUserIDs []string
}

func (Party) IsSearchResult() {}
func (Party) IsEntity()       {}

type Mu struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL string  `json:"avatarUrl"`
	Level     int32   `json:"level"`
	Hq        int32   `json:"hq"`
	Dorms     int32   `json:"dorms"`
	MercRep   float64 `json:"mercRep"`

	OwnerUserID   string
	RegionID      string
	MemberUserIDs []string
}

func (Mu) IsSearchResult() {}
func (Mu) IsEntity()       {}

type Battle struct {
	ID              string      `json:"id"`
	AttackerDamages int32       `json:"attackerDamages"`
	DefenderDamages int32       `json:"defenderDamages"`
	WinnerSide      *enums.Side `json:"winnerSide"`
	IsActive        bool        `json:"isActive"`
	EndedAt         *time.Time  `json:"endedAt"`
	LastUpdated     time.Time   `json:"lastUpdated"`

	AttackerCountryID string
	DefenderCountryID string
	AttackerRegionID  *string
	DefenderRegionID  string
}

type Item struct {
	ID       string           `json:"id"`
	ItemCode string           `json:"itemCode"`
	State    int32            `json:"state"`
	Status   enums.ItemStatus `json:"status"`

	OwnerUserID string
	SkillsMap   map[string]float64
}

type Damage struct {
	ID           string     `json:"id"`
	Side         enums.Side `json:"side"`
	MilitaryRank int32      `json:"militaryRank"`
	Damages      int32      `json:"damages"`
	At           time.Time  `json:"at"`
	Ammo         *string    `json:"ammo"`

	BattleID  string
	UserID    string
	CountryID string
	MuID      *string
	PartyID   *string
	SkillID   string
	WeaponID  *string
	HelmetID  *string
	ChestID   *string
	PantsID   *string
	BootsID   *string
	GlovesID  *string
}

type Skill struct {
	ID    string    `json:"id"`
	Since time.Time `json:"since"`
	Set   *SkillSet `json:"set"`

	UserID string
}

type Company struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ItemCode string `json:"itemCode"`

	UserID   string
	RegionID string
}

type Employee struct {
	ID                     string    `json:"id"`
	Wage                   float64   `json:"wage"`
	Fidelity               int32     `json:"fidelity"`
	JoinedAt               time.Time `json:"joinedAt"`
	LastFidelityIncreaseAt time.Time `json:"lastFidelityIncreaseAt"`

	UserID     string
	CompanyID  string
	EmployerID string
}

type TradeOffer struct {
	ID        string          `json:"id"`
	ItemCode  string          `json:"itemCode"`
	Side      enums.TradeSide `json:"side"`
	Quantity  int32           `json:"quantity"`
	Fulfilled int32           `json:"fulfilled"`
	Cancelled bool            `json:"cancelled"`
	Price     float64         `json:"price"`
	Since     time.Time       `json:"since"`

	UserID    string
	CountryID *string
	MuID      *string
}

// ---------- transactions (Activity union) ----------

type MarketTransaction struct {
	ID    string  `json:"id"`
	Money float64 `json:"money"`

	SellerID string
	BuyerID  string
	ItemID   string
}

func (MarketTransaction) IsActivity() {}

type WageTransaction struct {
	ID               string  `json:"id"`
	Money            float64 `json:"money"`
	ProductionPoints int32   `json:"productionPoints"`

	EmployeeID string
	EmployerID string
}

func (WageTransaction) IsActivity() {}

type CaseTransaction struct {
	ID   string `json:"id"`
	Case string `json:"case"`

	UserID string
	ItemID string
}

func (CaseTransaction) IsActivity() {}

type TradeTransaction struct {
	ID           string  `json:"id"`
	ItemCode     string  `json:"itemCode"`
	Money        float64 `json:"money"`
	Quantity     int32   `json:"quantity"`
	TimeTillSale int64   `json:"timeTillSale"`

	SellerID        string
	BuyerID         string
	SellerMuID      *string
	BuyerMuID       *string
	SellerCountryID *string
	BuyerCountryID  *string
}

func (TradeTransaction) IsActivity() {}

type CraftTransaction struct {
	ID         string `json:"id"`
	ScrapsCost int32  `json:"scrapsCost"`

	UserID string
	ItemID string
}

func (CraftTransaction) IsActivity() {}

type DismantleTransaction struct {
	ID             string `json:"id"`
	ScrapsReceived int32  `json:"scrapsReceived"`

	UserID string
	ItemID string
}

func (DismantleTransaction) IsActivity() {}

type LootTransaction struct {
	ID string `json:"id"`

	UserID string
	ItemID string
}

func (LootTransaction) IsActivity() {}

// ActivityConnection is a cursor-paginated page of a user's activity feed.
type ActivityConnection struct {
	Edges       []Activity `json:"edges"`
	EndCursor   *string    `json:"endCursor"`
	HasNextPage bool       `json:"hasNextPage"`
}

// ---------- events ----------

type UserNameChange struct {
	ID       string    `json:"id"`
	At       time.Time `json:"at"`
	Username string    `json:"username"`
	UserID   string
}

type UserCountryChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	UserID    string
	CountryID *string
}

type UserPartyChange struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	UserID  string
	PartyID *string
}

type UserCompanyChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	UserID    string
	CompanyID *string
}

type UserMuChange struct {
	ID     string    `json:"id"`
	At     time.Time `json:"at"`
	UserID string
	MuID   *string
}

type UserSkillChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	UserID    string
	SkillsMap map[string]int
}

type PartyNameChange struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	Name    string    `json:"name"`
	PartyID string
}

type PartyLeaderChange struct {
	ID           string    `json:"id"`
	At           time.Time `json:"at"`
	PartyID      string
	LeaderUserID string
}

type PartyDescriptionChange struct {
	ID          string    `json:"id"`
	At          time.Time `json:"at"`
	Description string    `json:"description"`
	PartyID     string
}

type PartyEthicsChange struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	Ethics  *Ethics   `json:"ethics"`
	PartyID string
}

type CountryRulingPartyChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	CountryID string
	PartyID   *string
}

type CountrySpecialisationChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	ItemCode  *string   `json:"itemCode"`
	CountryID string
}

type RegionOwnerChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	RegionID  string
	CountryID string
}

type RegionDeposit struct {
	ID           string     `json:"id"`
	At           time.Time  `json:"at"`
	Type         *string    `json:"type"`
	StartsAt     *time.Time `json:"startsAt"`
	EndsAt       *time.Time `json:"endsAt"`
	BonusPercent *float64   `json:"bonusPercent"`
	RegionID     string
}

type RegionStrategicResource struct {
	ID       string    `json:"id"`
	At       time.Time `json:"at"`
	Resource *string   `json:"resource"`
	RegionID string
}

type CompanyRegionChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	CompanyID string
	RegionID  string
}

type CompanyItemCodeChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	ItemCode  string    `json:"itemCode"`
	CompanyID string
}

type EmployeeWageChange struct {
	ID        string    `json:"id"`
	At        time.Time `json:"at"`
	Wage      float64   `json:"wage"`
	UserID    string
	CompanyID string
}

type MuNameChange struct {
	ID   string    `json:"id"`
	At   time.Time `json:"at"`
	Name string    `json:"name"`
	MuID string
}

type MuOwnerChange struct {
	ID          string    `json:"id"`
	At          time.Time `json:"at"`
	MuID        string
	OwnerUserID string
}

type MuMercenaryReputationChange struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	MercRep float64   `json:"mercRep"`
	MuID    string
}

type BattleOrderChange struct {
	ID       string    `json:"id"`
	At       time.Time `json:"at"`
	Side     string    `json:"side"`
	Kind     string    `json:"kind"`
	Action   string    `json:"action"`
	BattleID string
	EntityID string
}

// ---------- processed: user ----------

type UserFinanceReport struct {
	ID              string    `json:"id"`
	DayStart        time.Time `json:"dayStart"`
	WagesPaid       float64   `json:"wagesPaid"`
	WagesEarned     float64   `json:"wagesEarned"`
	ItemsBought     float64   `json:"itemsBought"`
	ItemsSold       float64   `json:"itemsSold"`
	EquipBought     float64   `json:"equipBought"`
	EquipSold       float64   `json:"equipSold"`
	ValueDismantled float64   `json:"valueDismantled"`
	CasesOpened     int32     `json:"casesOpened"`
	CasesNet        float64   `json:"casesNet"`
	UserID          string
}

type UserFlipEvent struct {
	ID          string    `json:"id"`
	ItemCode    string    `json:"itemCode"`
	Quantity    int32     `json:"quantity"`
	BuyCost     float64   `json:"buyCost"`
	SellRevenue float64   `json:"sellRevenue"`
	Profit      float64   `json:"profit"`
	HeldMs      int64     `json:"heldMs"`
	At          time.Time `json:"at"`
	UserID      string
}

type FlipLot struct {
	TradeID   string    `json:"tradeId"`
	Quantity  int32     `json:"quantity"`
	UnitPrice float64   `json:"unitPrice"`
	BoughtAt  time.Time `json:"boughtAt"`
}

type FlipLotGroup struct {
	ItemCode string     `json:"itemCode"`
	Lots     []*FlipLot `json:"lots"`
}

type UserFlipState struct {
	TotalFlips  int32     `json:"totalFlips"`
	TotalProfit float64   `json:"totalProfit"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserID      string
	OpenLotsMap map[string][]FlipLot
}

type UserInventory struct {
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    string
	ItemsMap  map[string]int
}

type UserBattleParticipation struct {
	TotalDamage            int64     `json:"totalDamage"`
	BattlesParticipated    int32     `json:"battlesParticipated"`
	NegativeDamage         int64     `json:"negativeDamage"`
	OwnCountryBattles      int32     `json:"ownCountryBattles"`
	OwnCountryParticipated int32     `json:"ownCountryParticipated"`
	MuOrderBattles         int32     `json:"muOrderBattles"`
	MuOrderParticipated    int32     `json:"muOrderParticipated"`
	UpdatedAt              time.Time `json:"updatedAt"`
	UserID                 string
}

// ---------- processed: country ----------

type TaxHijack struct {
	Amount    float64 `json:"amount"`
	CountryID string
}

type TaxSource struct {
	Total     float64 `json:"total"`
	Hijacked  float64 `json:"hijacked"`
	CorePct   float64 `json:"corePct"`
	CountryID string
}

type CountryTaxFlow struct {
	ID            string       `json:"id"`
	HourStart     time.Time    `json:"hourStart"`
	TotalTax      float64      `json:"totalTax"`
	HijackedIn    float64      `json:"hijackedIn"`
	CoreEarned    float64      `json:"coreEarned"`
	NonCoreEarned float64      `json:"nonCoreEarned"`
	HijackedOut   float64      `json:"hijackedOut"`
	Hijackers     []*TaxHijack `json:"hijackers"`
	Sources       []*TaxSource `json:"sources"`
	CountryID     string
}

type CountryFlipEvent struct {
	ID          string    `json:"id"`
	ItemCode    string    `json:"itemCode"`
	Quantity    int32     `json:"quantity"`
	BuyCost     float64   `json:"buyCost"`
	SellRevenue float64   `json:"sellRevenue"`
	Profit      float64   `json:"profit"`
	At          time.Time `json:"at"`
	CountryID   string
}

type CountryFlipState struct {
	TotalTrades int32     `json:"totalTrades"`
	Profitable  int32     `json:"profitable"`
	TotalProfit float64   `json:"totalProfit"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CountryID   string
}

type InventoryLot struct {
	TradeID   string    `json:"tradeId"`
	Quantity  int32     `json:"quantity"`
	UnitPrice float64   `json:"unitPrice"`
	BoughtAt  time.Time `json:"boughtAt"`
}

type InventoryLotGroup struct {
	ItemCode string          `json:"itemCode"`
	Lots     []*InventoryLot `json:"lots"`
}

type CountryInventory struct {
	UpdatedAt time.Time `json:"updatedAt"`
	CountryID string
	LotsMap   map[string][]InventoryLot
}

// ---------- processed: entity / battle ----------

type EntityWealthReport struct {
	ID          string    `json:"id"`
	DayStart    time.Time `json:"dayStart"`
	MemberCount int32     `json:"memberCount"`
	TotalDamage int64     `json:"totalDamage"`
	TotalWealth float64   `json:"totalWealth"`
	WagesPaid   float64   `json:"wagesPaid"`
	WagesEarned float64   `json:"wagesEarned"`
	EntityType  string
	EntityID    string
}

type EquipmentUsage struct {
	ItemCode string  `json:"itemCode"`
	Count    float64 `json:"count"`
	Value    float64 `json:"value"`
}

type BattleDamageReport struct {
	ID            string            `json:"id"`
	IntervalStart time.Time         `json:"intervalStart"`
	Side          string            `json:"side"`
	Damage        int64             `json:"damage"`
	DamagePct     float64           `json:"damagePct"`
	Equipment     []*EquipmentUsage `json:"equipment"`
	BattleID      string
	EntityType    string
	EntityID      string
}

// ---------- processed: market / items ----------

type ItemCandle struct {
	ID          string    `json:"id"`
	ItemCode    string    `json:"itemCode"`
	BucketStart time.Time `json:"bucketStart"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Avg         float64   `json:"avg"`
	Volume      int32     `json:"volume"`
	Money       float64   `json:"money"`
	Count       int32     `json:"count"`
}

type WageCandle struct {
	ID          string    `json:"id"`
	BucketStart time.Time `json:"bucketStart"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Avg         float64   `json:"avg"`
	Volume      int32     `json:"volume"`
	Money       float64   `json:"money"`
	Count       int32     `json:"count"`
}

type OrderbookLevel struct {
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type EffectivePrice struct {
	Size     int32   `json:"size"`
	AvgPrice float64 `json:"avgPrice"`
}

type ItemMarketReport struct {
	ItemCode       string            `json:"itemCode"`
	Volume24h      int32             `json:"volume24h"`
	AvgWeighted24h float64           `json:"avgWeighted24h"`
	PctChange24h   float64           `json:"pctChange24h"`
	Low24h         float64           `json:"low24h"`
	High24h        float64           `json:"high24h"`
	Spread         float64           `json:"spread"`
	Bids           []*OrderbookLevel `json:"bids"`
	Asks           []*OrderbookLevel `json:"asks"`
	EffectiveBuy   []*EffectivePrice `json:"effectiveBuy"`
	EffectiveSell  []*EffectivePrice `json:"effectiveSell"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type EquipmentWindowPrice struct {
	WeightedAvg float64   `json:"weightedAvg"`
	Volume      int32     `json:"volume"`
	Count       int32     `json:"count"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EquipmentSkillPrice struct {
	SkillKey  string        `json:"skillKey"`
	Skills    []*FloatEntry `json:"skills"`
	Min       float64       `json:"min"`
	Max       float64       `json:"max"`
	Avg       float64       `json:"avg"`
	Volume    int32         `json:"volume"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type EquipmentPricing struct {
	ItemCode   string                 `json:"itemCode"`
	WindowDays int32                  `json:"windowDays"`
	Window     *EquipmentWindowPrice  `json:"window"`
	Skills     []*EquipmentSkillPrice `json:"skills"`
}

type MarketState struct {
	At              time.Time `json:"at"`
	AvgWage24h      float64   `json:"avgWage24h"`
	WageVolume24h   int32     `json:"wageVolume24h"`
	MarketVolume24h float64   `json:"marketVolume24h"`
	WageMin         float64   `json:"wageMin"`
	WageMax         float64   `json:"wageMax"`
	WageAvgWeighted float64   `json:"wageAvgWeighted"`
}

type WagePaidUser struct {
	TotalPaid float64 `json:"totalPaid"`
	UserID    string
}

type WageMarketState struct {
	At             time.Time       `json:"at"`
	AvgWeighted14d float64         `json:"avgWeighted14d"`
	Min14d         float64         `json:"min14d"`
	Max14d         float64         `json:"max14d"`
	TotalPaid14d   float64         `json:"totalPaid14d"`
	Top10Least24h  []*WagePaidUser `json:"top10Least24h"`
	Top10Most24h   []*WagePaidUser `json:"top10Most24h"`
}

type InflationPoint struct {
	ID         string        `json:"id"`
	DayStart   time.Time     `json:"dayStart"`
	IndexValue float64       `json:"indexValue"`
	PctChange  float64       `json:"pctChange24h"`
	PerItem    []*FloatEntry `json:"perItem"`
}

type DismantleReport struct {
	ID           string      `json:"id"`
	HourStart    time.Time   `json:"hourStart"`
	Count        int32       `json:"count"`
	StateBuckets []*IntEntry `json:"stateBuckets"`
}

type CaseItemStat struct {
	ItemCode       string      `json:"itemCode"`
	AvgWeighted14d float64     `json:"avgWeighted14d"`
	TotalDrops     int32       `json:"totalDrops"`
	PerSkillRoll   []*IntEntry `json:"perSkillRoll"`
}

type CasesReport struct {
	Case             string          `json:"case"`
	TotalOpened      int32           `json:"totalOpened"`
	UniqueItemCodes  []string        `json:"uniqueItemCodes"`
	ExpectedValue14d float64         `json:"expectedValue14d"`
	PerItem          []*CaseItemStat `json:"perItem"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

// ---------- leaderboards / orderbook ----------

type DamageRanking struct {
	TotalDamage int64 `json:"totalDamage"`
	BattleCount int32 `json:"battleCount"`
	UserID      string
}

type WageRanking struct {
	Total  float64 `json:"total"`
	Count  int32   `json:"count"`
	UserID string
}

type OrderBookLevel struct {
	Price     float64 `json:"price"`
	Remaining int32   `json:"remaining"`
}

// OrderBook carries the item code so bids/asks are resolved on demand.
type OrderBook struct {
	ItemCode string `json:"itemCode"`
}
