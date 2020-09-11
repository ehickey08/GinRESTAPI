package recommendations

import (
	"context"
	"fmt"
	"github.com/ehickey08/GinRESTAPI/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math"
	"net/http"
	"sort"
	"time"
)

const recommendationPath = "recommendations"

func SetupRoutes(router *gin.Engine, basePath string) {
	router.GET(fmt.Sprintf("%s/%s/:userID", basePath, recommendationPath), getRecommendations)
}

type User struct {
	HomeState         string
	DesiredState      string
	CostRank          int
	IncomeBracketRank int
	DiversityRank     int
	ADAComplianceRank int
	FutureIncomeRank  int
	Subfield          string
	Basefield         string
	BaseFieldCode     string
	SubfieldCode      string
	Auth0ID           string
}

func getRecommendations(c *gin.Context) {
	userID := c.Param("userID")
	mongoID, err := primitive.ObjectIDFromHex(userID)
	collection := database.Client.Database("student_analytics").Collection("users")
	var user User
	err = collection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	schools, err := getSchools(user.DesiredState)
	normalizedSchoolData(schools)
	finalSchools := scoreData(schools, user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, finalSchools)
}

func normalizedSchoolData(data []*School) {
	settingOptions := map[float64]string{
		1: "City",
		2: "Suburb",
		3: "Town",
		4: "Rural",
	}
	for _, school := range data {
		settingVal := math.Floor(float64(school.Setting) / 10)
		school.SettingString = settingOptions[settingVal]
		if school.SchoolType == 1 {
			school.TypeString = "Public"
		} else {
			school.TypeString = "Private"
		}
	}
}

func scoreData(data []*School, user User) []*School {
	for _, school := range data {
		determineCost(school, user.HomeState == user.DesiredState)
		determineDiversity(school)
		addADACompliance(school)
	}
	cleanedData := removeNulls(data)
	getFieldData(data, user)
	addFinancialData(cleanedData)
	getRanks(cleanedData, user)
	getFinalScores(cleanedData)
	return grabTopSchools(cleanedData)
}

func determineCost(school *School, stateMatch bool) {
	if stateMatch {
		school.ActualCost = school.CostInStateOnCampus
	} else {
		school.ActualCost = school.CostOutStateOnCampus
	}
}

func determineDiversity(school *School) {
	indSum := school.IndianOrNative.Float64*(school.IndianOrNative.Float64-1) + school.AsianOrPacific.Float64*(school.AsianOrPacific.Float64-1) + school.Black.Float64*(school.Black.Float64-1) + school.Hispanic.Float64*(school.Hispanic.Float64-1) + school.White.Float64*(school.White.Float64-1) + school.Multi.Float64*(school.Multi.Float64-1) + school.Unknown.Float64*(school.Unknown.Float64-1)
	nonResAmt := 100.0 - school.International.Float64
	totalSum := nonResAmt * (nonResAmt - 1.0)
	if totalSum > 0 {
		divisionRes := indSum / totalSum
		school.Diversity = 1.0 - divisionRes
	} else {
		school.Diversity = 0
	}

}

func addADACompliance(school *School) {
	if school.LessThan3PCTDisability == 1 {
		school.ADACompliance = 2
	} else {
		school.ADACompliance = int(school.PctWithDisability.Float64)
	}
}

func addFinancialData(data []*School) {
	for _, school := range data {
		postTaxInc := float64(school.MedianEarningsAt34.Int64/12) * 0.72
		loanPmt := calcLoanPmt(school.ActualCost)
		school.DisposableIncome = postTaxInc - loanPmt
	}
}

func calcLoanPmt(cost NullInt) float64 {
	pv := cost.Int64 * 4
	r := 0.0725 / 12.0
	n := float64(12 * 15)
	num := float64(pv) * r
	den := 1.0 - math.Pow(1+r, -n)
	return num / den
}

func removeNulls(data []*School) []*School {
	cleanedData := make([]*School, 0)
	for _, school := range data {
		if school.ActualCost.Valid && school.White.Valid && school.ParentalIncomeTopTenPct.Valid && school.MedianEarningsAt34.Valid {
			cleanedData = append(cleanedData, school)
		}
	}
	return cleanedData
}

func getRanks(data []*School, user User) {
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].DisposableIncome < data[j].DisposableIncome
	})
	baseValue := float64(0)
	currRank, offset := 1, 0
	for i, school := range data {
		currValue := school.DisposableIncome
		if currValue == baseValue && i != 0 {
			offset++
		} else {
			currRank += offset
			offset = 1
		}

		school.DisposableIncomeRank = currRank
		school.DisposableIncomeScore = currRank * user.FutureIncomeRank

		baseValue = currValue
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].ActualCost.Int64 > data[j].ActualCost.Int64
	})

	var intBaseValue int64 = 0
	currRank, offset = 1, 0
	for i, school := range data {
		currValue := school.ActualCost.Int64
		if currValue == intBaseValue && i != 0 {
			offset++
		} else {
			currRank += offset
			offset = 1
		}

		school.CostRank = currRank
		school.CostScore = currRank * user.CostRank

		intBaseValue = currValue
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Diversity < data[j].Diversity
	})
	baseValue = float64(0)
	currRank, offset = 1, 0
	for i, school := range data {
		currValue := school.Diversity
		if currValue == baseValue && i != 0 {
			offset++
		} else {
			currRank += offset
			offset = 1
		}

		school.DiversityRank = currRank
		school.DiversityScore = currRank * user.DiversityRank

		baseValue = currValue
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].ADACompliance < data[j].ADACompliance
	})
	justIntBaseValue := 0
	currRank, offset = 1, 0
	for i, school := range data {
		currValue := school.ADACompliance
		if currValue == justIntBaseValue && i != 0 {
			offset++
		} else {
			currRank += offset
			offset = 1
		}

		school.ADARank = currRank
		school.ADAScore = currRank * user.ADAComplianceRank

		justIntBaseValue = currValue
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].ParentalIncomeTopTenPct.String > data[j].ParentalIncomeTopTenPct.String
	})
	strBaseValue := ""
	currRank, offset = 1, 0
	for i, school := range data {
		currValue := school.ParentalIncomeTopTenPct.String
		if currValue == strBaseValue && i != 0 {
			offset++
		} else {
			currRank += offset
			offset = 1
		}

		school.ParentalIncomeRank = currRank
		school.ParentalIncomeScore = currRank * user.IncomeBracketRank

		strBaseValue = currValue
	}
}

func getFinalScores(data []*School) {
	for _, school := range data {
		school.FinalScore = school.CostScore + school.DiversityScore + school.DisposableIncomeScore + school.ADAScore + school.ParentalIncomeScore
	}
}

func grabTopSchools(data []*School) []*School {
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].FinalScore > data[j].FinalScore
	})
	if len(data) <= 6 {
		return data
	}
	haveSubfield := filter(data)
	if len(haveSubfield) > 5 {
		topSchools := make([]*School, 6)
		copy(topSchools, haveSubfield)
		return topSchools
	}
	haveBasefield := filter2(data)
	sort.SliceStable(haveBasefield, func(i, j int) bool {
		return haveBasefield[i].BasefieldOfStudy > haveBasefield[j].BasefieldOfStudy
	})
	schoolToAdd := 0
	for len(haveSubfield) < 6 && schoolToAdd < len(haveBasefield) {
		haveSubfield = append(haveSubfield, haveBasefield[schoolToAdd])
		schoolToAdd++
	}
	noBasefield := filter3(data)
	schoolToAdd = 0
	for len(haveSubfield) < 6 {
		haveSubfield = append(haveSubfield, noBasefield[schoolToAdd])
		schoolToAdd++
	}
	return haveSubfield

}

func filter(data []*School) (ret []*School) {
	for _, s := range data {
		if s.SubfieldOfStudy > 0 {
			ret = append(ret, s)
		}
	}
	return
}
func filter2(data []*School) (ret []*School) {
	for _, s := range data {
		if s.SubfieldOfStudy == 0 && s.BasefieldOfStudy > 0 {
			ret = append(ret, s)
		}
	}
	return
}
func filter3(data []*School) (ret []*School) {
	for _, s := range data {
		if s.BasefieldOfStudy == 0 {
			ret = append(ret, s)
		}
	}
	return
}

type MongoSchool struct {
	UnitId int
	School string
	Cips   []string
	Data   []CipData
}

type CipData struct {
	CipCode   string
	PctOfDegs float64
	AmtOfDegs int
}

func getFieldData(data []*School, user User) {
	start := time.Now()
	for _, school := range data {
		go updateSchoolData(school, user)

	}
	fmt.Println(time.Since(start))
}

func updateSchoolData(school *School, user User) {
	collection := database.Client.Database("student_analytics").Collection("school_completion_data")
	mongoSchool := MongoSchool{}
	err := collection.FindOne(context.Background(), bson.M{"unitID": school.UnitID}).Decode(&mongoSchool)
	if err != nil {
		log.Fatal(err)
	} else {
		pctSubfield, pctBasefield := 0.0, 0.0
		for _, cip := range mongoSchool.Cips {
			if cip == user.SubfieldCode {
				for _, detailedCip := range mongoSchool.Data {
					if detailedCip.CipCode == user.SubfieldCode {
						pctSubfield = detailedCip.PctOfDegs
					}
				}
			}
			if cip == user.BaseFieldCode {
				for _, detailedCip := range mongoSchool.Data {
					if detailedCip.CipCode == user.BaseFieldCode {
						pctBasefield = detailedCip.PctOfDegs
					}
				}
			}
		}
		school.SubfieldOfStudy = pctSubfield
		school.BasefieldOfStudy = pctBasefield
	}
}
