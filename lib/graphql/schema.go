package graphql

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"harper/api/login"
	db "harper/database"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

// Ticket contains information about one product
type Ticket struct {
	ID            int    `json:"id"`
	Status        string `json:"status"`
	Tracking      string `json:"tracking"`
	Stars         int    `json:"stars"`
	UsersID       int    `json:"user_id"`
	TechniciansID int    `json:"technician_id"`
}

type ticketTrack struct {
	Status string `db:"status"`
	Stars  int    `db:"stars"`
}

var tickets = []Ticket{}

var ticketType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Ticket",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"tracking": &graphql.Field{
				Type: graphql.String,
			},
			"stars": &graphql.Field{
				Type: graphql.Int,
			},
			"user_id": &graphql.Field{
				Type: graphql.Int,
			},
			"technician_id": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"ticket": &graphql.Field{
				Type:        ticketType,
				Description: "Get ticket by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					var ticket Ticket
					result := db.Conn.First(&ticket, id)
					if result.Error != nil {
						return nil, result.Error
					}
					link := "http://localhost:8080/service/tracking?id=" + ticket.Tracking
					ticket = Ticket{
						ID:            ticket.ID,
						Status:        ticket.Status,
						Tracking:      link,
						Stars:         ticket.Stars,
						UsersID:       ticket.UsersID,
						TechniciansID: ticket.TechniciansID,
					}
					return ticket, nil
				},
			},
			"ticketList": &graphql.Field{
				Type:        graphql.NewList(ticketType),
				Description: "Get ticket list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					//check if the person is a member of the company
					splitEmail := strings.SplitAfter(email, "@")
					if splitEmail[1] == "harper.com" {
						rows, err := db.Conn.Model(&Ticket{}).Where("technicians_id = ?", userID).Rows()
						if err != nil {
							return nil, err
						}
						defer rows.Close()
						for rows.Next() {
							// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
							db.Conn.ScanRows(rows, &tickets)
						}
					} else {
						rows, err := db.Conn.Model(&Ticket{}).Where("users_id = ?", userID).Rows()
						if err != nil {
							return nil, err
						}
						defer rows.Close()
						for rows.Next() {
							// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
							db.Conn.ScanRows(rows, &tickets)
						}
					}

					return tickets, nil
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create": &graphql.Field{
			Type:        ticketType,
			Description: "Create new service request",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				//hash for tracking, using email and date
				date := time.Now()
				h := md5.New()
				io.WriteString(h, email)
				io.WriteString(h, date.String())
				hash := h.Sum(nil)
				hashString := hex.EncodeToString(hash[:])
				link := "http://localhost:8080/service/tracking?id=" + hashString

				//select a technician randomly
				var technicianID int
				result := db.Conn.Raw(`SELECT id FROM technicians ORDER BY random() LIMIT 1`).Scan(&technicianID)
				if result.Error != nil {
					return nil, result.Error
				}
				//create ticket
				ticket := Ticket{Status: "created", Tracking: hashString, Stars: 0, UsersID: userID, TechniciansID: technicianID}
				result = db.Conn.Create(&ticket)
				if result.Error != nil {
					return nil, result.Error
				}
				ticket = Ticket{
					ID:            ticket.ID,
					Status:        ticket.Status,
					Tracking:      link,
					Stars:         ticket.Stars,
					UsersID:       ticket.UsersID,
					TechniciansID: ticket.TechniciansID,
				}
				return ticket, nil
			},
		},
		"rate": &graphql.Field{
			Type:        ticketType,
			Description: "Rate service",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"stars": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var ticket Ticket
				id, _ := params.Args["id"].(int)
				stars, starsOk := params.Args["stars"].(int)

				if starsOk {
					if stars < 1 && stars > 5 {
						err := errors.New("invalid range")
						return nil, err
					}
					//enter service rating to database
					result := db.Conn.Model(&Ticket{}).Where("id = ?", id).Update("stars", stars)
					if result.Error != nil {
						return nil, result.Error
					}
					result = db.Conn.First(&ticket, id)
					if result.Error != nil {
						return nil, result.Error
					}
					return ticket, nil
				}

				err := errors.New("service rating error")
				return nil, err
			},
		},
		"updateStatus": &graphql.Field{
			Type:        ticketType,
			Description: "Rate service",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"status": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var ticket Ticket
				id, _ := params.Args["id"].(int)
				status, statusOK := params.Args["status"].(string)

				if statusOK {
					//check if the person is not a member of the company
					splitEmail := strings.SplitAfter(email, "@")
					if splitEmail[1] != "harper.com" {
						err := errors.New("unauthorized user")
						return nil, err
					}
					//truncate the input to only 2 possible values
					if status != "created" && status != "finished" {
						err := errors.New("invalid value")
						return nil, err
					}
					//update status. Only a member of the company can do it
					result := db.Conn.Model(&Ticket{}).Where("id = ?", id).Update("status", status)
					if result.Error != nil {
						return nil, result.Error
					}
					result = db.Conn.First(&ticket, id)
					if result.Error != nil {
						return nil, result.Error
					}
					return ticket, nil
				}

				err := errors.New("error updating status")
				return nil, err
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

type reqBody struct {
	Query string `json:"query"`
}

var userID int
var email string

// Service create service request
func Service(router *gin.Engine) {
	routerService := router.Group("/service")
	{
		routerService.POST("/", login.AuthMiddleware.MiddlewareFunc(), func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			userID = int(claims["ID"].(float64))
			email = claims["Email"].(string)
			var rBody reqBody
			c.ShouldBind(&rBody)
			result := executeQuery(rBody.Query, schema)
			json.NewEncoder(c.Writer).Encode(result)
		})

		routerService.GET("/tracking", func(c *gin.Context) {
			idString := c.Query("id")
			var ticket Ticket
			result := db.Conn.Where("tracking = ?", idString).First(&ticket)
			notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)
			if notFound {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "not found tracking",
				})
				return
			}
			ticketRes := ticketTrack{
				Status: ticket.Status,
				Stars:  ticket.Stars,
			}
			json.NewEncoder(c.Writer).Encode(ticketRes)
		})
	}
}
