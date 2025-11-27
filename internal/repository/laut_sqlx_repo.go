package repository

import (
    "context"
    "encoding/json"
    "golang_daerah/config"
    "github.com/jmoiron/sqlx"
)

type LautSQLXRepository struct {
    db *sqlx.DB
}

func NewLautSQLXRepository(db *sqlx.DB) *LautSQLXRepository {
    return &LautSQLXRepository{db: db}
}

func (r *LautSQLXRepository) GetPaginatedJSON(limit, offset int) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    query := `
        SELECT id, port_name, port_code, port_address, city, province, country,
               operator_name, operator_contact, harbor_master_name, harbor_master_id,
               harbor_master_rank, harbor_master_office_address, number_of_piers,
               main_pier_length, max_ship_draft, max_ship_length,
               terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
               emergency_contact, security_office_name, security_officer_id,
               security_level, checkin_counter_count, special_facilities
        FROM Laut
        ORDER BY id ASC
        LIMIT ? OFFSET ?
    `

    rows, err := r.db.QueryxContext(ctx, query, limit, offset)
    if err != nil {
        return nil, handleQueryError(err)
    }
    defer rows.Close()

    var results []map[string]interface{}
    for rows.Next() {
        row := make(map[string]interface{})
        if err := rows.MapScan(row); err != nil {
            return nil, err
        }
        results = append(results, row)
    }

    return json.Marshal(results)
}

func (r *LautSQLXRepository) InsertJSON(jsonData []byte) error {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    var items []map[string]interface{}
    if err := json.Unmarshal(jsonData, &items); err != nil {
        return err
    }

    query := `
        INSERT INTO Laut (
            port_name, port_code, port_address, city, province, country,
            operator_name, operator_contact, harbor_master_name, harbor_master_id,
            harbor_master_rank, harbor_master_office_address, number_of_piers,
            main_pier_length, max_ship_draft, max_ship_length,
            terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
            emergency_contact, security_office_name, security_officer_id,
            security_level, checkin_counter_count, special_facilities
        ) VALUES (
            :port_name, :port_code, :port_address, :city, :province, :country,
            :operator_name, :operator_contact, :harbor_master_name, :harbor_master_id,
            :harbor_master_rank, :harbor_master_office_address, :number_of_piers,
            :main_pier_length, :max_ship_draft, :max_ship_length,
            :terminal_capacity_passenger, :terminal_capacity_cargo, :operational_hours,
            :emergency_contact, :security_office_name, :security_officer_id,
            :security_level, :checkin_counter_count, :special_facilities
        )
    `

    for _, item := range items {
        if _, err := r.db.NamedExecContext(ctx, query, item); err != nil {
            return handleQueryError(err)
        }
    }

    return nil
}