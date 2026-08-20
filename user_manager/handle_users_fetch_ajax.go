package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/neat"
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

	// Collect blind index search results as separate ID sets so we can
	// intersect them (AND semantics). A user must match ALL active
	// filters to be included. If any filter returns zero results, the
	// overall result is empty.
	var filterIDSets [][]string

	if firstName != "" && u.BlindIndexFirstName() != nil {
		ids, err := u.BlindIndexFirstName().Search(r.Context(), firstName, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUsersFetchAjax blind index first_name", slog.String("error", err.Error()))
			}
		}
		if len(ids) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		filterIDSets = append(filterIDSets, ids)
	}

	if lastName != "" && u.BlindIndexLastName() != nil {
		ids, err := u.BlindIndexLastName().Search(r.Context(), lastName, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUsersFetchAjax blind index last_name", slog.String("error", err.Error()))
			}
		}
		if len(ids) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		filterIDSets = append(filterIDSets, ids)
	}

	if email != "" && u.BlindIndexEmail() != nil {
		ids, err := u.BlindIndexEmail().Search(r.Context(), email, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			if u.Logger() != nil {
				u.Logger().Error("userManagerController.handleUsersFetchAjax blind index email", slog.String("error", err.Error()))
			}
		}
		if len(ids) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		filterIDSets = append(filterIDSets, ids)
	}

	// Intersect all filter ID sets (AND semantics). When only one
	// filter is active, the intersection is that set itself.
	if len(filterIDSets) > 0 {
		intersected := intersectIDSets(filterIDSets)
		if len(intersected) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		query.SetIDIn(intersected)
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
	if len(filterIDSets) > 0 {
		intersected := intersectIDSets(filterIDSets)
		countQuery.SetIDIn(intersected)
	}

	userCount, err := u.UserStore().UserCount(r.Context(), countQuery)
	if err != nil {
		if u.Logger() != nil {
			u.Logger().Error("userManagerController.handleUsersFetchAjax UserCount", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to count users"))
		return ""
	}

	users := make([]map[string]interface{}, 0, len(userList))
	for _, user := range userList {
		firstNameVal := user.GetFirstName()
		lastNameVal := user.GetLastName()
		emailVal := user.GetEmail()

		if u.VaultTokenizer() != nil {
			var err error
			firstNameVal, lastNameVal, emailVal, _, _, err = u.VaultTokenizer().Untokenize(r.Context(), user)
			if err != nil {
				if u.Logger() != nil {
					u.Logger().Error("userManagerController.handleUsersFetchAjax Untokenize", slog.String("error", err.Error()))
				}
				firstNameVal = "n/a"
				lastNameVal = "n/a"
				emailVal = "n/a"
			}
		}

		users = append(users, map[string]interface{}{
			FieldID:        user.GetID(),
			FieldFirstName: firstNameVal,
			FieldLastName:  lastNameVal,
			FieldEmail:     emailVal,
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

// intersectIDSets computes the intersection of multiple ID slices.
// Returns IDs that appear in ALL input slices. If any slice is empty,
// the result is empty. Order follows the first slice.
func intersectIDSets(sets [][]string) []string {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0]
	}

	// Build a count map: how many sets each ID appears in.
	counts := make(map[string]int, len(sets[0]))
	order := make([]string, 0, len(sets[0]))
	for _, id := range sets[0] {
		if counts[id] == 0 {
			order = append(order, id)
		}
		counts[id]++
	}
	for _, s := range sets[1:] {
		seen := make(map[string]bool, len(s))
		for _, id := range s {
			if !seen[id] {
				seen[id] = true
				counts[id]++
			}
		}
	}

	// Keep only IDs that appear in all sets.
	result := make([]string, 0, len(order))
	for _, id := range order {
		if counts[id] == len(sets) {
			result = append(result, id)
		}
	}
	return result
}
