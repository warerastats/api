package graph

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/warerastats/api/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ---------------- Battle ----------------

func (r *battleResolver) AttackerCountry(ctx context.Context, obj *model.Battle) (*model.Country, error) {
	return loadCountry(ctx, obj.AttackerCountryID)
}

func (r *battleResolver) DefenderCountry(ctx context.Context, obj *model.Battle) (*model.Country, error) {
	return loadCountry(ctx, obj.DefenderCountryID)
}

func (r *battleResolver) AttackerRegion(ctx context.Context, obj *model.Battle) (*model.Region, error) {
	return loadRegionP(ctx, obj.AttackerRegionID)
}

func (r *battleResolver) DefenderRegion(ctx context.Context, obj *model.Battle) (*model.Region, error) {
	return loadRegion(ctx, obj.DefenderRegionID)
}

func (r *battleResolver) Damages(ctx context.Context, obj *model.Battle, first *int32, after *string, side *string, userID *string) ([]*model.Damage, error) {
	bid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	uid, err := oidArgPtr(userID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Damage.GetByBattlePaged(ctx, bid, before, limitOf(first, 50), sidePtr(side), uid)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toDamage), nil
}

func (r *battleResolver) Mus(ctx context.Context, obj *model.Battle) ([]*model.Mu, error) {
	bid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	orders, err := r.Colls.Events.BattleOrderChange.MuOrdersByBattles(ctx, []bson.ObjectID{bid})
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(orders))
	out := make([]*model.Mu, 0, len(orders))
	for _, o := range orders {
		hex := o.MuID.Hex()
		if _, dup := seen[hex]; dup {
			continue
		}
		seen[hex] = struct{}{}
		m, err := loadMu(ctx, hex)
		if err != nil {
			return nil, err
		}
		if m != nil {
			out = append(out, m)
		}
	}
	return out, nil
}

func (r *battleResolver) OrderChanges(ctx context.Context, obj *model.Battle, first *int32, after *string) ([]*model.BattleOrderChange, error) {
	bid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.BattleOrderChange.ListByBattle(ctx, bid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toBattleOrderChange), nil
}

func (r *battleResolver) DamageReports(ctx context.Context, obj *model.Battle, from time.Time, to time.Time) ([]*model.BattleDamageReport, error) {
	bid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.BattleDamageReport.GetByBattle(ctx, bid, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toBattleDamageReport), nil
}

// ---------------- Company ----------------

func (r *companyResolver) Owner(ctx context.Context, obj *model.Company) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *companyResolver) Region(ctx context.Context, obj *model.Company) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *companyResolver) Employees(ctx context.Context, obj *model.Company) ([]*model.Employee, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Employee.GetByCompany(ctx, cid)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toEmployee), nil
}

func (r *companyResolver) RegionHistory(ctx context.Context, obj *model.Company, first *int32, after *string) ([]*model.CompanyRegionChange, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.CompanyRegionChange.ListByCompany(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCompanyRegionChange), nil
}

func (r *companyResolver) ItemCodeHistory(ctx context.Context, obj *model.Company, first *int32, after *string) ([]*model.CompanyItemCodeChange, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.CompanyItemCodeChange.ListByCompany(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCompanyItemCodeChange), nil
}

// ---------------- Country ----------------

func (r *countryResolver) RulingParty(ctx context.Context, obj *model.Country) (*model.Party, error) {
	return loadPartyP(ctx, obj.RulingPartyID)
}

func (r *countryResolver) Users(ctx context.Context, obj *model.Country, first *int32, after *string) ([]*model.User, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.User.GetByCountryPaged(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toUser), nil
}

func (r *countryResolver) UserCount(ctx context.Context, obj *model.Country) (int32, error) {
	aggs, err := r.Colls.Trackers.User.CountByCountry(ctx)
	if err != nil {
		return 0, err
	}
	for _, a := range aggs {
		if a.CountryID.Hex() == obj.ID {
			return int32(a.MemberCount), nil
		}
	}
	return 0, nil
}

func (r *countryResolver) Parties(ctx context.Context, obj *model.Country, first *int32, after *string) ([]*model.Party, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Party.GetByCountryPaged(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toParty), nil
}

func (r *countryResolver) Regions(ctx context.Context, obj *model.Country) ([]*model.Region, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Region.GetByCountry(ctx, cid)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toRegion), nil
}

func (r *countryResolver) Battles(ctx context.Context, obj *model.Country, first *int32, after *string) ([]*model.Battle, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Battle.GetByCountryPaged(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toBattle), nil
}

func (r *countryResolver) RulingPartyHistory(ctx context.Context, obj *model.Country, first *int32, after *string) ([]*model.CountryRulingPartyChange, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.CountryRulingPartyChange.ListByCountry(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCountryRulingPartyChange), nil
}

func (r *countryResolver) SpecialisationHistory(ctx context.Context, obj *model.Country, first *int32, after *string) ([]*model.CountrySpecialisationChange, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.CountrySpecialisationChange.ListByCountry(ctx, cid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCountrySpecialisationChange), nil
}

func (r *countryResolver) TaxFlows(ctx context.Context, obj *model.Country, from time.Time, to time.Time) ([]*model.CountryTaxFlow, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.CountryTaxFlow.GetByCountryRange(ctx, cid, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCountryTaxFlow), nil
}

func (r *countryResolver) FlipEvents(ctx context.Context, obj *model.Country, from time.Time, to time.Time, first *int32, after *string) ([]*model.CountryFlipEvent, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Estimators.CountryFlipEvent.GetByCountryRange(ctx, cid, from, to, before, limitOf(first, 50))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCountryFlipEvent), nil
}

func (r *countryResolver) FlipState(ctx context.Context, obj *model.Country) (*model.CountryFlipState, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	st, ok, err := r.Colls.Processed.Estimators.CountryFlipState.Get(ctx, cid)
	if err != nil || !ok {
		return nil, err
	}
	return toCountryFlipState(st), nil
}

func (r *countryResolver) Inventory(ctx context.Context, obj *model.Country) (*model.CountryInventory, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	inv, ok, err := r.Colls.Processed.Estimators.CountryInventory.Get(ctx, cid)
	if err != nil || !ok {
		return nil, err
	}
	return toCountryInventory(inv), nil
}

func (r *countryResolver) WealthReports(ctx context.Context, obj *model.Country, from time.Time, to time.Time) ([]*model.EntityWealthReport, error) {
	cid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.EntityWealthReport.GetByEntityRange(ctx, "country", cid, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toEntityWealthReport), nil
}

// ---------------- Damage ----------------

func (r *damageResolver) Battle(ctx context.Context, obj *model.Damage) (*model.Battle, error) {
	return loadBattle(ctx, obj.BattleID)
}

func (r *damageResolver) User(ctx context.Context, obj *model.Damage) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *damageResolver) Country(ctx context.Context, obj *model.Damage) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *damageResolver) Mu(ctx context.Context, obj *model.Damage) (*model.Mu, error) {
	return loadMuP(ctx, obj.MuID)
}

func (r *damageResolver) Party(ctx context.Context, obj *model.Damage) (*model.Party, error) {
	return loadPartyP(ctx, obj.PartyID)
}

// Skill is not resolved: the skills collection has no by-id lookup; the snapshot
// id is exposed only as a relation placeholder. Returns null.
func (r *damageResolver) Skill(ctx context.Context, obj *model.Damage) (*model.Skill, error) {
	return nil, nil
}

func (r *damageResolver) Weapon(ctx context.Context, obj *model.Damage) (*model.Item, error) {
	return loadItemP(ctx, obj.WeaponID)
}

func (r *damageResolver) Helmet(ctx context.Context, obj *model.Damage) (*model.Item, error) {
	return loadItemP(ctx, obj.HelmetID)
}

func (r *damageResolver) Chest(ctx context.Context, obj *model.Damage) (*model.Item, error) {
	return loadItemP(ctx, obj.ChestID)
}

func (r *damageResolver) Pants(ctx context.Context, obj *model.Damage) (*model.Item, error) {
	return loadItemP(ctx, obj.PantsID)
}

func (r *damageResolver) Boots(ctx context.Context, obj *model.Damage) (*model.Item, error) {
	return loadItemP(ctx, obj.BootsID)
}

func (r *damageResolver) Gloves(ctx context.Context, obj *model.Damage) (*model.Item, error) {
	return loadItemP(ctx, obj.GlovesID)
}

// ---------------- Employee ----------------

func (r *employeeResolver) User(ctx context.Context, obj *model.Employee) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *employeeResolver) Company(ctx context.Context, obj *model.Employee) (*model.Company, error) {
	return loadCompany(ctx, obj.CompanyID)
}

func (r *employeeResolver) Employer(ctx context.Context, obj *model.Employee) (*model.User, error) {
	return loadUser(ctx, obj.EmployerID)
}

func (r *employeeResolver) WageHistory(ctx context.Context, obj *model.Employee, first *int32, after *string) ([]*model.EmployeeWageChange, error) {
	uid, err := oidOf(obj.UserID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.EmployeeWageChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toEmployeeWageChange), nil
}

// ---------------- Item ----------------

func (r *itemResolver) Skills(ctx context.Context, obj *model.Item) ([]*model.FloatEntry, error) {
	return floatEntries(obj.SkillsMap), nil
}

func (r *itemResolver) Owner(ctx context.Context, obj *model.Item) (*model.User, error) {
	return loadUser(ctx, obj.OwnerUserID)
}

func (r *itemResolver) MarketTransactions(ctx context.Context, obj *model.Item, first *int32, after *string) ([]*model.MarketTransaction, error) {
	iid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Transactions.MarketTransaction.GetByItemPaged(ctx, iid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toMarket), nil
}

// ---------------- Mu ----------------

func (r *muResolver) Owner(ctx context.Context, obj *model.Mu) (*model.User, error) {
	return loadUser(ctx, obj.OwnerUserID)
}

func (r *muResolver) Region(ctx context.Context, obj *model.Mu) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *muResolver) Members(ctx context.Context, obj *model.Mu) ([]*model.User, error) {
	return r.loadMembers(ctx, obj.MemberUserIDs)
}

func (r *muResolver) Battles(ctx context.Context, obj *model.Mu, first *int32, after *string) ([]*model.Battle, error) {
	mid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	ids, err := r.Colls.Events.BattleOrderChange.ListBattleIDsByMu(ctx, mid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	out := make([]*model.Battle, 0, len(ids))
	for _, id := range ids {
		b, err := loadBattle(ctx, id.Hex())
		if err != nil {
			return nil, err
		}
		if b != nil {
			out = append(out, b)
		}
	}
	return out, nil
}

func (r *muResolver) NameHistory(ctx context.Context, obj *model.Mu, first *int32, after *string) ([]*model.MuNameChange, error) {
	mid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.MuNameChange.ListByMu(ctx, mid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toMuNameChange), nil
}

func (r *muResolver) OwnerHistory(ctx context.Context, obj *model.Mu, first *int32, after *string) ([]*model.MuOwnerChange, error) {
	mid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.MuOwnerChange.ListByMu(ctx, mid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toMuOwnerChange), nil
}

func (r *muResolver) MercReputationHistory(ctx context.Context, obj *model.Mu, first *int32, after *string) ([]*model.MuMercenaryReputationChange, error) {
	mid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.MuMercenaryReputationChange.ListByMu(ctx, mid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toMuMercenaryReputationChange), nil
}

func (r *muResolver) WealthReports(ctx context.Context, obj *model.Mu, from time.Time, to time.Time) ([]*model.EntityWealthReport, error) {
	mid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.EntityWealthReport.GetByEntityRange(ctx, "mu", mid, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toEntityWealthReport), nil
}

// ---------------- Party ----------------

func (r *partyResolver) Country(ctx context.Context, obj *model.Party) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *partyResolver) Region(ctx context.Context, obj *model.Party) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *partyResolver) Leader(ctx context.Context, obj *model.Party) (*model.User, error) {
	return loadUser(ctx, obj.LeaderUserID)
}

func (r *partyResolver) Members(ctx context.Context, obj *model.Party) ([]*model.User, error) {
	return r.loadMembers(ctx, obj.MemberUserIDs)
}

func (r *partyResolver) RulesCountries(ctx context.Context, obj *model.Party) ([]*model.Country, error) {
	pid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Country.GetByRulingParty(ctx, pid)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toCountry), nil
}

func (r *partyResolver) NameHistory(ctx context.Context, obj *model.Party, first *int32, after *string) ([]*model.PartyNameChange, error) {
	pid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.PartyNameChange.ListByParty(ctx, pid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toPartyNameChange), nil
}

func (r *partyResolver) LeaderHistory(ctx context.Context, obj *model.Party, first *int32, after *string) ([]*model.PartyLeaderChange, error) {
	pid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.PartyLeaderChange.ListByParty(ctx, pid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toPartyLeaderChange), nil
}

func (r *partyResolver) DescriptionHistory(ctx context.Context, obj *model.Party, first *int32, after *string) ([]*model.PartyDescriptionChange, error) {
	pid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.PartyDescriptionChange.ListByParty(ctx, pid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toPartyDescriptionChange), nil
}

func (r *partyResolver) EthicsHistory(ctx context.Context, obj *model.Party, first *int32, after *string) ([]*model.PartyEthicsChange, error) {
	pid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.PartyEthicsChange.ListByParty(ctx, pid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toPartyEthicsChange), nil
}

func (r *partyResolver) WealthReports(ctx context.Context, obj *model.Party, from time.Time, to time.Time) ([]*model.EntityWealthReport, error) {
	pid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.EntityWealthReport.GetByEntityRange(ctx, "party", pid, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toEntityWealthReport), nil
}

// ---------------- Region ----------------

func (r *regionResolver) Country(ctx context.Context, obj *model.Region) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *regionResolver) InitialCountry(ctx context.Context, obj *model.Region) (*model.Country, error) {
	return loadCountry(ctx, obj.InitialCountryID)
}

func (r *regionResolver) Neighbors(ctx context.Context, obj *model.Region) ([]*model.Region, error) {
	ids := make([]bson.ObjectID, 0, len(obj.NeighborIDs))
	for _, h := range obj.NeighborIDs {
		id, err := oidOf(h)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	rows, err := r.Colls.Trackers.Region.GetMany(ctx, ids)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toRegion), nil
}

func (r *regionResolver) Parties(ctx context.Context, obj *model.Region) ([]*model.Party, error) {
	rid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Party.GetByRegion(ctx, rid)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toParty), nil
}

func (r *regionResolver) Mus(ctx context.Context, obj *model.Region) ([]*model.Mu, error) {
	rid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Mu.GetByRegion(ctx, rid)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toMu), nil
}

func (r *regionResolver) Companies(ctx context.Context, obj *model.Region, first *int32, after *string) ([]*model.Company, error) {
	rid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Company.GetByRegionPaged(ctx, rid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toCompany), nil
}

func (r *regionResolver) OwnerHistory(ctx context.Context, obj *model.Region, first *int32, after *string) ([]*model.RegionOwnerChange, error) {
	rid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.RegionOwnerChange.ListByRegion(ctx, rid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toRegionOwnerChange), nil
}

func (r *regionResolver) Deposits(ctx context.Context, obj *model.Region, first *int32, after *string) ([]*model.RegionDeposit, error) {
	rid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.RegionDeposit.ListByRegion(ctx, rid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toRegionDeposit), nil
}

func (r *regionResolver) StrategicResources(ctx context.Context, obj *model.Region, first *int32, after *string) ([]*model.RegionStrategicResource, error) {
	rid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.RegionStrategicResource.ListByRegion(ctx, rid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toRegionStrategicResource), nil
}

// ---------------- Skill ----------------

func (r *skillResolver) User(ctx context.Context, obj *model.Skill) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

// ---------------- TradeOffer ----------------

func (r *tradeOfferResolver) User(ctx context.Context, obj *model.TradeOffer) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *tradeOfferResolver) Country(ctx context.Context, obj *model.TradeOffer) (*model.Country, error) {
	return loadCountryP(ctx, obj.CountryID)
}

func (r *tradeOfferResolver) Mu(ctx context.Context, obj *model.TradeOffer) (*model.Mu, error) {
	return loadMuP(ctx, obj.MuID)
}

// ---------------- User ----------------

func (r *userResolver) Wealth(ctx context.Context, obj *model.User) ([]*model.FloatEntry, error) {
	return floatEntries(obj.WealthMap), nil
}

func (r *userResolver) Skills(ctx context.Context, obj *model.User) ([]*model.IntEntry, error) {
	return intEntries(obj.SkillsMap), nil
}

func (r *userResolver) Country(ctx context.Context, obj *model.User) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *userResolver) Company(ctx context.Context, obj *model.User) (*model.Company, error) {
	return loadCompanyP(ctx, obj.CompanyID)
}

func (r *userResolver) Party(ctx context.Context, obj *model.User) (*model.Party, error) {
	return loadPartyP(ctx, obj.PartyID)
}

func (r *userResolver) Mu(ctx context.Context, obj *model.User) (*model.Mu, error) {
	return loadMuP(ctx, obj.MuID)
}

func (r *userResolver) Transactions(ctx context.Context, obj *model.User, first *int32, after *string) (*model.ActivityConnection, error) {
	limit := limitOf(first, 20)
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}

	type edge struct {
		id  string
		act model.Activity
	}
	var edges []edge

	markets, err := r.Colls.Transactions.MarketTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range markets {
		edges = append(edges, edge{t.ID.Hex(), toMarket(t)})
	}
	wages, err := r.Colls.Transactions.WageTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range wages {
		edges = append(edges, edge{t.ID.Hex(), toWage(t)})
	}
	cases, err := r.Colls.Transactions.CaseTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range cases {
		edges = append(edges, edge{t.ID.Hex(), toCase(t)})
	}
	trades, err := r.Colls.Transactions.TradeTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range trades {
		edges = append(edges, edge{t.ID.Hex(), toTrade(t)})
	}
	crafts, err := r.Colls.Transactions.CraftTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range crafts {
		edges = append(edges, edge{t.ID.Hex(), toCraft(t)})
	}
	dismantles, err := r.Colls.Transactions.DismantleTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range dismantles {
		edges = append(edges, edge{t.ID.Hex(), toDismantle(t)})
	}
	loots, err := r.Colls.Transactions.LootTransaction.ByUser(ctx, uid, before, limit)
	if err != nil {
		return nil, err
	}
	for _, t := range loots {
		edges = append(edges, edge{t.ID.Hex(), toLoot(t)})
	}

	// ObjectID hex sorts lexicographically in timestamp order, so a descending
	// hex sort is a newest-first time sort across the merged collections.
	sort.Slice(edges, func(i, j int) bool { return edges[i].id > edges[j].id })

	more := len(edges) > limit
	if more {
		edges = edges[:limit]
	}
	conn := &model.ActivityConnection{
		Edges:       make([]model.Activity, len(edges)),
		HasNextPage: more,
	}
	for i, e := range edges {
		conn.Edges[i] = e.act
	}
	if len(edges) > 0 {
		cur := edges[len(edges)-1].id
		conn.EndCursor = &cur
	}
	return conn, nil
}

func (r *userResolver) Damages(ctx context.Context, obj *model.User, battleID string) ([]*model.Damage, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	bid, err := oidOf(battleID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Damage.ByUserAndBattle(ctx, uid, bid)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toDamage), nil
}

func (r *userResolver) AllDamages(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.Damage, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Damage.GetByUserPaged(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toDamage), nil
}

func (r *userResolver) Items(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.Item, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Item.GetByOwnerPaged(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toItem), nil
}

func (r *userResolver) OwnedCompanies(ctx context.Context, obj *model.User) ([]*model.Company, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Company.GetByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toCompany), nil
}

func (r *userResolver) Employment(ctx context.Context, obj *model.User) (*model.Employee, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	e, ok, err := r.Colls.Trackers.Employee.GetByUserID(ctx, uid)
	if err != nil || !ok {
		return nil, err
	}
	return toEmployee(e), nil
}

func (r *userResolver) TradeOffers(ctx context.Context, obj *model.User, first *int32, after *string, itemCode *string, side *model.TradeSide) ([]*model.TradeOffer, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.TradeOffer.GetByUserPaged(ctx, uid, before, limitOf(first, 20), itemCode, tradeSidePtr(side))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toTradeOffer), nil
}

func (r *userResolver) LatestSkills(ctx context.Context, obj *model.User) (*model.Skill, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	s, err := r.Colls.Trackers.Skill.GetLatestForUser(ctx, uid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return toSkill(s), nil
}

func (r *userResolver) SkillSnapshots(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.Skill, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Skill.GetHistoryForUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toSkill), nil
}

func (r *userResolver) NameHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.UserNameChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.UserNameChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserNameChange), nil
}

func (r *userResolver) CountryHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.UserCountryChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.UserCountryChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserCountryChange), nil
}

func (r *userResolver) PartyHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.UserPartyChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.UserPartyChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserPartyChange), nil
}

func (r *userResolver) CompanyHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.UserCompanyChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.UserCompanyChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserCompanyChange), nil
}

func (r *userResolver) MuHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.UserMuChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.UserMUChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserMuChange), nil
}

func (r *userResolver) SkillChangeHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.UserSkillChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.UserSkillChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserSkillChange), nil
}

func (r *userResolver) WageHistory(ctx context.Context, obj *model.User, first *int32, after *string) ([]*model.EmployeeWageChange, error) {
	uid, before, err := r.userCursor(obj, after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Events.EmployeeWageChange.ListByUser(ctx, uid, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toEmployeeWageChange), nil
}

func (r *userResolver) FinanceReports(ctx context.Context, obj *model.User, from time.Time, to time.Time) ([]*model.UserFinanceReport, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.UserFinanceReport.GetByUserRange(ctx, uid, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserFinanceReport), nil
}

func (r *userResolver) FlipEvents(ctx context.Context, obj *model.User, from time.Time, to time.Time, first *int32, after *string) ([]*model.UserFlipEvent, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Estimators.UserFlipEvent.GetByUserRange(ctx, uid, from, to, before, limitOf(first, 50))
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toUserFlipEvent), nil
}

func (r *userResolver) FlipState(ctx context.Context, obj *model.User) (*model.UserFlipState, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	st, ok, err := r.Colls.Processed.Estimators.UserFlipState.Get(ctx, uid)
	if err != nil || !ok {
		return nil, err
	}
	return toUserFlipState(st), nil
}

func (r *userResolver) Inventory(ctx context.Context, obj *model.User) (*model.UserInventory, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	inv, ok, err := r.Colls.Processed.Estimators.UserInventory.Get(ctx, uid)
	if err != nil || !ok {
		return nil, err
	}
	return toUserInventory(inv), nil
}

func (r *userResolver) BattleParticipation(ctx context.Context, obj *model.User) (*model.UserBattleParticipation, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return nil, err
	}
	p, ok, err := r.Colls.Processed.Estimators.BattleParticipation.Get(ctx, uid)
	if err != nil || !ok {
		return nil, err
	}
	return toUserBattleParticipation(p), nil
}

// userCursor parses a user object id and an optional after cursor for history lists.
func (r *userResolver) userCursor(obj *model.User, after *string) (bson.ObjectID, *bson.ObjectID, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return bson.ObjectID{}, nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return bson.ObjectID{}, nil, err
	}
	return uid, before, nil
}

// Battle returns BattleResolver implementation.
func (r *Resolver) Battle() BattleResolver { return &battleResolver{r} }

// Company returns CompanyResolver implementation.
func (r *Resolver) Company() CompanyResolver { return &companyResolver{r} }

// Country returns CountryResolver implementation.
func (r *Resolver) Country() CountryResolver { return &countryResolver{r} }

// Damage returns DamageResolver implementation.
func (r *Resolver) Damage() DamageResolver { return &damageResolver{r} }

// Employee returns EmployeeResolver implementation.
func (r *Resolver) Employee() EmployeeResolver { return &employeeResolver{r} }

// Item returns ItemResolver implementation.
func (r *Resolver) Item() ItemResolver { return &itemResolver{r} }

// Mu returns MuResolver implementation.
func (r *Resolver) Mu() MuResolver { return &muResolver{r} }

// Party returns PartyResolver implementation.
func (r *Resolver) Party() PartyResolver { return &partyResolver{r} }

// Region returns RegionResolver implementation.
func (r *Resolver) Region() RegionResolver { return &regionResolver{r} }

// Skill returns SkillResolver implementation.
func (r *Resolver) Skill() SkillResolver { return &skillResolver{r} }

// TradeOffer returns TradeOfferResolver implementation.
func (r *Resolver) TradeOffer() TradeOfferResolver { return &tradeOfferResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type battleResolver struct{ *Resolver }
type companyResolver struct{ *Resolver }
type countryResolver struct{ *Resolver }
type damageResolver struct{ *Resolver }
type employeeResolver struct{ *Resolver }
type itemResolver struct{ *Resolver }
type muResolver struct{ *Resolver }
type partyResolver struct{ *Resolver }
type regionResolver struct{ *Resolver }
type skillResolver struct{ *Resolver }
type tradeOfferResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
