package handlers

import (
	"net/http"

	"github.com/rainclear/accroo/pkg/config"
	"github.com/rainclear/accroo/pkg/models"
	"github.com/rainclear/accroo/pkg/render"
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

func (m *Repository) AccountTypes(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "account_types.page.html", &models.TemplateData{})
}
