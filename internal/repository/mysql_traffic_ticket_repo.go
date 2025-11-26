package repository

import (
	"context"
	"database/sql"
	"log"

	"golang_daerah/config"
	"golang_daerah/internal/entities"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLTrafficTicketRepository struct {
	db *sql.DB
}

func NewMySQLTrafficTicketRepository(db *sql.DB) *MySQLTrafficTicketRepository {
	return &MySQLTrafficTicketRepository{db: db}
}

func (r *MySQLTrafficTicketRepository) Insert(ticket *entities.TrafficTicket) error {
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		ticket.DetectedSpeed,
		ticket.LegalSpeed,
		ticket.ViolationLocation,
		ticket.ViolationDate,
		ticket.ViolationTime,
		ticket.ViolationType,
		ticket.LicensePlateNumber,
		ticket.VehicleProductionID,
		ticket.VehicleFactory,
		ticket.VehicleModel,
		ticket.VehicleColor,
		ticket.VehicleBrand,
		ticket.OfficerName,
		ticket.OfficerID,
		ticket.OfficerRank,
		ticket.SuspectName,
		ticket.SuspectID,
		ticket.SuspectAge,
		ticket.OfficerAge,
		ticket.SuspectJob,
		ticket.SuspectAddress,
		ticket.SuspectBirthPlace,
		ticket.OfficerBranchOfficeAddress,
	)

	if err != nil {
		log.Printf("[MySQL] InsertTrafficTicket error: %v", err)
		return handleQueryError(err)
	}

	return err
}

func (r *MySQLTrafficTicketRepository) List(limit, offset int) ([]*entities.TrafficTicket, error) {
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
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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

// GetPaginated returns tickets with pagination controls, ordered by id ASC (starting from ID 1)
func (r *MySQLTrafficTicketRepository) GetPaginated(limit, offset int) ([]entities.TrafficTicket, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `
		SELECT id, detected_speed, legal_speed, violation_location, violation_date, violation_time,
		       violation_type, license_plate_number, vehicle_production_id, vehicle_factory, vehicle_model,
		       vehicle_color, vehicle_brand, officer_name, officer_id, officer_rank, suspect_name,
		       suspect_id, suspect_age, officer_age, suspect_job, suspect_address, suspect_birth_place,
		       officer_branch_office_address
		FROM traffic_tickets
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, handleQueryError(err)
	}
	defer rows.Close()

	var tickets []entities.TrafficTicket
	for rows.Next() {
		var t entities.TrafficTicket
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
