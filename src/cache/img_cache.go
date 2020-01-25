package cache

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GradeyCullins/GoogleVisionFilter/src/types"
)

// func CheckImgCache(db *sql.DB, imgURIs []string) (*types.ImgFilterRes, []string, error) {
// 	return nil, nil, nil
// }

// CheckImgCache checks the database cache for
func CheckImgCache(db *sql.DB, imgURIList []string) (*types.BatchImgFilterRes, error) {
	rows, err := db.Query("SELECT * from users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user types.User
		if err = rows.Scan(&user.UID, &user.Email, &user.Password); err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
	}

	return nil, nil
}
