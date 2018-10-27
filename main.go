package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"math/rand"
	"net/http"
	"time"
)

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Info  string  `json:"info"`
	Price float64 `json:"price"`
}

var products []Product

var productType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"info": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"product": &graphql.Field{
				Type:        productType,
				Description: "Get product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, product := range products {
							if int(product.ID) == id {
								return product, nil
							}
						}
					}
					return nil, nil
				},
			},
			"list": &graphql.Field{
				Type:        graphql.NewList(productType),
				Description: "Get product list",
				Resolve: func(parms graphql.ResolveParams) (interface{}, error) {
					return products, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type:        productType,
				Description: "Crete nre product",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"info": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					rand.Seed(time.Now().UnixNano())
					product := Product{
						ID:    int64(rand.Intn(100000)),
						Name:  params.Args["name"].(string),
						Info:  params.Args["info"].(string),
						Price: params.Args["price"].(float64),
					}
					products = append(products, product)
					return product, nil
				},
			},
			"update": &graphql.Field{
				Type:        productType,
				Description: "Update product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"info": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.Float,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					name, nameOk := params.Args["name"].(string)
					info, infoOk := params.Args["info"].(string)
					price, priceOk := params.Args["price"].(float64)
					product := Product{}
					for i, p := range products {
						if int64(id) == p.ID {
							if nameOk {
								products[i].Name = name
							}

							if infoOk {
								products[i].Info = info
							}

							if priceOk {
								products[i].Price = price
							}
							product = products[i]
							break
						}
					}
					return product, nil
				},
			},
			"delete": &graphql.Field{
				Type:        productType,
				Description: "Delete product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},

				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					product := Product{}
					for i, p := range products {
						if int64(id) == p.ID {
							product = products[i]
							products = append(products[:i], products[i+1:]...)
						}
					}
					return product, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, shcema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func initProductsData(p *[]Product) {
	product1 := Product{ID: 1, Name: "abc", Info: "abc", Price: 70}
	producr2 := Product{ID: 2, Name: "efg", Info: "efg", Price: 10}
	product3 := Product{ID: 3, Name: "hij", Info: "hij", Price: 20}
	*p = append(*p, product1, producr2, product3)
}

func main() {
	initProductsData(&products)
	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}
