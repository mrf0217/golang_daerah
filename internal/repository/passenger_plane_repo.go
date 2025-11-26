package repository

import (
	"context"
	"database/sql"
	"log"

	"golang_daerah/config"
	"golang_daerah/internal/entities"

	_ "github.com/go-sql-driver/mysql"
)

type PassengerPlaneRepository struct {
	db *sql.DB
}

func NewPassengerPlaneRepository(db *sql.DB) *PassengerPlaneRepository {
	return &PassengerPlaneRepository{db: db}
}

func (r *PassengerPlaneRepository) Insert(p *entities.Passenger) error {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `
        INSERT INTO passenger_plane (
            passenger_name, passenger_id, age, gender, passport_number, nationality,
            flight_number, departure_airport, arrival_airport, departure_date, departure_time,
            arrival_time, seat_number, ticket_class, baggage_weight, airline, gate, boarding_status,
            officer_name, officer_id, officer_rank, officer_branch_office_address, checkin_counter, special_request
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err := r.db.ExecContext(ctx, query,
		p.PassengerName,
		p.PassengerID,
		p.Age,
		p.Gender,
		p.PassportNumber,
		p.Nationality,
		p.FlightNumber,
		p.DepartureAirport,
		p.ArrivalAirport,
		p.DepartureDate,
		p.DepartureTime,
		p.ArrivalTime,
		p.SeatNumber,
		p.TicketClass,
		p.BaggageWeight,
		p.Airline,
		p.Gate,
		p.BoardingStatus,
		p.OfficerName,
		p.OfficerID,
		p.OfficerRank,
		p.OfficerBranchOfficeAddress,
		p.CheckinCounter,
		p.SpecialRequest,
	)

	if err != nil {
		log.Printf("[Passenger Plane MySQL] Insert error: %v", err)
		return handleQueryError(err)
	}

	return err
}

func (r *PassengerPlaneRepository) List(limit, offset int) ([]*entities.Passenger, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	limit, offset = normalizePagination(limit, offset)

	query := `
        SELECT id, passenger_name, passenger_id, age, gender, passport_number, nationality,
               flight_number, departure_airport, arrival_airport, departure_date, departure_time,
               arrival_time, seat_number, ticket_class, baggage_weight, airline, gate, boarding_status,
               officer_name, officer_id, officer_rank, officer_branch_office_address, checkin_counter, special_request
        FROM passenger_plane
        ORDER BY id DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, handleQueryError(err)
	}
	defer rows.Close()

	var tickets []*entities.Passenger
	for rows.Next() {
		t := &entities.Passenger{}
		if err := rows.Scan(
			&t.ID,
			&t.PassengerName,
			&t.PassengerID,
			&t.Age,
			&t.Gender,
			&t.PassportNumber,
			&t.Nationality,
			&t.FlightNumber,
			&t.DepartureAirport,
			&t.ArrivalAirport,
			&t.DepartureDate,
			&t.DepartureTime,
			&t.ArrivalTime,
			&t.SeatNumber,
			&t.TicketClass,
			&t.BaggageWeight,
			&t.Airline,
			&t.Gate,
			&t.BoardingStatus,
			&t.OfficerName,
			&t.OfficerID,
			&t.OfficerRank,
			&t.OfficerBranchOfficeAddress,
			&t.CheckinCounter,
			&t.SpecialRequest,
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

// GetPaginated returns passenger records with pagination controls, ordered by id ASC (starting from ID 1)
func (r *PassengerPlaneRepository) GetPaginated(limit, offset int) ([]entities.Passenger, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `
        SELECT id, passenger_name, passenger_id, age, gender, passport_number, nationality,
               flight_number, departure_airport, arrival_airport, departure_date, departure_time,
               arrival_time, seat_number, ticket_class, baggage_weight, airline, gate, boarding_status,
               officer_name, officer_id, officer_rank, officer_branch_office_address, checkin_counter, special_request
        FROM passenger_plane
        ORDER BY id
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, handleQueryError(err)
	}
	defer rows.Close()

	var passengers []entities.Passenger
	for rows.Next() {
		var p entities.Passenger
		if err := rows.Scan(
			&p.ID,
			&p.PassengerName,
			&p.PassengerID,
			&p.Age,
			&p.Gender,
			&p.PassportNumber,
			&p.Nationality,
			&p.FlightNumber,
			&p.DepartureAirport,
			&p.ArrivalAirport,
			&p.DepartureDate,
			&p.DepartureTime,
			&p.ArrivalTime,
			&p.SeatNumber,
			&p.TicketClass,
			&p.BaggageWeight,
			&p.Airline,
			&p.Gate,
			&p.BoardingStatus,
			&p.OfficerName,
			&p.OfficerID,
			&p.OfficerRank,
			&p.OfficerBranchOfficeAddress,
			&p.CheckinCounter,
			&p.SpecialRequest,
		); err != nil {
			return nil, err
		}
		passengers = append(passengers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return passengers, nil
}
