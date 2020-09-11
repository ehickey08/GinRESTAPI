package recommendations

import (
	"context"
	"github.com/ehickey08/GinRESTAPI/database"
	"log"
	"time"
)

func getSchools(desiredState string) ([]*School, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := database.DbConn.QueryContext(ctx, Sql, desiredState)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	schools := make([]*School, 0)
	for rows.Next() {
		school := &School{}
		if err := rows.Scan(&school.UnitID,
			&school.School,
			&school.City,
			&school.State,
			&school.Setting,
			&school.SchoolType,
			&school.Url,
			&school.ReligiousAffiliation,
			&school.StudentPopulation,
			&school.CostInStateOnCampus,
			&school.CostOutStateOnCampus,
			&school.AcceptanceRate,
			&school.Act,
			&school.Sat,
			&school.ParentalIncomeTopTenPct,
			&school.ParentalIncomeTopEightyPct,
			&school.MedianEarningsAt34,
			&school.IndianOrNative,
			&school.AsianOrPacific,
			&school.Black,
			&school.Hispanic,
			&school.White,
			&school.International,
			&school.Multi,
			&school.Unknown,
			&school.Female,
			&school.PctWithDisability,
			&school.LessThan3PCTDisability,
		); err != nil {
			log.Fatal(err)
			return nil, err
		}
		schools = append(schools, school)
	}
	return schools, nil
}
