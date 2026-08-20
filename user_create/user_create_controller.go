package user_create

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/useradmin/shared"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

// UiInterface defines the user create controller's UI interface
type UiInterface interface {
	shared.UiInterface
	UserCreate(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new user create controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

type userCreateControllerData struct {
	firstName      string
	lastName       string
	email          string
	successMessage string
}

// UserCreate handles the user create controller requests
func (u *ui) UserCreate(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the user create controller request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := u.prepareDataAndValidate(r)

	if errorMessage != "" {
		return hb.Swal(hb.SwalOptions{
			Icon: "error",
			Text: errorMessage,
		}).ToHTML()
	}

	if data.successMessage != "" {
		return hb.Wrap().
			Child(hb.Swal(hb.SwalOptions{
				Icon: "success",
				Text: data.successMessage,
			})).
			Child(hb.Script("setTimeout(() => {window.location.href = window.location.href}, 2000)")).
			ToHTML()
	}

	return u.
		modal(r, data).
		ToHTML()
}

func (u *ui) modal(r *http.Request, data userCreateControllerData) hb.TagInterface {
	submitUrl := shared.NewLinksFromRequest(r).UserCreate(nil)

	formGroupFirstName := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("First name")).
		Child(bs.FormInput().Name("user_first_name").Value(data.firstName))

	formGroupLastName := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("Last name")).
		Child(bs.FormInput().Name("user_last_name").Value(data.lastName))

	formGroupEmail := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("Email")).
		Child(bs.FormInput().Name("user_email").Value(data.email))

	modalID := "ModaluserCreate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("New user Create").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModaluserCreate').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	buttonSend := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Create & Edit").
		Class("btn btn-primary float-end").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModaluserCreate").
		HxTarget("body").
		HxSwap("beforeend")

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Close").
		Class("btn btn-secondary float-start").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(formGroupFirstName).
						Child(formGroupLastName).
						Child(formGroupEmail)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					Child(buttonSend)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})
}

func (u *ui) prepareDataAndValidate(r *http.Request) (data userCreateControllerData, errorMessage string) {
	if u.UserStore() == nil {
		return data, "User store is not configured"
	}

	// When AuthUser callback is provided, enforce authentication.
	// When nil, the host project is expected to gate the route with
	// its own auth middleware — skip the check.
	authUser := u.AuthUser(r)
	if authUser == nil && u.AuthUserField != nil {
		return data, "You are not logged in. Please login to continue."
	}

	data.firstName = strings.TrimSpace(req.GetStringTrimmed(r, "user_first_name"))
	data.lastName = strings.TrimSpace(req.GetStringTrimmed(r, "user_last_name"))
	data.email = strings.TrimSpace(req.GetStringTrimmed(r, "user_email"))

	if r.Method != http.MethodPost {
		return data, ""
	}

	if data.firstName == "" {
		return data, "user first name is required"
	}

	if data.lastName == "" {
		return data, "user last name is required"
	}

	if data.email == "" {
		return data, "user email is required"
	}

	if !govalidator.IsEmail(data.email) {
		return data, "user email is invalid"
	}

	// Check email uniqueness before creating.
	if u.BlindIndexEmail() != nil && u.VaultTokenizer() != nil {
		ids, err := u.BlindIndexEmail().Search(r.Context(), data.email, "equals")
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userCreateController email uniqueness check", slog.String("error", err.Error()))
			}
			return data, "Failed to verify email uniqueness. Please contact an administrator."
		}
		if len(ids) > 0 {
			return data, "A user with this email already exists."
		}
	} else {
		query := userstore.NewUserQuery().SetEmail(data.email)
		existing, err := u.UserStore().UserList(r.Context(), query)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userCreateController email uniqueness check", slog.String("error", err.Error()))
			}
			return data, "Failed to verify email uniqueness. Please contact an administrator."
		}
		if len(existing) > 0 {
			return data, "A user with this email already exists."
		}
	}

	user := userstore.NewUser()
	user.SetFirstName(data.firstName)
	user.SetLastName(data.lastName)
	user.SetEmail(data.email)

	err := u.UserStore().UserCreate(r.Context(), user)

	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("Error. At userCreateController > prepareDataAndValidate", slog.String("error", err.Error()))
		}
		return data, "Creating user failed. Please contact an administrator."
	}

	data.successMessage = "user created successfully."

	return data, ""
}
