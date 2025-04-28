package storage

import (
    "context"
    "github.com/google/uuid"
    "lowerthirdsapi/internal/entities"
    "lowerthirdsapi/internal/helpers"
)

func (s lowerThirdsService) CreateHymn(ctx context.Context, h *entities.Hymn) error {
    s.logger.Debug("CreateHymn")

    // TODO: put some user-level security on this query
    _, err := s.MySqlDB.ExecContext(
        ctx,
        `INSERT INTO Hymns (id, language, page, name, translation_id) VALUES (?, ?, ?, ?, ?)`,
        h.HymnID,
        h.Language,
        h.Page,
        h.Name,
        h.TranslationID,
    )
    if err != nil {
        s.logger.Error("CreateHymn Error", err)
        return err
    }
    return nil
}

func (s lowerThirdsService) DeleteHymn(ctx context.Context, hymnID uuid.UUID) error {
    s.logger.Debug("DeleteHymn for hymnID ", hymnID)

    // TODO: put some user level security on this query
    result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE Hymns 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE id = ?
		  AND deleted_dt IS NULL`,
        hymnID,
    )
    if err != nil {
        s.logger.Error("DeleteHymn error ", err)
        return err
    }
    affectedRows, err := result.RowsAffected()
    if err == nil {
        s.logger.Info("DeleteHymn affected rows: ", affectedRows)
    }
    return nil
}

func (s lowerThirdsService) GetHymn(ctx context.Context, hymnID uuid.UUID) (*entities.Hymn, error) {
    s.logger.Debug("GetHymn by hymnID ", hymnID)

    var hymn entities.Hymn

    err := s.MySqlDB.Get(&hymn, `SELECT * FROM Hymns WHERE id = ?`, hymnID)
    if err != nil {
        s.logger.Error("GetHymn Error", err)
        return nil, err
    }

    verses, err := s.GetVerses(ctx, hymnID)
    if err != nil {
        s.logger.Error("GetHymn verses error", err)
        return nil, err
    }
    hymn.Verses = *verses

    return &hymn, nil
}

func (s lowerThirdsService) GetHymns(ctx context.Context) (*[]entities.Hymn, error) {
    s.logger.Debug("GetHymns")

    var hymn []entities.Hymn

    qp := helpers.GetQueryParams(ctx)
    err := s.MySqlDB.Select(&hymn, `SELECT * FROM Hymns WHERE language = ?`, qp.Language)
    if err != nil {
        s.logger.Error("GetHymn Error", err)
        return nil, err
    }
    return &hymn, nil
}

func (s lowerThirdsService) UpdateHymn(ctx context.Context, hymnID uuid.UUID, h *entities.Hymn) error {
    s.logger.Debug("UpdateHymn")

    // TODO: put some user level security on this query
    result, err := s.MySqlDB.ExecContext(
        ctx,
        `UPDATE Hymns SET
		  id = ?,
		  language = ?,
		  page = ?,
		  name = ?,
		  translation_id = ?
		WHERE id = ?`,
        h.HymnID,
        h.Language,
        h.Page,
        h.Name,
        h.TranslationID,
        hymnID,
    )
    if err != nil {
        s.logger.Error("UpdateHymn Error", err)
        return err
    }
    affectedRows, err := result.RowsAffected()
    if err == nil {
        s.logger.Info("UpdateHymn affected rows: ", affectedRows)
    }
    return nil
}
