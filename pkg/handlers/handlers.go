package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

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

	accountNameMap := make(map[int64]string)
	for _, account := range accounts {
		accound_id := account.ID
		accountNameMap[accound_id] = account.AccountName
	}

	render.RenderTemplate(w, r, "accounts.page.html", &models.TemplateData{
		IntKeyMap: accountNameMap,
	})
}

func (m *Repository) NewAccount(w http.ResponseWriter, r *http.Request) {
	log.Printf("In NewAccount...")
	render.RenderTemplate(w, r, "new_account.page.html", &models.TemplateData{})
}

func (m *Repository) PostNewAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
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

	http.Redirect(w, r, "accounts", http.StatusSeeOther)
}

func (m *Repository) ModifyAccount(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// var accound_id int64
	// for key, value := range r.Form {
	// 	fmt.Printf("%s = %s\n", key, value)
	// 	if key == "id" {
	// 		num, err := strconv.ParseInt(value[0], 10, 64)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 			return
	// 		}
	// 		accound_id = num
	// 	}
	// }

	id := r.URL.Query().Get("id")
	var account models.Account
	isEdit := false

	log.Printf("ModifyAccount, id=%s", id)

	if id != "" {
		accound_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
			return
		}

		account, err = dbm.GetAccountInfo(accound_id)
		if err != nil {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		isEdit = true
	} else {
		// For adding a new account, use an empty account struct
		account = models.Account{}
		isEdit = false
	}

	log.Printf("AccountName:%s, AccountNameAtCRA:%s", account.AccountName, account.AccountNameAtCRA)
	// Render the form
	tmpl := template.Must(template.ParseFiles("modify_account.page.html"))
	tmpl.Execute(w, struct {
		Account models.Account
		IsEdit  bool
	}{
		Account: account,
		IsEdit:  isEdit,
	})
}

func (m *Repository) EditAnAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}

	var accound_id int64
	var accountname = ""
	var accountname_at_cra = ""

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
		if key == "id" {
			num, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				log.Fatal(err)
				return
			}
			accound_id = num
		}
		if key == "inputAccountName" {
			accountname = value[0]
		}
		if key == "inputAccountNameAtCRA" {
			accountname_at_cra = value[0]
		}
	}

	var account models.Account

	if accound_id > 0 {
		account, err = dbm.GetAccountInfo(accound_id)

		if err != nil {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Account Id not existing", http.StatusNotFound)
		return
	}

	if accountname != "" && accountname != account.AccountName &&
		accountname_at_cra != "" && accountname_at_cra != account.AccountNameAtCRA {
		err = dbm.UpdateAnAccount(accound_id, accountname, accountname_at_cra)

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("accountname: %s, accountname_at_cra: %s\n", accountname, accountname_at_cra)
		}
	}

	http.Redirect(w, r, "accounts", http.StatusSeeOther)
}

func (m *Repository) DeleteAnAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}

	var accound_id int64

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
		if key == "id" {
			num, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				log.Fatal(err)
				return
			}
			accound_id = num
		}
	}

	if accound_id > 0 {
		err = dbm.DeleteAnAccount(accound_id)
		if err != nil {
			log.Fatal(err)
			return
		} else {
			fmt.Printf("account id to be deleted: %d\n", accound_id)
		}
	}

	http.Redirect(w, r, "accounts", http.StatusSeeOther)
}
