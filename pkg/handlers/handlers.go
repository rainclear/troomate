package handlers

import (
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
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{})
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

	render.RenderTemplate(w, "accounts.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) NewAccount(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "new_account.page.html", &models.TemplateData{})
}
