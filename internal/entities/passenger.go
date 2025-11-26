package entities

type Passenger struct {
    ID                        int    `json:"id"`
    PassengerName             string `json:"passenger_name"`
    PassengerID               string `json:"passenger_id"`
    Age                       int    `json:"age"`
    Gender                    string `json:"gender"`
    PassportNumber            string `json:"passport_number"`
    Nationality               string `json:"nationality"`
    FlightNumber              string `json:"flight_number"`
    DepartureAirport          string `json:"departure_airport"`
    ArrivalAirport            string `json:"arrival_airport"`
    DepartureDate             string `json:"departure_date"`
    DepartureTime             string `json:"departure_time"`
    ArrivalTime               string `json:"arrival_time"`
    SeatNumber                string `json:"seat_number"`
    TicketClass               string `json:"ticket_class"`
    BaggageWeight             int    `json:"baggage_weight"`
    Airline                   string `json:"airline"`
    Gate                      string `json:"gate"`
    BoardingStatus            string `json:"boarding_status"`
    OfficerName               string `json:"officer_name"`
    OfficerID                 string `json:"officer_id"`
    OfficerRank               string `json:"officer_rank"`
    OfficerBranchOfficeAddress string `json:"officer_branch_office_address"`
    CheckinCounter            string `json:"checkin_counter"`
    SpecialRequest            string `json:"special_request"`
}


