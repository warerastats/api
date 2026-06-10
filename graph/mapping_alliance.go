package graph

import (
	"github.com/warerastats/api/graph/model"
	processedreports "github.com/warerastats/models/models/stores/processed/reports"
)

func toCountryAllianceMoneyFlowReport(r processedreports.CountryAllianceMoneyFlowReport) *model.CountryAllianceMoneyFlowReport {
	counterparts := make([]*model.CountryAllianceMoneyFlowCounterpart, len(r.Counterparts))
	for i, cp := range r.Counterparts {
		counterparts[i] = &model.CountryAllianceMoneyFlowCounterpart{
			InEquipment:  cp.InEquipment,
			OutEquipment: cp.OutEquipment,
			InItems:      cp.InItems,
			OutItems:     cp.OutItems,
			InWages:      cp.InWages,
			OutWages:     cp.OutWages,
			AllianceID:   cp.AllianceID.Hex(),
		}
	}
	return &model.CountryAllianceMoneyFlowReport{
		ID:                          r.ID,
		DayStart:                    r.DayStart,
		InEquipment:                 r.InEquipment,
		OutEquipment:                r.OutEquipment,
		InItems:                     r.InItems,
		OutItems:                    r.OutItems,
		InWages:                     r.InWages,
		OutWages:                    r.OutWages,
		InEquipmentInAlliance:       r.InEquipmentInAlliance,
		OutEquipmentInAlliance:      r.OutEquipmentInAlliance,
		InItemsInAlliance:           r.InItemsInAlliance,
		OutItemsInAlliance:          r.OutItemsInAlliance,
		InWagesInAlliance:           r.InWagesInAlliance,
		OutWagesInAlliance:          r.OutWagesInAlliance,
		InEquipmentOutsideAlliance:  r.InEquipmentOutsideAlliance,
		OutEquipmentOutsideAlliance: r.OutEquipmentOutsideAlliance,
		InItemsOutsideAlliance:      r.InItemsOutsideAlliance,
		OutItemsOutsideAlliance:     r.OutItemsOutsideAlliance,
		InWagesOutsideAlliance:      r.InWagesOutsideAlliance,
		OutWagesOutsideAlliance:     r.OutWagesOutsideAlliance,
		Counterparts:                counterparts,
		CountryID:                   r.CountryID.Hex(),
	}
}

func toAllianceMoneyFlowReport(r processedreports.AllianceMoneyFlowReport) *model.AllianceMoneyFlowReport {
	counterparts := make([]*model.AllianceMoneyFlowCounterpart, len(r.Counterparts))
	for i, cp := range r.Counterparts {
		counterparts[i] = &model.AllianceMoneyFlowCounterpart{
			InEquipment:  cp.InEquipment,
			OutEquipment: cp.OutEquipment,
			InItems:      cp.InItems,
			OutItems:     cp.OutItems,
			InWages:      cp.InWages,
			OutWages:     cp.OutWages,
			AllianceID:   cp.AllianceID.Hex(),
		}
	}
	return &model.AllianceMoneyFlowReport{
		ID:                          r.ID,
		DayStart:                    r.DayStart,
		InEquipment:                 r.InEquipment,
		OutEquipment:                r.OutEquipment,
		InItems:                     r.InItems,
		OutItems:                    r.OutItems,
		InWages:                     r.InWages,
		OutWages:                    r.OutWages,
		InEquipmentInAlliance:       r.InEquipmentInAlliance,
		OutEquipmentInAlliance:      r.OutEquipmentInAlliance,
		InItemsInAlliance:           r.InItemsInAlliance,
		OutItemsInAlliance:          r.OutItemsInAlliance,
		InWagesInAlliance:           r.InWagesInAlliance,
		OutWagesInAlliance:          r.OutWagesInAlliance,
		InEquipmentOutsideAlliance:  r.InEquipmentOutsideAlliance,
		OutEquipmentOutsideAlliance: r.OutEquipmentOutsideAlliance,
		InItemsOutsideAlliance:      r.InItemsOutsideAlliance,
		OutItemsOutsideAlliance:     r.OutItemsOutsideAlliance,
		InWagesOutsideAlliance:      r.InWagesOutsideAlliance,
		OutWagesOutsideAlliance:     r.OutWagesOutsideAlliance,
		Counterparts:                counterparts,
		AllianceID:                  r.AllianceID.Hex(),
	}
}
