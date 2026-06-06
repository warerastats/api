package graph

import (
	"context"

	"github.com/warerastats/api/graph/model"
)

func (r *battleOrderChangeResolver) Battle(ctx context.Context, obj *model.BattleOrderChange) (*model.Battle, error) {
	return loadBattle(ctx, obj.BattleID)
}

func (r *battleOrderChangeResolver) Entity(ctx context.Context, obj *model.BattleOrderChange) (model.Entity, error) {
	return loadEntity(ctx, obj.Kind, obj.EntityID)
}

func (r *companyItemCodeChangeResolver) Company(ctx context.Context, obj *model.CompanyItemCodeChange) (*model.Company, error) {
	return loadCompany(ctx, obj.CompanyID)
}

func (r *companyRegionChangeResolver) Company(ctx context.Context, obj *model.CompanyRegionChange) (*model.Company, error) {
	return loadCompany(ctx, obj.CompanyID)
}

func (r *companyRegionChangeResolver) Region(ctx context.Context, obj *model.CompanyRegionChange) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *countryRulingPartyChangeResolver) Country(ctx context.Context, obj *model.CountryRulingPartyChange) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *countryRulingPartyChangeResolver) Party(ctx context.Context, obj *model.CountryRulingPartyChange) (*model.Party, error) {
	return loadPartyP(ctx, obj.PartyID)
}

func (r *countrySpecialisationChangeResolver) Country(ctx context.Context, obj *model.CountrySpecialisationChange) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *employeeWageChangeResolver) User(ctx context.Context, obj *model.EmployeeWageChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *employeeWageChangeResolver) Company(ctx context.Context, obj *model.EmployeeWageChange) (*model.Company, error) {
	return loadCompany(ctx, obj.CompanyID)
}

func (r *muMercenaryReputationChangeResolver) Mu(ctx context.Context, obj *model.MuMercenaryReputationChange) (*model.Mu, error) {
	return loadMu(ctx, obj.MuID)
}

func (r *muNameChangeResolver) Mu(ctx context.Context, obj *model.MuNameChange) (*model.Mu, error) {
	return loadMu(ctx, obj.MuID)
}

func (r *muOwnerChangeResolver) Mu(ctx context.Context, obj *model.MuOwnerChange) (*model.Mu, error) {
	return loadMu(ctx, obj.MuID)
}

func (r *muOwnerChangeResolver) Owner(ctx context.Context, obj *model.MuOwnerChange) (*model.User, error) {
	return loadUser(ctx, obj.OwnerUserID)
}

func (r *partyDescriptionChangeResolver) Party(ctx context.Context, obj *model.PartyDescriptionChange) (*model.Party, error) {
	return loadParty(ctx, obj.PartyID)
}

func (r *partyEthicsChangeResolver) Party(ctx context.Context, obj *model.PartyEthicsChange) (*model.Party, error) {
	return loadParty(ctx, obj.PartyID)
}

func (r *partyLeaderChangeResolver) Party(ctx context.Context, obj *model.PartyLeaderChange) (*model.Party, error) {
	return loadParty(ctx, obj.PartyID)
}

func (r *partyLeaderChangeResolver) Leader(ctx context.Context, obj *model.PartyLeaderChange) (*model.User, error) {
	return loadUser(ctx, obj.LeaderUserID)
}

func (r *partyNameChangeResolver) Party(ctx context.Context, obj *model.PartyNameChange) (*model.Party, error) {
	return loadParty(ctx, obj.PartyID)
}

func (r *regionDepositResolver) Region(ctx context.Context, obj *model.RegionDeposit) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *regionOwnerChangeResolver) Region(ctx context.Context, obj *model.RegionOwnerChange) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *regionOwnerChangeResolver) Country(ctx context.Context, obj *model.RegionOwnerChange) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *regionStrategicResourceResolver) Region(ctx context.Context, obj *model.RegionStrategicResource) (*model.Region, error) {
	return loadRegion(ctx, obj.RegionID)
}

func (r *userCompanyChangeResolver) User(ctx context.Context, obj *model.UserCompanyChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userCompanyChangeResolver) Company(ctx context.Context, obj *model.UserCompanyChange) (*model.Company, error) {
	return loadCompanyP(ctx, obj.CompanyID)
}

func (r *userCountryChangeResolver) User(ctx context.Context, obj *model.UserCountryChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userCountryChangeResolver) Country(ctx context.Context, obj *model.UserCountryChange) (*model.Country, error) {
	return loadCountryP(ctx, obj.CountryID)
}

func (r *userMuChangeResolver) User(ctx context.Context, obj *model.UserMuChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userMuChangeResolver) Mu(ctx context.Context, obj *model.UserMuChange) (*model.Mu, error) {
	return loadMuP(ctx, obj.MuID)
}

func (r *userNameChangeResolver) User(ctx context.Context, obj *model.UserNameChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userPartyChangeResolver) User(ctx context.Context, obj *model.UserPartyChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userPartyChangeResolver) Party(ctx context.Context, obj *model.UserPartyChange) (*model.Party, error) {
	return loadPartyP(ctx, obj.PartyID)
}

func (r *userSkillChangeResolver) User(ctx context.Context, obj *model.UserSkillChange) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userSkillChangeResolver) Skills(ctx context.Context, obj *model.UserSkillChange) ([]*model.IntEntry, error) {
	return intEntries(obj.SkillsMap), nil
}

// BattleOrderChange returns BattleOrderChangeResolver implementation.
func (r *Resolver) BattleOrderChange() BattleOrderChangeResolver {
	return &battleOrderChangeResolver{r}
}

// CompanyItemCodeChange returns CompanyItemCodeChangeResolver implementation.
func (r *Resolver) CompanyItemCodeChange() CompanyItemCodeChangeResolver {
	return &companyItemCodeChangeResolver{r}
}

// CompanyRegionChange returns CompanyRegionChangeResolver implementation.
func (r *Resolver) CompanyRegionChange() CompanyRegionChangeResolver {
	return &companyRegionChangeResolver{r}
}

// CountryRulingPartyChange returns CountryRulingPartyChangeResolver implementation.
func (r *Resolver) CountryRulingPartyChange() CountryRulingPartyChangeResolver {
	return &countryRulingPartyChangeResolver{r}
}

// CountrySpecialisationChange returns CountrySpecialisationChangeResolver implementation.
func (r *Resolver) CountrySpecialisationChange() CountrySpecialisationChangeResolver {
	return &countrySpecialisationChangeResolver{r}
}

// EmployeeWageChange returns EmployeeWageChangeResolver implementation.
func (r *Resolver) EmployeeWageChange() EmployeeWageChangeResolver {
	return &employeeWageChangeResolver{r}
}

// MuMercenaryReputationChange returns MuMercenaryReputationChangeResolver implementation.
func (r *Resolver) MuMercenaryReputationChange() MuMercenaryReputationChangeResolver {
	return &muMercenaryReputationChangeResolver{r}
}

// MuNameChange returns MuNameChangeResolver implementation.
func (r *Resolver) MuNameChange() MuNameChangeResolver { return &muNameChangeResolver{r} }

// MuOwnerChange returns MuOwnerChangeResolver implementation.
func (r *Resolver) MuOwnerChange() MuOwnerChangeResolver { return &muOwnerChangeResolver{r} }

// PartyDescriptionChange returns PartyDescriptionChangeResolver implementation.
func (r *Resolver) PartyDescriptionChange() PartyDescriptionChangeResolver {
	return &partyDescriptionChangeResolver{r}
}

// PartyEthicsChange returns PartyEthicsChangeResolver implementation.
func (r *Resolver) PartyEthicsChange() PartyEthicsChangeResolver {
	return &partyEthicsChangeResolver{r}
}

// PartyLeaderChange returns PartyLeaderChangeResolver implementation.
func (r *Resolver) PartyLeaderChange() PartyLeaderChangeResolver {
	return &partyLeaderChangeResolver{r}
}

// PartyNameChange returns PartyNameChangeResolver implementation.
func (r *Resolver) PartyNameChange() PartyNameChangeResolver { return &partyNameChangeResolver{r} }

// RegionDeposit returns RegionDepositResolver implementation.
func (r *Resolver) RegionDeposit() RegionDepositResolver { return &regionDepositResolver{r} }

// RegionOwnerChange returns RegionOwnerChangeResolver implementation.
func (r *Resolver) RegionOwnerChange() RegionOwnerChangeResolver {
	return &regionOwnerChangeResolver{r}
}

// RegionStrategicResource returns RegionStrategicResourceResolver implementation.
func (r *Resolver) RegionStrategicResource() RegionStrategicResourceResolver {
	return &regionStrategicResourceResolver{r}
}

// UserCompanyChange returns UserCompanyChangeResolver implementation.
func (r *Resolver) UserCompanyChange() UserCompanyChangeResolver {
	return &userCompanyChangeResolver{r}
}

// UserCountryChange returns UserCountryChangeResolver implementation.
func (r *Resolver) UserCountryChange() UserCountryChangeResolver {
	return &userCountryChangeResolver{r}
}

// UserMuChange returns UserMuChangeResolver implementation.
func (r *Resolver) UserMuChange() UserMuChangeResolver { return &userMuChangeResolver{r} }

// UserNameChange returns UserNameChangeResolver implementation.
func (r *Resolver) UserNameChange() UserNameChangeResolver { return &userNameChangeResolver{r} }

// UserPartyChange returns UserPartyChangeResolver implementation.
func (r *Resolver) UserPartyChange() UserPartyChangeResolver { return &userPartyChangeResolver{r} }

// UserSkillChange returns UserSkillChangeResolver implementation.
func (r *Resolver) UserSkillChange() UserSkillChangeResolver { return &userSkillChangeResolver{r} }

type battleOrderChangeResolver struct{ *Resolver }
type companyItemCodeChangeResolver struct{ *Resolver }
type companyRegionChangeResolver struct{ *Resolver }
type countryRulingPartyChangeResolver struct{ *Resolver }
type countrySpecialisationChangeResolver struct{ *Resolver }
type employeeWageChangeResolver struct{ *Resolver }
type muMercenaryReputationChangeResolver struct{ *Resolver }
type muNameChangeResolver struct{ *Resolver }
type muOwnerChangeResolver struct{ *Resolver }
type partyDescriptionChangeResolver struct{ *Resolver }
type partyEthicsChangeResolver struct{ *Resolver }
type partyLeaderChangeResolver struct{ *Resolver }
type partyNameChangeResolver struct{ *Resolver }
type regionDepositResolver struct{ *Resolver }
type regionOwnerChangeResolver struct{ *Resolver }
type regionStrategicResourceResolver struct{ *Resolver }
type userCompanyChangeResolver struct{ *Resolver }
type userCountryChangeResolver struct{ *Resolver }
type userMuChangeResolver struct{ *Resolver }
type userNameChangeResolver struct{ *Resolver }
type userPartyChangeResolver struct{ *Resolver }
type userSkillChangeResolver struct{ *Resolver }
