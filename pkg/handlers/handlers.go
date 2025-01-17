package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rainclear/troomate/pkg/config"
	"github.com/rainclear/troomate/pkg/dbm"
	"github.com/rainclear/troomate/pkg/models"
	"github.com/rainclear/troomate/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{})
}

func (m *Repository) Accounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := dbm.ListAccounts()
	if err != nil {
		log.Fatal("Db Error")
	}

	stringMap := make(map[string]string)
	for _, account := range accounts {
		account_info := strings.Split(account, ",")
		accound_id := account_info[0]
		stringMap[accound_id] = account
	}

	render.RenderTemplate(w, r, "accounts.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) NewAccount(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "new_account.page.html", &models.TemplateData{})
}

func (m *Repository) PostNewAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var accountname = ""
	var accountname_at_cra = ""

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
		if key == "inputAccountName" {
			accountname = value[0]
		}
		if key == "inputAccountNameAtCRA" {
			accountname_at_cra = value[0]
		}
	}

	if accountname != "" && accountname_at_cra != "" {
		err = dbm.AddNewAccount(accountname, accountname_at_cra)

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("accountname: %s, accountname_at_cra: %s\n", accountname, accountname_at_cra)
		}
	}

	Repo.Accounts(w, r)
}
