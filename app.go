package main

import (
	"database/sql"
	"go-mux/models"
	"go-mux/repos"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// tom: for Initialize
	"fmt"
	"log"

	"encoding/json"
	// tom: for route handlers
	"net/http"
	"strconv"

	// tom: go get required
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	GormDB *gorm.DB
	Repo   repos.Repos
}

// tom: initial function is empty, it's filled afterwards
// func (a *App) Initialize(user, password, dbname string) { }

// tom: added "sslmode=disable" to connection string

func (a *App) Initialize(user, password, dbname string) {
	fmt.Println("Initialization of DB Connection...")
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error

	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.GormDB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	a.GormDB.AutoMigrate(&models.Store{}, &models.StoreProducts{})
	fmt.Println("DB Connection Successful")
	a.Router = mux.NewRouter()

	// tom: this line is added after initializeRoutes is created later on
	a.initializeRoutes()
}

// tom: initial version
// func (a *App) Run(addr string) { }
// improved version

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

// tom: these are added later
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting product")
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	fmt.Println("Test123")
	products, err := getProducts(a.DB, start, count)
	fmt.Println("Post test working fine")
	if err != nil {
		fmt.Println("Error occured here")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a new Product")
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) createStore(w http.ResponseWriter, r *http.Request) {
	var s models.Store
	decoder := json.NewDecoder(r.Body)
	fmt.Println(decoder)
	err := decoder.Decode(&s)
	fmt.Println(err)
	fmt.Println("Creating new Store")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err = a.Repo.CreateNewStore(a.GormDB, &s)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) addStoreProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeId, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	//fmt.Println(storeId)
	var prods []*models.StoreProducts
	//fmt.Println(prods)
	decoder := json.NewDecoder(r.Body)
	//fmt.Println(&decoder)
	err = decoder.Decode(&prods)
	//fmt.Println(err)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	err = a.Repo.AddProductsToStore(a.GormDB, storeId, prods)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, prods)

}

func (a *App) getStoreProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeId, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	prods, err := a.Repo.GetAllProductsFromStore(a.GormDB, storeId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, prods)
	return
}

func (a *App) initializeRoutes() {
	fmt.Println("Setting up the routes")
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	a.Router.HandleFunc("/store", a.createStore).Methods("POST")
	a.Router.HandleFunc("/store/{id:[0-9]+}/products", a.getStoreProducts).Methods("GET")
	a.Router.HandleFunc("/store/{id:[0-9]+}", a.addStoreProducts).Methods("POST")

}
