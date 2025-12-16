package repository

func (r *Repository) DeleteGrowthRequest_() error {
	creatorID := 1
	query := `
		UPDATE growth_requests 
		SET status = 'удалён'
		WHERE creator_id = ? AND status = 'черновик'
	`

	if err := r.db.Exec(query, creatorID).Error; err != nil {
		return err
	}

	return nil
}
