package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"roersla.no/askeladden/internal/config"
)

type DatabaseIface interface {
	AddQuestion(question, authorID, authorName, messageID, channelID string) (int64, error)
	GetQuestionByMessageID(messageID string) (*Question, error)
	ApproveQuestion(questionID int, approverID string) error
	RejectQuestion(questionID int, rejectorID string) error
	GetPendingQuestion() (*Question, error)
	UpdateApprovalMessageID(questionID int, approvalMessageID string) error
	GetQuestionByApprovalMessageID(approvalMessageID string) (*Question, error)
	GetPendingQuestionByID(questionID int) (*Question, error)
	GetApprovalStats() (int, int, int, error)
	GetLeastAskedApprovedQuestion() (*Question, error)
	IncrementQuestionUsage(questionID int) error
	GetApprovedQuestionStats() (int, int, int, error)
	AddBannedWord(word, reason, authorID string) error
	RemoveBannedWord(word string) error
	IsBannedWord(word string) (bool, *BannedWord, error)
	GetBannedWords() ([]*BannedWord, error)
	Close() error
	ClearDatabase() error
}

// DB struct implements the DatabaseIface
var _ DatabaseIface = (*DB)(nil)

type DB struct {
	conn              *sql.DB
	tableName         string // Dynamic table name (daily_questions or daily_questions_testing)
	bannedWordsTable  string // banned_bokmal_words or banned_bokmal_words_testing
}

// New creates a new database connection
func New(cfg *config.Config) (*DB, error) {
	log.Printf("Connecting to database at %s:%d", cfg.Database.Host, cfg.Database.Port)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully")
	
	// Determine table names based on config
	tableName := "daily_questions"
	bannedWordsTable := "banned_bokmal_words"
	
	if cfg.TableSuffix != "" {
		tableName += cfg.TableSuffix
		bannedWordsTable += cfg.TableSuffix
		log.Printf("Using beta table names: %s, %s", tableName, bannedWordsTable)
	}
	
	db := &DB{
		conn:              conn,
		tableName:         tableName,
		bannedWordsTable: bannedWordsTable,
	}
	
	// Create tables if they don't exist
	log.Println("Creating database tables if they don't exist")
	if err := db.createTables(); err != nil {
		log.Printf("Failed to create tables: %v", err)
		return nil, err
	}

	// Create banned Bokmål words table
	bannedWordsQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INT AUTO_INCREMENT PRIMARY KEY,
		word VARCHAR(255) NOT NULL UNIQUE,
		reason TEXT,
		author_id VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`, db.bannedWordsTable)

	log.Printf("Creating table: %s", db.bannedWordsTable)
	if _, err := db.conn.Exec(bannedWordsQuery); err != nil {
		return nil, fmt.Errorf("failed to create %s table: %w", db.bannedWordsTable, err)
	}

	log.Println("Database initialization completed")
	return db, nil
}

// createTables creates the necessary database tables
func (db *DB) createTables() error {
	// Create questions table with approval column and usage tracking
	// Use dynamic table name (daily_questions or daily_questions_testing)
	questionsQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INT AUTO_INCREMENT PRIMARY KEY,
		question TEXT NOT NULL,
		author_id VARCHAR(255) NOT NULL,
		author_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		times_asked INT DEFAULT 0,
		last_asked_at TIMESTAMP NULL,
		message_id VARCHAR(255),
		channel_id VARCHAR(255),
		approval_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
		approval_message_id VARCHAR(255),
		approved_by VARCHAR(255),
		approved_at TIMESTAMP NULL
	);`, db.tableName)

	log.Printf("Creating table: %s", db.tableName)
	if _, err := db.conn.Exec(questionsQuery); err != nil {
		return fmt.Errorf("failed to create %s table: %w", db.tableName, err)
	}

	return nil
}

// Question represents a question from the database
type Question struct {
	ID               int
	Question         string
	AuthorID         string
	AuthorName       string
	CreatedAt        time.Time
	TimesAsked       int
	LastAskedAt      *time.Time
	MessageID        string
	ChannelID        string
	ApprovalStatus   string
	ApprovalMessageID *string
	ApprovedBy       *string
	ApprovedAt       *time.Time
}

// BannedWord represents a banned word from the database
type BannedWord struct {
	ID        int
	Word      string
	Reason    string
	AuthorID  string
	CreatedAt time.Time
}

// AddQuestion adds a new question to the database
func (db *DB) AddQuestion(question, authorID, authorName, messageID, channelID string) (int64, error) {
	log.Printf("Adding question from user %s (ID: %s): %s", authorName, authorID, question)
	query := fmt.Sprintf("INSERT INTO %s (question, author_id, author_name, message_id, channel_id) VALUES (?, ?, ?, ?, ?)", db.tableName)
	result, err := db.conn.Exec(query, question, authorID, authorName, messageID, channelID)
	if err != nil {
		log.Printf("Failed to add question: %v", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		return 0, err
	}
	log.Printf("Successfully added question with ID %d", id)
	return id, nil
}

// AddQuestionFromMessage adds a new question to the database from a message
func (db *DB) AddQuestionFromMessage(message *discordgo.Message) (int64, error) {
	return db.AddQuestion(message.Content, message.Author.ID, message.Author.Username, message.ID, message.ChannelID)
}

// GetQuestionByMessageID gets a question by its Discord message ID
func (db *DB) GetQuestionByMessageID(messageID string) (*Question, error) {
	query := fmt.Sprintf("SELECT id, question, author_id, author_name, created_at, times_asked, last_asked_at, message_id, channel_id, approval_status, approval_message_id, approved_by, approved_at FROM %s WHERE message_id = ?", db.tableName)
	var q Question
	err := db.conn.QueryRow(query, messageID).Scan(
		&q.ID, &q.Question, &q.AuthorID, &q.AuthorName, &q.CreatedAt, &q.TimesAsked, &q.LastAskedAt, &q.MessageID, &q.ChannelID,
		&q.ApprovalStatus, &q.ApprovalMessageID, &q.ApprovedBy, &q.ApprovedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// ApproveQuestion updates the approval status for a question
func (db *DB) ApproveQuestion(questionID int, approverID string) error {
    log.Printf("Approving question ID %d by approver %s", questionID, approverID)
    query := fmt.Sprintf("UPDATE %s SET approval_status = 'approved', approved_by = ?, approved_at = NOW() WHERE id = ?", db.tableName)
    _, err := db.conn.Exec(query, approverID, questionID)
    if err != nil {
        log.Printf("Failed to approve question ID %d: %v", questionID, err)
        return err
    }
    log.Printf("Successfully approved question ID %d", questionID)
    return nil
}

// RejectQuestion updates the approval status for a question to rejected
func (db *DB) RejectQuestion(questionID int, rejectorID string) error {
    log.Printf("Rejecting question ID %d by rejector %s", questionID, rejectorID)
    query := fmt.Sprintf("UPDATE %s SET approval_status = 'rejected', approved_by = ?, approved_at = NOW() WHERE id = ?", db.tableName)
    _, err := db.conn.Exec(query, rejectorID, questionID)
    if err != nil {
        log.Printf("Failed to reject question ID %d: %v", questionID, err)
        return err
    }
    log.Printf("Successfully rejected question ID %d", questionID)
    return nil
}

// GetPendingQuestion retrieves the next pending question for approval
func (db *DB) GetPendingQuestion() (*Question, error) {
    log.Println("Retrieving next pending question")
    query := fmt.Sprintf("SELECT id, question, author_id, author_name, created_at, times_asked, last_asked_at, message_id, channel_id, approval_status, approval_message_id, approved_by, approved_at FROM %s WHERE approval_status = 'pending' ORDER BY created_at ASC LIMIT 1", db.tableName)
    var q Question
    err := db.conn.QueryRow(query).Scan(&q.ID, &q.Question, &q.AuthorID, &q.AuthorName, &q.CreatedAt, &q.TimesAsked, &q.LastAskedAt, &q.MessageID, &q.ChannelID, &q.ApprovalStatus, &q.ApprovalMessageID, &q.ApprovedBy, &q.ApprovedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("No pending questions found")
            return nil, nil
        }
        log.Printf("Failed to get pending question: %v", err)
        return nil, err
    }
    log.Printf("Retrieved pending question ID %d: %s", q.ID, q.Question)
    return &q, nil
}

// UpdateApprovalMessageID updates the approval message ID for a question
func (db *DB) UpdateApprovalMessageID(questionID int, approvalMessageID string) error {
	log.Printf("Updating approval message ID for question %d: %s", questionID, approvalMessageID)
	query := fmt.Sprintf("UPDATE %s SET approval_message_id = ? WHERE id = ?", db.tableName)
	_, err := db.conn.Exec(query, approvalMessageID, questionID)
	if err != nil {
		log.Printf("Failed to update approval message ID for question %d: %v", questionID, err)
		return err
	}
	log.Printf("Successfully updated approval message ID for question %d", questionID)
	return nil
}

// GetQuestionByApprovalMessageID gets a question by its approval message ID
func (db *DB) GetQuestionByApprovalMessageID(approvalMessageID string) (*Question, error) {
	log.Printf("Looking up question by approval message ID: %s", approvalMessageID)
	query := fmt.Sprintf("SELECT id, question, author_id, author_name, created_at, times_asked, last_asked_at, message_id, channel_id, approval_status, approval_message_id, approved_by, approved_at FROM %s WHERE approval_message_id = ?", db.tableName)
	var q Question
	err := db.conn.QueryRow(query, approvalMessageID).Scan(
		&q.ID, &q.Question, &q.AuthorID, &q.AuthorName, &q.CreatedAt, &q.TimesAsked, &q.LastAskedAt, &q.MessageID, &q.ChannelID,
		&q.ApprovalStatus, &q.ApprovalMessageID, &q.ApprovedBy, &q.ApprovedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No question found for approval message ID: %s", approvalMessageID)
		} else {
			log.Printf("Failed to get question by approval message ID %s: %v", approvalMessageID, err)
		}
		return nil, err
	}
	log.Printf("Found question ID %d for approval message %s", q.ID, approvalMessageID)
	return &q, nil
}

// GetPendingQuestionByID gets a pending question by its question ID
func (db *DB) GetPendingQuestionByID(questionID int) (*Question, error) {
	log.Printf("Looking up pending question by question ID: %d", questionID)
	query := fmt.Sprintf("SELECT id, question, author_id, author_name, created_at, times_asked, last_asked_at, message_id, channel_id, approval_status, approval_message_id, approved_by, approved_at FROM %s WHERE id = ? AND approval_status = 'pending'", db.tableName)
	log.Printf("[DEBUG] SQL Query: %s", query)
	var q Question
	err := db.conn.QueryRow(query, questionID).Scan(
		&q.ID, &q.Question, &q.AuthorID, &q.AuthorName, &q.CreatedAt, &q.TimesAsked, &q.LastAskedAt, &q.MessageID, &q.ChannelID,
		&q.ApprovalStatus, &q.ApprovalMessageID, &q.ApprovedBy, &q.ApprovedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No pending question found with ID: %d", questionID)
		} else {
			log.Printf("Failed to get pending question by ID %d: %v", questionID, err)
		}
		return nil, err
	}
	log.Printf("[DATABASE] Found pending question ID %d", q.ID)
	return &q, nil
}


// GetApprovalStats returns statistics about question approvals
func (db *DB) GetApprovalStats() (int, int, int, error) {
	var pending, approved, rejected int
	
	// Get pending count
	pendingQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE approval_status = 'pending'", db.tableName)
	err := db.conn.QueryRow(pendingQuery).Scan(&pending)
	if err != nil {
		return 0, 0, 0, err
	}
	
	// Get approved count
	approvedQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE approval_status = 'approved'", db.tableName)
	err = db.conn.QueryRow(approvedQuery).Scan(&approved)
	if err != nil {
		return 0, 0, 0, err
	}
	
	// Get rejected count
	rejectedQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE approval_status = 'rejected'", db.tableName)
	err = db.conn.QueryRow(rejectedQuery).Scan(&rejected)
	if err != nil {
		return 0, 0, 0, err
	}
	
	return pending, approved, rejected, nil
}

// GetLeastAskedApprovedQuestion gets the least asked approved question for equal distribution
func (db *DB) GetLeastAskedApprovedQuestion() (*Question, error) {
	log.Println("Retrieving least asked approved question")
	query := fmt.Sprintf(`SELECT id, question, author_id, author_name, created_at, times_asked, last_asked_at, message_id, channel_id, approval_status, approval_message_id, approved_by, approved_at 
		  FROM %s 
		  WHERE approval_status = 'approved' 
		  ORDER BY times_asked ASC, created_at ASC 
		  LIMIT 1`, db.tableName)
	var q Question
	err := db.conn.QueryRow(query).Scan(
		&q.ID, &q.Question, &q.AuthorID, &q.AuthorName, &q.CreatedAt, &q.TimesAsked, &q.LastAskedAt,
		&q.MessageID, &q.ChannelID, &q.ApprovalStatus, &q.ApprovalMessageID, &q.ApprovedBy, &q.ApprovedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[DATABASE] No approved questions found")
			return nil, nil
		} else {
			log.Printf("[DATABASE] Failed to get least asked approved question: %v", err)
			return nil, err
		}
	}
	log.Printf("[DATABASE] Retrieved least asked approved question (asked %d times): %s", q.TimesAsked, q.Question)
	return &q, nil
}

// IncrementQuestionUsage increments the times_asked count and updates last_asked_at for a question
func (db *DB) IncrementQuestionUsage(questionID int) error {
	log.Printf("[DATABASE] Incrementing usage count for question ID %d", questionID)
	query := fmt.Sprintf("UPDATE %s SET times_asked = times_asked + 1, last_asked_at = NOW() WHERE id = ?", db.tableName)
	_, err := db.conn.Exec(query, questionID)
	if err != nil {
		log.Printf("[DATABASE] Failed to increment usage count for question ID %d: %v", questionID, err)
		return err
	}
	log.Printf("[DATABASE] Successfully incremented usage count for question ID %d", questionID)
	return nil
}

// GetApprovedQuestionStats returns stats about approved questions usage
func (db *DB) GetApprovedQuestionStats() (int, int, int, error) {
	var totalApproved, totalAsked, minAsked int
	
	// Get total approved questions count
	totalApprovedQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE approval_status = 'approved'", db.tableName)
	err := db.conn.QueryRow(totalApprovedQuery).Scan(&totalApproved)
	if err != nil {
		return 0, 0, 0, err
	}
	
	// Get total times questions have been asked
	totalAskedQuery := fmt.Sprintf("SELECT COALESCE(SUM(times_asked), 0) FROM %s WHERE approval_status = 'approved'", db.tableName)
	err = db.conn.QueryRow(totalAskedQuery).Scan(&totalAsked)
	if err != nil {
		return 0, 0, 0, err
	}
	
	// Get minimum times asked (for equal distribution tracking)
	minAskedQuery := fmt.Sprintf("SELECT COALESCE(MIN(times_asked), 0) FROM %s WHERE approval_status = 'approved'", db.tableName)
	err = db.conn.QueryRow(minAskedQuery).Scan(&minAsked)
	if err != nil {
		return 0, 0, 0, err
	}
	
	return totalApproved, totalAsked, minAsked, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// AddBannedWord adds a new banned word to the database
func (db *DB) AddBannedWord(word, reason, authorID string) error {
	log.Printf("Adding banned word: %s by %s", word, authorID)
	query := fmt.Sprintf("INSERT INTO %s (word, reason, author_id) VALUES (?, ?, ?)", db.bannedWordsTable)
	_, err := db.conn.Exec(query, word, reason, authorID)
	if err != nil {
		log.Printf("Failed to add banned word: %v", err)
		return err
	}
	log.Printf("Successfully added banned word: %s", word)
	return nil
}

// RemoveBannedWord removes a banned word from the database
func (db *DB) RemoveBannedWord(word string) error {
	log.Printf("Removing banned word: %s", word)
	query := fmt.Sprintf("DELETE FROM %s WHERE word = ?", db.bannedWordsTable)
	_, err := db.conn.Exec(query, word)
	if err != nil {
		log.Printf("Failed to remove banned word: %v", err)
		return err
	}
	log.Printf("Successfully removed banned word: %s", word)
	return nil
}

// IsBannedWord checks if a word is banned
func (db *DB) IsBannedWord(word string) (bool, *BannedWord, error) {
	query := fmt.Sprintf("SELECT id, word, reason, author_id, created_at FROM %s WHERE word = ?", db.bannedWordsTable)
	var bw BannedWord
	err := db.conn.QueryRow(query, word).Scan(
		&bw.ID, &bw.Word, &bw.Reason, &bw.AuthorID, &bw.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &bw, nil
}

// GetBannedWords returns all banned words
func (db *DB) GetBannedWords() ([]*BannedWord, error) {
	query := fmt.Sprintf("SELECT id, word, reason, author_id, created_at FROM %s ORDER BY created_at DESC", db.bannedWordsTable)
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*BannedWord
	for rows.Next() {
		var bw BannedWord
		err := rows.Scan(&bw.ID, &bw.Word, &bw.Reason, &bw.AuthorID, &bw.CreatedAt)
		if err != nil {
			return nil, err
		}
		words = append(words, &bw)
	}
	return words, nil
}

// ClearDatabase drops all tables from the database
func (db *DB) ClearDatabase() error {
	log.Println("Clearing the database")
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", db.tableName)
	_, err := db.conn.Exec(query)
	if err != nil {
		log.Printf("Failed to clear the database: %v", err)
		return err
	}
	log.Println("Database cleared successfully")
	return nil
}
