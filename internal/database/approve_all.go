package database

import "log"

// ApproveAllPendingQuestions approves all questions with status 'pending'
func (db *DB) ApproveAllPendingQuestions(approverID string) error {
	log.Printf("Approving ALL pending questions by approver %s", approverID)
	query := `UPDATE daily_questions SET approval_status='approved', approved_by=?, approved_at=NOW() WHERE approval_status='pending'`
	_, err := db.conn.Exec(query, approverID)
	if err != nil {
		log.Printf("Failed to approve all pending questions: %v", err)
		return err
	}
	log.Printf("Successfully approved ALL pending questions!")
	return nil
}
