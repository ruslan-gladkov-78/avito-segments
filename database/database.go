package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type API struct {
	DB *sql.DB
}

func (api API) CreateSegment(slug string) error {
	var alreadyExists bool
	err := api.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Segments WHERE slug=$1)",
		slug).Scan(&alreadyExists)
	if err != nil {
		return err
	}
	if alreadyExists {
		return nil
	}
	_, err = api.DB.Exec(
		"INSERT INTO Segments (slug) VALUES($1)", slug,
	)
	return err
}

func (api API) DeleteSegment(slug string) error {
	var exists bool
	err := api.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Segments WHERE slug=$1)",
		slug).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	var segmentID int
	err = api.DB.QueryRow("SELECT segment_id FROM Segments WHERE slug=$1",
		slug).Scan(&segmentID)
	if err != nil {
		return err
	}
	rows, e := api.DB.Query("SELECT user_id, segments_ids FROM Users WHERE $1 = any (segments_ids)",
		segmentID)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var segmentsIDs []sql.NullInt32
		if err := rows.Scan(&id, pq.Array(&segmentsIDs)); err != nil {
			return err
		}
		newSegmentsIDs := make([]int, 0)
		for i := 0; i < len(segmentsIDs); i++ {
			if segmentID != int(segmentsIDs[i].Int32) {
				newSegmentsIDs = append(newSegmentsIDs, int(segmentsIDs[i].Int32))
			}
		}
		_, err = api.DB.Exec(
			"UPDATE Users SET segments_ids=$1 WHERE user_id=$2",
			newSegmentsIDs, id,
		)
		if err != nil {
			return err
		}
	}
	_, err = api.DB.Exec(
		"DELETE FROM Segments WHERE slug=$1", slug,
	)
	return err
}

func (api API) ChangeSegments(toAdd []string, toDelete []string, id int, TTL int) error {
	var userExists bool
	err := api.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE user_id=$1)",
		id).Scan(&userExists)
	if err != nil {
		return err
	}
	toAddIDs := make([]int, len(toAdd))
	for i, toAddSlug := range toAdd {
		err = api.DB.QueryRow("SELECT segment_id FROM Segments WHERE Slug=$1",
			toAddSlug).Scan(&(toAddIDs[i]))
		if err != nil {
			return err
		}
	}
	if !userExists {
		_, err = api.DB.Exec(
			"INSERT INTO Users VALUES($1, $2)",
			id, toAddIDs,
		)
		if TTL != 0 {
			go func() {
				time.Sleep(time.Duration(TTL) * time.Second)
				api.ChangeSegments([]string{}, toAdd, id, 0)
			}()
		}
		return err
	} else {
		toDeleteIDs := make([]int, len(toDelete))
		var oldIDsSQL []sql.NullInt32
		for i, toDeleteSlug := range toDelete {
			err = api.DB.QueryRow("SELECT segment_id FROM Segments WHERE slug=$1",
				toDeleteSlug).Scan(&(toDeleteIDs[i]))
			if err != nil {
				return err
			}
		}
		err = api.DB.QueryRow("SELECT segments_ids FROM Users WHERE user_id =$1",
			id).Scan(pq.Array(&oldIDsSQL))
		if err != nil {
			return err
		}
		oldIDs := make([]int, len(oldIDsSQL))
		for i := 0; i < len(oldIDsSQL); i++ {
			oldIDs[i] = int(oldIDsSQL[i].Int32)
		}
		newSegmentsIDsMap := make(map[int]bool, 0)
		newSegmentsIDsSlice := make([]int, 0)
		for _, segmentID := range append(toAddIDs, oldIDs...) {
			newSegmentsIDsMap[segmentID] = true
		}
		for _, segmentID := range toDeleteIDs {
			newSegmentsIDsMap[segmentID] = false
		}
		for k, v := range newSegmentsIDsMap {
			if v {
				newSegmentsIDsSlice = append(newSegmentsIDsSlice, k)
			}
		}
		_, err = api.DB.Exec(
			"UPDATE Users SET segments_ids=$1 WHERE user_id=$2",
			newSegmentsIDsSlice, id,
		)
		if TTL != 0 {
			go func() {
				time.Sleep(time.Duration(TTL) * time.Second)
				api.ChangeSegments([]string{}, toAdd, id, 0)
			}()
		}
		return err
	}
}

func (api API) GetSegments(id int) ([]string, error) {
	var userExists bool
	err := api.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE user_id=$1)",
		id).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, fmt.Errorf("user does not exist")
	}

	var segmentsIDs []sql.NullInt32
	slugs := make([]string, 0)
	err = api.DB.QueryRow("SELECT segments_ids FROM Users WHERE user_id =$1",
		id).Scan(pq.Array(&segmentsIDs))
	if err != nil {
		return nil, err
	}
	for _, segmentID := range segmentsIDs {
		var slug string
		err = api.DB.QueryRow("SELECT slug FROM Segments WHERE segment_id =$1",
			int(segmentID.Int32)).Scan(&slug)
		if err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, nil
}
