package user_delete

import (
	"log/slog"
	"net/http"

	"github.com/dracory/useradmin/shared"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

// UiInterface defines the user delete controller's UI interface
type UiInterface interface {
	shared.UiInterface
	UserDelete(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new user delete controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

type userDeleteControllerData struct {
	userID         string
	user           userstore.UserInterface
	successMessage string
}

// UserDelete handles the user delete controller requests
func (u *ui) UserDelete(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the user delete controller request and returns HTML
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

func (u *ui) modal(r *http.Request, data userDeleteControllerData) hb.TagInterface {
	submitUrl := shared.NewLinksFromRequest(r).UserDelete(map[string]string{
		"user_id": data.userID,
	})

	modalID := "ModalUserDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupUserId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("user_id").
		Value(data.userID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-danger float-end").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalUserDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete User").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalUserDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

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
						Child(hb.Paragraph().Text("Are you sure you want to delete this user?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupUserId)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(
						hb.Button().HTML("Close").
							Class("btn btn-secondary float-start").
							Data("bs-dismiss", "modal").
							OnClick(modalCloseScript)).
					Child(buttonDelete)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().
		Children([]hb.TagInterface{
			modal,
			backdrop,
		})
}

func (u *ui) prepareDataAndValidate(r *http.Request) (data userDeleteControllerData, errorMessage string) {
	if u.UserStore() == nil {
		return data, "User store is not configured"
	}

	// When AuthUser callback is provided, enforce authentication.
	// When nil, the host project is expected to gate the route with
	// its own auth middleware — skip the check.
	authUser := u.AuthUser(r)
	data.userID = req.GetString(r, "user_id")

	if authUser == nil && u.AuthUserField != nil {
		return data, "You are not logged in. Please login to continue."
	}

	if data.userID == "" {
		return data, "user id is required"
	}

	user, err := u.UserStore().UserFindByID(r.Context(), data.userID)

	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("Error. At userDeleteController > prepareDataAndValidate", slog.String("error", err.Error()))
		}
		return data, "User not found"
	}

	if user == nil {
		return data, "User not found"
	}

	data.user = user

	if r.Method != "POST" {
		return data, ""
	}

	err = u.UserStore().UserSoftDelete(r.Context(), user)

	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("Error. At userDeleteController > prepareDataAndValidate", slog.String("error", err.Error()))
		}
		return data, "Deleting user failed. Please contact an administrator."
	}

	data.successMessage = "user deleted successfully."

	return data, ""
}
