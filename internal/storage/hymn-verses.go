package storage

import (
	"context"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) CreateVerse(ctx context.Context, hymnID uuid.UUID, v *entities.HymnVerse) error {
	s.logger.Debug("CreateVerse")

	_, err := s.MySqlDB.ExecContext(
		ctx,
		`INSERT INTO HymnVerses (hymn_id, verse_number, verse_lines, optional) VALUES (?, ?, ?, ?)`,
		v.HymnID,
		v.VerseNumber,
		v.VerseLines,
		v.Optional,
	)
	if err != nil {
		s.logger.Error("CreateVerse Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) DeleteVerse(ctx context.Context, hymnID uuid.UUID, verseNum int) error {
	s.logger.Debug("DeleteVerse for hymnID ", hymnID)

	result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE HymnVerses 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE HymnVerses.hymn_id = ?
		  AND verse_number = ?
		  AND deleted_dt IS NULL`,
		hymnID,
		verseNum,
	)
	if err != nil {
		s.logger.Error("DeleteVerse error ", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("DeleteVerse affected rows: ", affectedRows)
	}
	return nil
}

func (s lowerThirdsService) GetVerse(ctx context.Context, hymnID uuid.UUID, verseNum int) (*entities.HymnVerse, error) {
	s.logger.Debug("GetVerse by hymnID ", hymnID)

	var verse entities.HymnVerse

	err := s.MySqlDB.Get(
		&verse,
		`SELECT * FROM HymnVerses WHERE hymn_id = ? AND verse_number = ?`,
		hymnID,
		verseNum,
	)
	if err != nil {
		s.logger.Error("GetVerse Error", err)
		return nil, err
	}
	return &verse, nil
}

func (s lowerThirdsService) GetVerses(ctx context.Context, hymnID uuid.UUID) (*[]entities.HymnVerse, error) {
	s.logger.Debug("GetVerses by hymnID ", hymnID)

	var verses []entities.HymnVerse

	err := s.MySqlDB.Select(
		&verses,
		`SELECT * FROM HymnVerses WHERE hymn_id = ? ORDER BY verse_number`,
		hymnID,
	)
	if err != nil {
		s.logger.Error("GetVerses Error", err)
		return nil, err
	}
	return &verses, nil
}

func (s lowerThirdsService) UpdateVerse(ctx context.Context, hymnID uuid.UUID, verseNum int, v *entities.HymnVerse) error {
	s.logger.Debug("UpdateVerse")

	_, err := s.MySqlDB.ExecContext(
		ctx,
		`UPDATE HymnVerses SET 
		  hymn_id = ?, 
		  verse_number = ?,
		  verse_lines = ?,
		  optional = ?
		WHERE hymn_id = ? 
		  AND verse_number = ?`,
		v.HymnID,
		v.VerseNumber,
		v.VerseLines,
		v.Optional,
		hymnID,
		verseNum,
	)
	if err != nil {
		s.logger.Error("UpdateVerse Error", err)
		return err
	}
	return nil
}
