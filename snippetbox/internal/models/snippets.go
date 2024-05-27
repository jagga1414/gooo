package models

import ( 
	"database/sql"
	"time"
	"errors"
)
// Define a Snippet type to hold the data for an individual snippet. Notice how // the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID int 
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

//SnippetModel type which wraps a sql.DB connection pool.
// Define a
type SnippetModel struct { 
	DB *sql.DB
}
// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) { 
	stmt := `INSERT INTO snippets (title, content, created, expires)
	 VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
// Use the Exec() method on the embedded connection pool to execute the 
// statement. The first parameter is the SQL statement, followed by the 
// values for the placeholder parameters: title, content and expiry in
// that order. This method returns a sql.Result type, which contains some 
// basic information about what happened when the statement was executed. 
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil { 
		return 0, err
	}
	// Use the LastInsertId() method on the result to get the ID of our // newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
	return 0, err }
	// The ID returned has the type int64, so we convert it to an int type // before returning.
	return int(id), nil
}
// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) { 
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which 
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)
	// Initialize a new zeroed Snippet struct.
	var s Snippet
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into, 
	// and the number of arguments must be exactly the same as the number of 
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
	// If the query returns no rows, then row.Scan() will return a
	// sql.ErrNoRows error. We use the errors.Is() function check for that 
	// error specifically, and return our own ErrNoRecord error
	// instead (we'll create this in a moment).
	if errors.Is(err, sql.ErrNoRows) {
		return Snippet{}, ErrNoRecord 
	} else{
		return Snippet{}, err 
	}
}
// If everything went OK, then return the filled Snippet struct.
return s, nil
}
// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) { 
	stmt := `SELECT id, title, content, created, expires FROM snippets LIMIT 3`
// Use the Query() method on the connection pool to execute our
// SQL statement. This returns a sql.Rows resultset containing the result of // our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err 
	}
	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer // statement should come *after* you check for an error from the Query() // method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()
	var snippets []Snippet
	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the // rows.Scan() method. If iteration over all the rows completes then the // resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
	// Create a pointer to a new zeroed Snippet struct.
		var s Snippet
		// Use rows.Scan() to copy the values from each field in the row to the // new Snippet object that we created. Again, the arguments to row.Scan() // must be pointers to the place you want to copy the data into, and the // number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err 
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s) 
	}
	if err = rows.Err(); err != nil { 
		return nil, err
	}
	return snippets,nil
}