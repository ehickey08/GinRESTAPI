package recommendations

import "fmt"

var ColumnsWanted = []map[string]string{
	{UnitID: "main.UNITID"},
	{SchoolName: "main.INSTNM"},
	{City: "main.CITY"},
	{State: "main.STABBR"},
	{SchoolSetting: "main.LOCALE"},
	{SchoolType: "main.SECTOR"},
	{SchoolURL: "main.webaddr"},
	{ReligiousAffiliation: "ic.RELAFFIL"},
	{StudentPopulation: "dv.EFUGFT"},
	{CostInStateOnCampus: "cost.CINSON"},
	{CostOutStateOnCampus: "cost.COTSON"},
	{AcceptanceRate: "adm.ADMSSN/adm.APPLCN"},
	{Act: "adm.ACTCM75"},
	{Sat: "(adm.SATMT75 + adm.SATVR75)/2"},
	{ParentalIncomeTopTenPct: "mr2.fractionParTop10pctile"},
	{ParentalIncomeTopEightyPct: "(1 - mr2.fractionParBottomQuint1)"},
	{MedianEarningsAt34: "mr2.medianStudIndivEarnings2014"},
	{IndianOrNative: "dv.PCTENRAN"},
	{AsianOrPacific: "dv.PCTENRAP"},
	{Black: "dv.PCTENRBK"},
	{Hispanic: "dv.PCTENRHS"},
	{White: "dv.PCTENRWH"},
	{International: "dv.PCTENRNR"},
	{Multi: "dv.PCTENR2M"},
	{Unknown: "dv.PCTENRUN"},
	{PctStudentBodyFemale: "dv.PCTENRW"},
	{PctWithDisability: "ic.DISABPCT"},
	{LessThan3PCTDisability: "ic.DISAB"},
}

var TablesNeeded = []map[string][]string{
	{"som": {fmt.Sprintf("%s.%s", CompanyData, IDConversion),
		"round(main.opeid/100)",
		"som.opeid"}},
	{"mr2": {fmt.Sprintf("%s.%s", CompanyData, MobilityRateTwo),
		"som.super_opeid",
		"mr2.super_opeid"}},
	{"adm": {fmt.Sprintf("%s.%s", IpedsData, Admissions), "main.unitID", "adm.unitID"}},
	{"dv": {fmt.Sprintf("%s.%s", IpedsData, Diversity), "main.unitID", "dv.unitID"}},
	{"cost": {fmt.Sprintf("%s.%s", IpedsData, Cost), "main.unitID", "cost.unitID"}},
	{"ic": {fmt.Sprintf("%s.%s", IpedsData, IC), "main.unitID", "ic.UNITID"}},
}
