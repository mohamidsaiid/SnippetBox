package mysql

import (
	"database/sql"
	"mohamidsaiid.com/snippetbox/pkg/models"
)

// define a snippet model to wrap the sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// to insert a new snippet into the database and return the id or error if any error occuerd
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// using the exec method to execute the stmt followed by the needed data and returns sql.results
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil{
		return 0, err
	}

	id, err :=  result.LastInsertId()
	if err != nil{
		return 0, err
	}
	return int(id), nil
}

// to retrive the specifed snippet throught its id
// this method return struct the Snippet contains all the date or and error
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	
	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows{
		return nil, models.ErrNoRecord
	} else if err != nil{
		return nil, err
	}

	return s, nil
}

// this method returns the latest 10 snippets or an error if any error occured
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	
	rows, err := m.DB.Query(stmt)
	if err != nil{
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next(){

		snippet := &models.Snippet{}

		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil{
			return nil, err
		}
		
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return snippets, nil
}
