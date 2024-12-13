package models


import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel type wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert method adds a new record to the snippets table
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
   stmt := `INSERT INTO snippets (title, content, created, expires)
           VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

   result, err := m.DB.Exec(stmt, title, content, expires)
   if err != nil {
	   return 0, err
   }

   id, err := result.LastInsertId()
   if err != nil {
	   return 0, err
   }
   
   return int(id), nil
}

// Get method returns a specific snippet based on its ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
  
   stmt := `SELECT id, title, content, created, expires FROM snippets
             WHERE expires > UTC_TIMESTAMP() AND id = ?`
   
   row := m.DB.QueryRow(stmt, id)
   s := &Snippet{}

   err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
   if err != nil {
	   if errors.Is(err, sql.ErrNoRows) {
		   return nil, ErrNoRecord
	   } else {
		   return nil, err
	   }
   }

   return s, nil
}

// Latest method returns the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	
	stmt := `SELECT id, title, content, created, expires FROM snippets
         WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
    
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

// func (m *ExampleModel) ExampleTransaction() error {
// 	tx, err := m.DB.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback() // This will be ignored if the tx is committed first

// 	_, err = tx.Exec("INSERT INTO ...")
// 	if err != nil {
// 	   return err
// 	}

// 	// Carry out another transaction in exactly the same way.
// 	_, err = tx.Exec("UPDATE ...")
// 	if err != nil {
// 	  return err
// 	}

// 	err = tx.Commit()
// 	return err;
// }