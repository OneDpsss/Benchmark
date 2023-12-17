package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/jackc/pgx/v4"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

func PostgreSQL() {
	fmt.Println("PostgreSQL:")
	connString := "user=postgres password=admin host=localhost dbname=postgres sslmode=disable"
	conn, _ := pgx.Connect(context.Background(), connString)
	times := []int64{}
	// 1 query
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT cab_type, count(*) FROM taxi_rides GROUP BY 1;"
		rows, _ := conn.Query(context.Background(), query)
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 1 query: %d", times[len(times)/2])
	fmt.Println()
	// 2 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT passenger_count, avg(total_amount) FROM taxi_rides GROUP BY 1;"
		//println(query)
		//query := "SELECT cab_type FROM taxi_rides;"
		rows, _ := conn.Query(context.Background(), query)
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 2 query: %d", times[len(times)/2])
	fmt.Println()
	// 3 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT  passenger_count, extract(year from pickup_datetime), count(*) FROM taxi_rides GROUP BY 1, 2;"
		rows, _ := conn.Query(context.Background(), query)
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 3 query: %d", times[len(times)/2])
	fmt.Println()
	// 4 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT  passenger_count,  extract(year from pickup_datetime), round(trip_distance),  count(*) FROM taxi_rides GROUP BY 1, 2, 3 ORDER BY 2, 4 desc;"
		rows, _ := conn.Query(context.Background(), query)
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 4 query: %d", times[len(times)/2])
	fmt.Println()
	defer conn.Close(context.Background())

}

func SQlite() {
	fmt.Println("SQLite:")
	database, err := sql.Open("sqlite3", "Taxi_Trips.db")
	if err != nil {
		log.Fatal(err)
	}
	times := []int64{}
	// 1 query
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT cab_type, count(*) FROM trips GROUP BY 1;"
		rows, _ := database.Query(query)
		for rows.Next() {
		}
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 1 query: %d", times[len(times)/2])
	fmt.Println()
	// 2 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT passenger_count, avg(total_amount) FROM trips GROUP BY 1;"
		//println(query)
		//query := "SELECT cab_type FROM trips;"
		rows, _ := database.Query(query)
		for rows.Next() {
		}
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 2 query: %d", times[len(times)/2])
	fmt.Println()
	// 3 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := `SELECT passenger_count, strftime('%Y', pickup_datetime) AS pickup_year, count(*)
		FROM trips
		GROUP BY 1, pickup_year
		`
		rows, _ := database.Query(query)
		for rows.Next() {
		}
		err := rows.Close()
		if err != nil {
			return
		}
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 3 query: %d", times[len(times)/2])
	fmt.Println()
	// 4 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := `SELECT  passenger_count,  strftime('%Y', pickup_datetime), round(trip_distance),  count(*) FROM trips GROUP BY 1, 2, 3 ORDER BY 2, 4 desc;`
		rows, _ := database.Query(query)
		for rows.Next() {
		}
		rows.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})

	fmt.Printf("Median time for 4 query: %d", times[len(times)/2])
	fmt.Println()
}

func Gota() {
	// 1 query
	fmt.Println("Dataframe(Pandas):")
	times := []int64{}
	for c := 0; c < 20; c++ {
		cvs, _ := os.Open("dataset.csv")
		olddf := dataframe.ReadCSV(cvs)
		start := time.Now()
		olddf = olddf.Arrange(dataframe.Sort("cab_type"))
		cabTypeCol := olddf.Col("cab_type")
		df := dataframe.New(cabTypeCol)
		Countdf := df.GroupBy("cab_type")
		Countdf.GetGroups()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
		cvs.Close()
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 1 query: %d", times[len(times)/2])
	fmt.Println()
	// 2 query
	times = nil
	for c := 0; c < 20; c++ {
		cvs, _ := os.Open("dataset.csv")
		trips := dataframe.ReadCSV(cvs)
		start := time.Now()
		groupedDf := trips.Select([]string{"passenger_count", "total_amount"}).GroupBy("passenger_count").GetGroups()
		for _, val := range groupedDf {
			val.Col("total_amount").Mean()
		}
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
		cvs.Close()

	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 2 query: %d", times[len(times)/2])
	fmt.Println()
	times = nil
	// 3 query
	for c := 0; c < 20; c++ {
		cvs, _ := os.Open("dataset.csv")
		trips := dataframe.ReadCSV(cvs)
		start := time.Now()
		selectedDf := trips.Select([]string{"passenger_count", "pickup_datetime"})
		yearCol := []int{}
		for _, val := range selectedDf.Col("pickup_datetime").Records() {
			parsedTime, _ := time.Parse("2006-01-02 15:04:05", val)
			yearCol = append(yearCol, parsedTime.Year())
		}
		selectedDf = selectedDf.Mutate(series.New(yearCol, series.Int, "year"))
		groupDf := selectedDf.GroupBy("passenger_count", "year").GetGroups()
		for _, val := range groupDf {
			val.Nrow()
			//fmt.Println(ind, val.Nrow())
		}
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
		cvs.Close()

	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 3 query: %d", times[len(times)/2])
	fmt.Println()
	// 4 query
	times = nil
	for c := 0; c < 20; c++ {
		cvs, _ := os.Open("dataset.csv")
		trips := dataframe.ReadCSV(cvs)
		start := time.Now()
		selectedDf := trips.Select([]string{"passenger_count", "pickup_datetime", "trip_distance"})

		distance := []int{}
		for _, val := range selectedDf.Col("trip_distance").Records() {
			k, _ := strconv.ParseFloat(val, 64)
			c := math.Round(k)
			distance = append(distance, int(c))
		}
		yearCol := []int{}
		for _, val := range selectedDf.Col("pickup_datetime").Records() {
			parsedTime, _ := time.Parse("2006-01-02 15:04:05", val)
			yearCol = append(yearCol, parsedTime.Year())
		}
		selectedDf = selectedDf.Mutate(series.New(yearCol, series.Int, "year"))
		selectedDf = selectedDf.Mutate(series.New(distance, series.Int, "round"))
		groupDf := selectedDf.GroupBy("passenger_count", "year", "round").GetGroups()
		for _, val := range groupDf {
			val.Nrow()
		}
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
		cvs.Close()
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 4 query: %d", times[len(times)/2])
	fmt.Println()
}

func DuckDb() {
	db, err := sql.Open("duckdb", "/Users/valerashavlaygin/IdeaProjects/Benchmark/duck.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	times := []int64{}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT cab_type, count(*) FROM tripss GROUP BY 1;"
		row, _ := db.Query(query)
		row.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 1 query: %d", times[len(times)/2])
	fmt.Println()
	// 2 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT passenger_count, avg(total_amount) FROM tripss GROUP BY 1;"
		row, _ := db.Query(query)
		row.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 2 query: %d", times[len(times)/2])
	fmt.Println()
	// 3 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT  passenger_count, extract(year from pickup_datetime), count(*) FROM tripss GROUP BY 1, 2;"
		row, _ := db.Query(query)
		row.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 3 query: %d", times[len(times)/2])
	fmt.Println()
	// 4 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		query := "SELECT  passenger_count,  extract(year from pickup_datetime), round(trip_distance),  count(*) FROM tripss GROUP BY 1, 2, 3 ORDER BY 2, 4 desc;"
		row, _ := db.Query(query)
		row.Close()
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 4 query: %d", times[len(times)/2])
	fmt.Println()
}

type TaxiTrip struct {
	gorm.Model
	ID                   int
	CabType              int
	PickupDatetime       string
	DropoffDatetime      string
	PassengerCount       float64
	TripDistance         float64
	RatecodeID           float64
	StoreAndFwdFlag      string
	PULocationID         int
	DOLocationID         int
	PaymentType          int
	FareAmount           float64
	Extra                float64
	MTATax               float64
	TipAmount            float64
	TollsAmount          float64
	ImprovementSurcharge float64
	TotalAmount          float64
	CongestionSurcharge  float64
	AirportFee           float64
}

func help() {
	db, err := gorm.Open(sqlite.Open("Taxi_Trips.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	db.AutoMigrate(&TaxiTrip{})
	file, err := os.Open("dataset.csv")
	if err != nil {
		log.Fatal("Failed to open CSV file:", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV records:", err)
	}

	for _, record := range records {
		cabType, _ := strconv.Atoi(record[1])
		passengerCount, _ := strconv.ParseFloat(record[4], 64)
		tripDistance, _ := strconv.ParseFloat(record[5], 64)
		ratecodeID, _ := strconv.ParseFloat(record[6], 64)
		pulocationID, _ := strconv.Atoi(record[8])
		dolocationID, _ := strconv.Atoi(record[9])
		paymentType, _ := strconv.Atoi(record[10])
		fareAmount, _ := strconv.ParseFloat(record[11], 64)
		extra, _ := strconv.ParseFloat(record[12], 64)
		mtaTax, _ := strconv.ParseFloat(record[13], 64)
		tipAmount, _ := strconv.ParseFloat(record[14], 64)
		tollsAmount, _ := strconv.ParseFloat(record[15], 64)
		improvementSurcharge, _ := strconv.ParseFloat(record[16], 64)
		totalAmount, _ := strconv.ParseFloat(record[17], 64)
		congestionSurcharge, _ := strconv.ParseFloat(record[18], 64)
		airportFee, _ := strconv.ParseFloat(record[19], 64)
		newTrip := TaxiTrip{
			CabType:              cabType,
			PickupDatetime:       record[2],
			DropoffDatetime:      record[3],
			PassengerCount:       passengerCount,
			TripDistance:         tripDistance,
			RatecodeID:           ratecodeID,
			StoreAndFwdFlag:      record[7],
			PULocationID:         pulocationID,
			DOLocationID:         dolocationID,
			PaymentType:          paymentType,
			FareAmount:           fareAmount,
			Extra:                extra,
			MTATax:               mtaTax,
			TipAmount:            tipAmount,
			TollsAmount:          tollsAmount,
			ImprovementSurcharge: improvementSurcharge,
			TotalAmount:          totalAmount,
			CongestionSurcharge:  congestionSurcharge,
			AirportFee:           airportFee,
		}
		db.Create(&newTrip)
	}

	log.Println("Data loaded successfully.")
}

func Gorm() {
	db, err := gorm.Open(sqlite.Open("Taxi_Trips.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	// 1 query
	times := []int64{}
	for i := 0; i < 20; i++ {
		var result []struct {
			CabType int
			Count   int
		}
		start := time.Now()
		db.Table("trips").
			Select("CAST(TRIM(cab_type) as INTEGER), COUNT(*) as count").
			Group("cab_type").
			Scan(&result)
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 1 query: %d", times[len(times)/2])
	fmt.Println()
	// 2 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		var result []struct {
			PassengerCount float64
			AvgTotalAmount float64
		}
		db.Table("trips").
			Select("CAST(TRIM(passenger_count) as FLOAT), AVG(total_amount) as avg_total_amount").
			Group("passenger_count").
			Scan(&result)
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 2 query: %d", times[len(times)/2])
	fmt.Println()
	// 3 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		var result []struct {
			PassengerCount float64
			Year           float64
			Count          int
		}

		db.Table("trips").
			Select("passenger_count, strftime('%Y', pickup_datetime) AS year, COUNT(*) as count").
			Group("passenger_count, year").
			Scan(&result)
		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 3 query: %d", times[len(times)/2])
	fmt.Println()
	// 4 query
	times = nil
	for i := 0; i < 20; i++ {
		start := time.Now()
		var result []struct {
			PassengerCount float64
			Year           float64
			RoundTripDist  float64
			Count          int
		}
		db.Table("trips").
			Select("passenger_count, strftime('%Y', pickup_datetime) AS year, ROUND(trip_distance) AS round_trip_dist, COUNT(*) as count").
			Group("passenger_count, year, round_trip_dist").
			Order("year, count DESC").
			Scan(&result)

		duration := time.Since(start)
		times = append(times, duration.Milliseconds())
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] > times[j]
	})
	fmt.Printf("Median time for 4 query: %d", times[len(times)/2])
	fmt.Println()

}

func all() {
	PostgreSQL()
	SQlite()
	Gota()
	DuckDb()
	Gorm()
}

func main() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading cfg")
		return
	}
	functionName := viper.GetString("main.functionToRun")
	switch functionName {
	case "PostgreSQL":
		PostgreSQL()
	case "SQLite":
		SQlite()
	case "Pandas":
		Gota()
	case "DuckDB":
		DuckDb()
	case "Gorm":
		Gorm()
	case "All":
		all()
	default:
		fmt.Println("Error")
	}
}
