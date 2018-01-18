package apiforum

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/apiutil"
	"github.com/reo7sp/technopark-db/dbutil"
	"log"
	"net/http"
	"strconv"
)

func MakeShowUsersHandler(db *pgx.ConnPool) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showUsersRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showUsersAction(w, in, db)
	}
	return f
}

type showUsersInput struct {
	Slug string `json:"-"`

	Limit  int64  `json:"limit"`
	Since  string `json:"since"`
	IsDesc bool   `json:"desc"`

	Order    string `json:"-"`
	LimitSql string `json:"-"`
}

type showUsersOutputItem api.UserModel

type showUsersOutput []showUsersOutputItem

func showUsersRead(r *http.Request, ps map[string]string) (in showUsersInput, err error) {
	in.Slug = ps["slug"]

	query := r.URL.Query()

	in.Limit, err = strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		err = nil
		in.Limit = 0
		in.LimitSql = ""
	} else {
		in.LimitSql = "LIMIT " + strconv.FormatInt(in.Limit, 10) + "::bigint"
	}
	in.Since = query.Get("since")
	in.IsDesc = query.Get("desc") == "true"
	if in.IsDesc {
		in.Order = "DESC"
	} else {
		in.Order = "ASC"
	}
	return
}

func showUsersCheckForumExists(slug string, db *pgx.ConnPool) (bool, error) {
	sqlQuery := "SELECT slug::text FROM forums WHERE slug = $1::citext"
	var s string
	err := db.QueryRow(sqlQuery, slug).Scan(&s)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func showUsersAction(w http.ResponseWriter, in showUsersInput, db *pgx.ConnPool) {
	doesForumExists, err := showUsersCheckForumExists(in.Slug, db)
	if err != nil {
		log.Println("error: apiforum.showUsersAction: showUsersCheckForumExists:", err)
		w.WriteHeader(500)
		return
	}
	if !doesForumExists {
		errJson := api.ErrorModel{Message: "Can't find forum"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}

	out := make(showUsersOutput, 0, in.Limit)

	sqlQuery := fmt.Sprintf(`

	WITH thisForumUsers AS (
		SELECT fu.nickname
		FROM forumUsers fu
		WHERE fu.forumSlug = $1::citext
		AND (
			CASE WHEN $2::citext != ''
			THEN (
				CASE WHEN $3::boolean IS TRUE
				THEN fu.nickname < $2::citext
				ELSE fu.nickname > $2::citext
				END
			)
			ELSE TRUE
			END
		)
		ORDER BY fu.nickname %s
		%s
	)
	SELECT u.nickname::text, u.fullname, u.email::text, u.about FROM thisForumUsers tfu
	JOIN users u ON (u.nickname = tfu.nickname)

	`, in.Order, in.LimitSql)

	rows, err := db.Query(sqlQuery, in.Slug, in.Since, in.IsDesc)
	if err != nil {
		log.Println("error: apiforum.showUsersAction: SELECT start:", err)
		w.WriteHeader(500)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var outItem showUsersOutputItem
		err = rows.Scan(&outItem.Nickname, &outItem.Fullname, &outItem.Email, &outItem.About)
		if err != nil {
			log.Println("error: apiforum.showUsersAction: SELECT iter:", err)
			w.WriteHeader(500)
			return
		}
		out = append(out, outItem)
	}

	apiutil.WriteJsonObject(w, out, 200)
}
