package repository

import (
	"context"
	"database/sql"
	"golang_daerah/config"
	lautgen "golang_daerah/internal/repository/generated/laut"
	"log"
)

type LautRepository struct {
	queries *lautgen.Queries
}

func NewLautRepository(db *sql.DB) *LautRepository {
	return &LautRepository{
		queries: lautgen.New(db),
	}
}

// Insert adds a new port/terminal record to the database using SQLC-generated queries.
func (r *LautRepository) Insert(laut *lautgen.Laut) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	params := lautgen.InsertLautParams{
		PortName:                  laut.PortName,
		PortCode:                  laut.PortCode,
		PortAddress:               laut.PortAddress,
		City:                      laut.City,
		Province:                  laut.Province,
		Country:                   laut.Country,
		OperatorName:              laut.OperatorName,
		OperatorContact:           laut.OperatorContact,
		HarborMasterName:          laut.HarborMasterName,
		HarborMasterID:            laut.HarborMasterID,
		HarborMasterRank:          laut.HarborMasterRank,
		HarborMasterOfficeAddress: laut.HarborMasterOfficeAddress,
		NumberOfPiers:             laut.NumberOfPiers,
		MainPierLength:            laut.MainPierLength,
		MaxShipDraft:              laut.MaxShipDraft,
		MaxShipLength:             laut.MaxShipLength,
		TerminalCapacityPassenger: laut.TerminalCapacityPassenger,
		TerminalCapacityCargo:     laut.TerminalCapacityCargo,
		OperationalHours:          laut.OperationalHours,
		EmergencyContact:          laut.EmergencyContact,
		SecurityOfficeName:        laut.SecurityOfficeName,
		SecurityOfficerID:         laut.SecurityOfficerID,
		SecurityLevel:             laut.SecurityLevel,
		CheckinCounterCount:       laut.CheckinCounterCount,
		SpecialFacilities:         laut.SpecialFacilities,
	}

	if _, err := r.queries.InsertLaut(ctx, params); err != nil {
		log.Printf("[Laut] Insert error: %v", err)
		return handleQueryError(err)
	}

	return nil
}

// List returns port/terminal records with pagination controls (ordered by id DESC)
func (r *LautRepository) List(limit, offset int) ([]*lautgen.Laut, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	limit, offset = normalizePagination(limit, offset)

	results, err := r.queries.ListLauts(ctx, lautgen.ListLautsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, handleQueryError(err)
	}

	lauts := make([]*lautgen.Laut, len(results))
	for i := range results {
		lauts[i] = &results[i]
	}

	return lauts, nil
}

// GetPaginated returns port/terminal records with pagination controls, ordered by id ASC (starting from ID 1)
func (r *LautRepository) GetPaginated(limit, offset int) ([]lautgen.Laut, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	limit, offset = normalizePagination(limit, offset)

	results, err := r.queries.GetPaginatedLauts(ctx, lautgen.GetPaginatedLautsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, handleQueryError(err)
	}

	return results, nil
}
