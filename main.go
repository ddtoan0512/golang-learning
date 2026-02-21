package main

import (
	restaurantmodel "golang-learning/module/restaurant/model"
	"golang-learning/module/restaurant/transport/ginrestaurant"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"name" gorm:"column:addr;"`
}

func (Restaurant) TableName() string { return "restaurants" }

func main() {
	dsn := "root:ead8686ba57479778a76e@tcp(127.0.0.1:3306)/food_delivery?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := os.Getenv("MSYSQL_CONN_STRING")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(db)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// POST /restaurants

	v1 := r.Group("/v1")

	restaurants := v1.Group("/restaurants")

	restaurants.POST("", ginrestaurant.CreateRestaurant(db))

	restaurants.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var data Restaurant

		db.Where("id = ?", id).First(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.GET("", func(c *gin.Context) {
		var data []Restaurant

		type Paging struct {
			Page  int `json:"page" form:"page"`
			Limit int `json:"limit" form:"limit"`
		}

		var pagingData Paging

		if err := c.ShouldBind(&pagingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if pagingData.Page <= 0 {
			pagingData.Page = 1
		}

		if pagingData.Limit <= 0 {
			pagingData.Limit = 5
		}

		db.Offset((pagingData.Page - 1) * pagingData.Limit).
			Order("id desc").
			Limit(pagingData.Limit).
			Find(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.PATCH("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var data restaurantmodel.RestaurantUpdate

		db.Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Table(Restaurant{}.TableName()).Where("id = ?", id).Delete(nil)

		c.JSON(http.StatusOK, gin.H{})
	})

	r.Run()

	// newRestaurant := Restaurant{
	// 	Name: "Duks",
	// 	Addr: "Ha Noi",
	// }

	// if err := db.Create(&newRestaurant).Error; err != nil {
	// 	log.Println(err)
	// }

	//var myRestaurant Restaurant
	//
	//if err := db.Where("id = ?", 2).First(&myRestaurant).Error; err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(myRestaurant)
	//
	//myRestaurant.Name = "Tosa"
	//
	//if err := db.Where("id = ?", 2).Updates(&myRestaurant).Error; err != nil {
	//	log.Println(err)
	//}
}
