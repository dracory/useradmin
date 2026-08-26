package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/neat"
	"github.com/dracory/useradmin/shared"
	"github.com/dracory/userstore"
)

const maxPerPage = 500

func (u *ui) handleUsersFetchAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}
	if u.UserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}

	// Parse request body
	var reqBody struct {
		Page        int    `json:"page"`
		PerPage     int    `json:"per_page"`
		SortOrder   string `json:"sort_order"`
		SortBy      string `json:"sort_by"`
		Status      string `json:"status"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		UserID      string `json:"user_id"`
		CreatedFrom string `json:"created_from"`
		CreatedTo   string `json:"created_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	// Helper functions for trimmed values with defaults (similar to req.GetStringTrimmedOr)
	getPositiveInt := func(val int, defaultVal int) int {
		if val <= 0 {
			return defaultVal
		}
		return val
	}

	getStringTrimmed := func(val string, defaultVal string) string {
		val = strings.TrimSpace(val)
		if val == "" {
			return defaultVal
		}
		return val
	}

	page := getPositiveInt(reqBody.Page, 0)
	perPage := getPositiveInt(reqBody.PerPage, 10)
	if perPage > maxPerPage {
		perPage = maxPerPage
	}
	sortOrder := getStringTrimmed(reqBody.SortOrder, neat.SortDesc)
	sortBy := getStringTrimmed(reqBody.SortBy, userstore.COLUMN_CREATED_AT)
	status := getStringTrimmed(reqBody.Status, "")
	firstName := getStringTrimmed(reqBody.FirstName, "")
	lastName := getStringTrimmed(reqBody.LastName, "")
	email := getStringTrimmed(reqBody.Email, "")
	userID := getStringTrimmed(reqBody.UserID, "")
	createdFrom := getStringTrimmed(reqBody.CreatedFrom, "")
	createdTo := getStringTrimmed(reqBody.CreatedTo, "")

	query := userstore.NewUserQuery().
		SetSortDirection(sortOrder).
		SetOrderBy(sortBy).
		SetOffset(page * perPage).
		SetLimit(perPage)

	if status != "" {
		query.SetStatus(status)
	}

	if userID != "" {
		query.SetID(userID)
	}

	if createdFrom != "" {
		query.SetCreatedAtGte(createdFrom + " 00:00:00")
	}

	if createdTo != "" {
		query.SetCreatedAtLte(createdTo + " 23:59:59")
	}

	// When OnUserSearch is provided, use it to find matching user IDs
	// (e.g. blind index, Elasticsearch). When nil, fall back to
	// userstore query-based search (SetFirstNameLike, etc.).
	var filteredIDs []string
	if firstName != "" || lastName != "" || email != "" {
		if u.OnUserSearch() != nil {
			ids, err := u.OnUserSearch()(r.Context(), shared.UserSearchEvent{
				FirstName:  firstName,
				LastName:   lastName,
				Email:      email,
				ExactMatch: false,
			})
			if err != nil {
				if u.Logger() != nil {
					u.Logger().Error("userManagerController.handleUsersFetchAjax OnUserSearch", slog.String("error", err.Error()))
				}
			}
			if len(ids) == 0 {
				api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
				return ""
			}
			filteredIDs = ids
			query.SetIDIn(filteredIDs)
		} else {
			// Fallback: userstore query-based search
			if firstName != "" {
				query.SetFirstNameLike(firstName)
			}
			if lastName != "" {
				query.SetLastNameLike(lastName)
			}
			if email != "" {
				query.SetEmailLike(email)
			}
		}
	}

	userList, err := u.UserStore().UserList(r.Context(), query)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userManagerController.handleUsersFetchAjax UserList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load users"))
		return ""
	}

	// Build a separate count query without limit/offset so UserCount
	// returns the total matching count, not just the current page count.
	countQuery := userstore.NewUserQuery().
		SetSortDirection(sortOrder).
		SetOrderBy(sortBy)

	if status != "" {
		countQuery.SetStatus(status)
	}
	if userID != "" {
		countQuery.SetID(userID)
	}
	if createdFrom != "" {
		countQuery.SetCreatedAtGte(createdFrom + " 00:00:00")
	}
	if createdTo != "" {
		countQuery.SetCreatedAtLte(createdTo + " 23:59:59")
	}
	if len(filteredIDs) > 0 {
		countQuery.SetIDIn(filteredIDs)
	}

	userCount, err := u.UserStore().UserCount(r.Context(), countQuery)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userManagerController.handleUsersFetchAjax UserCount", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to count users"))
		return ""
	}

	// Unseal PII for display. Prefer the batch callback for
	// efficiency; fall back to per-user; fall back to plain text.
	if u.UsersPiiUnseal() != nil {
		unsealed, err := u.UsersPiiUnseal()(r.Context(), userList)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUsersFetchAjax UsersPiiUnseal", slog.String("error", err.Error()))
			}
		} else {
			userList = unsealed
		}
	} else if u.UserPiiUnseal() != nil {
		for i, user := range userList {
			unsealed, err := u.UserPiiUnseal()(r.Context(), user)
			if err != nil {
				if u.Logger() != nil {
					u.Logger().Error("userManagerController.handleUsersFetchAjax UserPiiUnseal", slog.String("error", err.Error()))
				}
			} else {
				userList[i] = unsealed
			}
		}
	}

	users := make([]map[string]interface{}, 0, len(userList))
	for _, user := range userList {
		users = append(users, map[string]interface{}{
			FieldID:        user.GetID(),
			FieldFirstName: user.GetFirstName(),
			FieldLastName:  user.GetLastName(),
			FieldEmail:     user.GetEmail(),
			FieldStatus:    user.GetStatus(),
			FieldCreatedAt: user.GetCreatedAtCarbon().Format("d M Y"),
			FieldUpdatedAt: user.GetUpdatedAtCarbon().Format("d M Y"),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{
		FieldUsers: users,
		FieldTotal: userCount,
	}))
	return ""
}
