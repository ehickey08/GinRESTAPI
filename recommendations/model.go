package recommendations

import (
	"database/sql"
	"encoding/json"
)

type NullFloat struct {
	sql.NullFloat64
}

type NullInt struct {
	sql.NullInt64
}

type NullString struct {
	sql.NullString
}

type NullAcceptance struct {
	sql.NullFloat64
}
type School struct {
	UnitID                     int        `json:"unitID"`
	School                     string     `json:"school"`
	City                       string     `json:"city"`
	State                      string     `json:"state"`
	Setting                    int        `json:"setting"`
	SchoolType                 int        `json:"schoolType"`
	Url                        string     `json:"url"`
	ReligiousAffiliation       int        `json:"religiousAffiliation"`
	StudentPopulation          NullInt    `json:"studentPopulation"`
	CostInStateOnCampus        NullInt    `json:"costInStateOnCampus"`
	CostOutStateOnCampus       NullInt    `json:"costOutStateOnCampus"`
	AcceptanceRate             NullFloat  `json:"acceptanceRate"`
	Act                        NullInt    `json:"act"`
	Sat                        NullFloat  `json:"sat"`
	ParentalIncomeTopTenPct    NullString `json:"parentalIncomeTopTenPct"`
	ParentalIncomeTopEightyPct NullString `json:"parentalIncomeTopEightyPct"`
	MedianEarningsAt34         NullInt    `json:"medianEarningsAt34"`
	IndianOrNative             NullFloat  `json:"Indian or Native"`
	AsianOrPacific             NullFloat  `json:"asianOrPacific"`
	Black                      NullFloat  `json:"black"`
	Hispanic                   NullFloat  `json:"hispanic"`
	White                      NullFloat  `json:"white"`
	International              NullFloat  `json:"international"`
	Multi                      NullFloat  `json:"multi"`
	Unknown                    NullFloat  `json:""unknown""`
	Female                     NullFloat  `json:"female"`
	PctWithDisability          NullFloat  `json:"PctWithDisability"`
	LessThan3PCTDisability     int        `json:"LessThan3PCTDisability"`
	SettingString              string     `json:"settingString"`
	TypeString                 string     `json:"typeString"`
	ActualCost                 NullInt    `json:"actualCost"`
	Diversity                  float64    `json:"diversity"`
	ADACompliance              int        `json:"adaCompliance"`
	DisposableIncome           float64    `json:"disposableIncome"`
	DisposableIncomeRank       int        `json:"disposableIncomeRank"`
	CostRank                   int        `json:"costRank"`
	ParentalIncomeRank         int        `json:"parentalIncomeRank"`
	ADARank                    int        `json:"adaRank"`
	DiversityRank              int        `json:"diversityRank"`
	DisposableIncomeScore      int        `json:"disposableIncomeScore"`
	CostScore                  int        `json:"costScore"`
	ParentalIncomeScore        int        `json:"parentalIncomeScore"`
	ADAScore                   int        `json:"adaScore"`
	DiversityScore             int        `json:"diversityScore"`
	FinalScore                 int        `json:"finalScore"`
	SubfieldOfStudy            float64    `json:"subfieldOfStudy"`
	BasefieldOfStudy           float64    `json:"basefieldOfStudy"`
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)
}
func (nf *NullFloat) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(nf.Float64)
}
func (ni *NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ni.Int64)
}
