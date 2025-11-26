package repository

import (
	"context"
	"database/sql"
	"log"

	"golang_daerah/config"
	"golang_daerah/internal/entities"

	_ "github.com/denisenkom/go-mssqldb"
)

type SQLServerTrafficTicketRepository struct {
	db *sql.DB
}

func NewSQLServerTrafficTicketRepository(db *sql.DB) *SQLServerTrafficTicketRepository {
	return &SQLServerTrafficTicketRepository{db: db}
}

func (r *SQLServerTrafficTicketRepository) Insert(ticket *entities.TrafficTicket) error {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `
		INSERT INTO traffic_tickets (
			detected_speed, legal_speed, violation_location, violation_date, violation_time,
			violation_type, license_plate_number, vehicle_production_id, vehicle_factory, vehicle_model,
			vehicle_color, vehicle_brand, officer_name, officer_id, officer_rank, suspect_name,
			suspect_id, suspect_age, officer_age, suspect_job, suspect_address, suspect_birth_place,
			officer_branch_office_address
		)
		VALUES (
			@detected_speed, @legal_speed, @violation_location, @violation_date, @violation_time,
			@violation_type, @license_plate_number, @vehicle_production_id, @vehicle_factory, @vehicle_model,
			@vehicle_color, @vehicle_brand, @officer_name, @officer_id, @officer_rank, @suspect_name,
			@suspect_id, @suspect_age, @officer_age, @suspect_job, @suspect_address, @suspect_birth_place,
			@officer_branch_office_address
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		sql.Named("detected_speed", ticket.DetectedSpeed),
		sql.Named("legal_speed", ticket.LegalSpeed),
		sql.Named("violation_location", ticket.ViolationLocation),
		sql.Named("violation_date", ticket.ViolationDate),
		sql.Named("violation_time", ticket.ViolationTime),
		sql.Named("violation_type", ticket.ViolationType),
		sql.Named("license_plate_number", ticket.LicensePlateNumber),
		sql.Named("vehicle_production_id", ticket.VehicleProductionID),
		sql.Named("vehicle_factory", ticket.VehicleFactory),
		sql.Named("vehicle_model", ticket.VehicleModel),
		sql.Named("vehicle_color", ticket.VehicleColor),
		sql.Named("vehicle_brand", ticket.VehicleBrand),
		sql.Named("officer_name", ticket.OfficerName),
		sql.Named("officer_id", ticket.OfficerID),
		sql.Named("officer_rank", ticket.OfficerRank),
		sql.Named("suspect_name", ticket.SuspectName),
		sql.Named("suspect_id", ticket.SuspectID),
		sql.Named("suspect_age", ticket.SuspectAge),
		sql.Named("officer_age", ticket.OfficerAge),
		sql.Named("suspect_job", ticket.SuspectJob),
		sql.Named("suspect_address", ticket.SuspectAddress),
		sql.Named("suspect_birth_place", ticket.SuspectBirthPlace),
		sql.Named("officer_branch_office_address", ticket.OfficerBranchOfficeAddress),
	)

	if err != nil {
		log.Printf("[SQL Server] InsertTrafficTicket error: %v", err)
		return handleQueryError(err)
	}

	return err
}

func (r *SQLServerTrafficTicketRepository) List(limit, offset int) ([]*entities.TrafficTicket, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	limit, offset = normalizePagination(limit, offset)

	query := `
		SELECT id, detected_speed, legal_speed, violation_location, violation_date, violation_time,
		       violation_type, license_plate_number, vehicle_production_id, vehicle_factory, vehicle_model,
		       vehicle_color, vehicle_brand, officer_name, officer_id, officer_rank, suspect_name,
		       suspect_id, suspect_age, officer_age, suspect_job, suspect_address, suspect_birth_place,
		       officer_branch_office_address
		FROM traffic_tickets
		ORDER BY id DESC
		OFFSET @offset ROWS
		FETCH NEXT @limit ROWS ONLY
	`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, handleQueryError(err)
	}
	defer rows.Close()

	var tickets []*entities.TrafficTicket
	for rows.Next() {
		t := &entities.TrafficTicket{}
		if err := rows.Scan(
			&t.ID,
			&t.DetectedSpeed,
			&t.LegalSpeed,
			&t.ViolationLocation,
			&t.ViolationDate,
			&t.ViolationTime,
			&t.ViolationType,
			&t.LicensePlateNumber,
			&t.VehicleProductionID,
			&t.VehicleFactory,
			&t.VehicleModel,
			&t.VehicleColor,
			&t.VehicleBrand,
			&t.OfficerName,
			&t.OfficerID,
			&t.OfficerRank,
			&t.SuspectName,
			&t.SuspectID,
			&t.SuspectAge,
			&t.OfficerAge,
			&t.SuspectJob,
			&t.SuspectAddress,
			&t.SuspectBirthPlace,
			&t.OfficerBranchOfficeAddress,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tickets, nil
}
