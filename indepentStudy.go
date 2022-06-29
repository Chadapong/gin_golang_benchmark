package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // import godotenv
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Username  string `gorm:"primary_key" json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

type UserDto struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type HealthCare struct {
	Index            int64   `json:"index"`
	Id               string  `gorm:"primary_key" json:"id"`
	HeartDisease     bool    `json:"heartDisease"`
	Bmi              int     `json:"bmi"`
	Smoking          bool    `json:"smoking"`
	Alcoholdrinking  bool    `json:"alcoholdrinking"`
	Stroke           bool    `json:"stroke"`
	PhysicalHealth   float64 `json:"physicalHealth"`
	MentalHealth     float64 `json:"mentalHealth"`
	DiffWalking      bool    `json:"diffWalking"`
	Sex              string  `json:"sex"`
	AgeCategory      string  `json:"ageCategory"`
	Race             string  `json:"race"`
	Diabetic         string  `json:"diabetic"`
	PhysicalActivity bool    `json:"physicalActivity"`
	GenHealth        string  `json:"genHealth"`
	SleepTime        float64 `json:"sleepTime"`
	Asthma           bool    `json:"asthma"`
	KidneyDisease    bool    `json:"kidneyDisease"`
	SkinCancer       bool    `json:"skinCancer"`
}

type ServiceModelResult struct {
	Records     []HealthCare `json:"records"`
	ExecuteTime float64      `json:"executeTime"`
}

type HealthCareHandler struct {
	DB *gorm.DB
}

func (h *HealthCareHandler) Initialize() {
	err := godotenv.Load("./conf.env")
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	// PORT := os.Getenv("PORT")
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", DB_HOST, DB_USER, DB_NAME, DB_PORT, DB_PASSWORD)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("%s", err)
	} else {
		h.DB = db
	}
}

func (h *HealthCareHandler) GetVeryComplexQuery(c *gin.Context) {
	returnResult := []HealthCare{}
	var avgExeTime float64 = 0.0

	for i := 0; i < 5; i++ {
		result := []HealthCare{}
		db_query := fmt.Sprintf(`SELECT * FROM health_cares_%d
		WHERE 
		heart_disease=true AND 
		bmi = (SELECT ROUND(avgBmi::numeric,0) FROM (SELECT AVG(bmi) as avgBmi FROM health_cares_%d WHERE sex='Female' AND age_category=(SELECT MIN(age_category) FROM health_cares_0) AND sleep_time=(SELECT ROUND(avgSleep::numeric,0) FROM (SELECT AVG(sleep_time) as avgSleep FROM health_cares_0) as bmi_avg_t)) as avgBmiT) AND 
		sleep_time = (SELECT ROUND(avg_sleep::numeric,0) FROM (SELECT AVG(sleep_time) as avg_sleep FROM health_cares_%d WHERE sex='Female' AND bmi=(SELECT AVG(avg_bmi) FROM (SELECT race, ROUND(avg_bmi::numeric,0) as avg_bmi FROM (SELECT race,AVG(bmi) as avg_bmi FROM health_cares_%d GROUP BY race) as complex_query_cond2) as round_AvgSleepT)) as roundAvgSleepT) AND 
		sex= (SELECT sex FROM(SELECT COUNT(sex) as numberPP,sex FROM health_cares_%d GROUP BY sex) as numberPPT WHERE numberPP= (SELECT MAX(totalPP) FROM (SELECT COUNT(sex) as totalPP,sex FROM health_cares_%d GROUP BY sex) as maxPPT)) AND  
		age_category = (SELECT MAX(DISTINCT age_category) FROM health_cares_%d WHERE bmi = (SELECT ROUND(avgBmi,0) FROM (SELECT AVG(bmi) as avgBmi FROM health_cares_0) as avgBmiT)) AND 
		physical_health = (SELECT ROUND(avgPh::numeric,0) FROM (SELECT AVG(physical_health) as avgPh FROM health_cares_0) as roundPH) ORDER BY index`, i, i, i, i, i, i, i)
		start := time.Now()
		resSql := h.DB.Raw(db_query).Scan(&result)
		if resSql.Error != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		end := time.Since(start)
		print(i)
		print("\n")
		print(end.Seconds())
		print("\n")
		print(avgExeTime)
		print("\n")
		avgExeTime += end.Seconds()
		print(avgExeTime)
		print("\n")
		if i == 4 {
			returnResult = result
		}
	}
	avgExeTime = avgExeTime / 5.0
	c.JSON(http.StatusOK, ServiceModelResult{Records: returnResult, ExecuteTime: avgExeTime})
}

func (h *HealthCareHandler) GetComplexQuery(c *gin.Context) {

	var avgExeTime float64 = 0.0

	returnResult := []HealthCare{}
	for i := 0; i < 5; i++ {
		result := []HealthCare{}
		db_query := fmt.Sprintf(`SELECT * FROM health_cares_%d 
		WHERE 
		bmi = (SELECT CAST(bmiValue as double precision) FROM (SELECT ROUND(avgBmi::numeric,0) as bmiValue FROM(SELECT AVG(bmi) as avgBmi FROM health_cares_%d  ) as bmiTable) as bmiT) AND 
		physical_health = (SELECT ROUND(avgPhy::numeric,0) FROM(SELECT AVG(physical_health) as avgPhy FROM health_cares_%d  ) as avgPhyT) AND 
		mental_health = (SELECT ROUND(avgMental::numeric,0) FROM(SELECT AVG(mental_health) as avgMental FROM health_cares_%d  ) as avgMentalT) AND 
		sleep_time = (SELECT ROUND(avgSleep::numeric,0) as avgRoundSleep FROM (SELECT gen_health, AVG(sleep_time) as avgSleep FROM health_cares_%d   GROUP BY gen_health) as T WHERE T.gen_health='Very good') 
		ORDER BY index`, i, i, i, i, i)
		start := time.Now()
		resSql := h.DB.Raw(db_query).Scan(&result)
		// temp, _ := (json.MarshalIndent(result, "", "\t"))
		// print(string(temp))
		if resSql.Error != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		end := time.Since(start)
		print(i)
		print("\n")
		print(end.Seconds())
		print("\n")
		print(avgExeTime)
		print("\n")
		avgExeTime += end.Seconds()
		print(avgExeTime)
		print("\n")
		if i == 4 {
			returnResult = result
		}
	}

	avgExeTime = avgExeTime / 5.0
	c.JSON(http.StatusOK, ServiceModelResult{Records: returnResult, ExecuteTime: avgExeTime})

}

// User DB
func (h *HealthCareHandler) CreateUser(c *gin.Context) {
	start := time.Now()
	var newAccount User
	if err := c.ShouldBindJSON(&newAccount); err != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAccount.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return
	}
	newAccount.Password = string(hashedPassword)

	if err := h.DB.Create(&newAccount).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	end := time.Since(start)
	c.JSON(http.StatusOK, map[string]interface{}{"executeTime": end.Seconds()})
}

func (h *HealthCareHandler) GetAllUser(c *gin.Context) {
	start := time.Now()
	users := []User{}
	if err := h.DB.Find(&users).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	end := time.Since(start)
	c.JSON(http.StatusOK, map[string]interface{}{"records": users, "executeTime": end.Seconds()})
}

func (h *HealthCareHandler) EditUser(c *gin.Context) {
	start := time.Now()
	var model UserDto
	if err := c.ShouldBindJSON(&model); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var existRecord User

	if err := h.DB.Find(&existRecord, User{Username: model.Username}).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	existRecord.Firstname = model.Firstname
	existRecord.Lastname = model.Lastname

	if err := h.DB.Model(&existRecord).Where("username =?", existRecord.Username).Omit("username", "password").Save(&existRecord).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	end := time.Since(start)
	c.JSON(http.StatusOK, map[string]interface{}{"executeTime": end.Seconds()})

}

func (h *HealthCareHandler) DeleteUser(c *gin.Context) {
	start := time.Now()
	username := c.Param("username")
	var existRecord User

	if err := h.DB.Find(&existRecord, User{Username: username}).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Where("username=?", username).Delete(&existRecord).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	end := time.Since(start)
	c.JSON(http.StatusOK, map[string]interface{}{"executeTime": end.Seconds()})

}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	h := HealthCareHandler{}
	h.Initialize()

	// healthcare
	r.GET("/very-complex-query", h.GetVeryComplexQuery)
	r.GET("/complex-query", h.GetComplexQuery)

	//user
	r.POST("/create", h.CreateUser)
	r.GET("/getUser", h.GetAllUser)
	r.PUT("/editUser", h.EditUser)
	r.DELETE("/deleteUser/:username", h.DeleteUser)

	return r
}

func main() {

	r := setupRouter()
	fmt.Println("Server Running on Port: ", 5000)

	r.Run("localhost:5000")
	// http.ListenAndServe(":5000", r)
}
