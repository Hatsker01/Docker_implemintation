package postgres

import (
	"time"

	pb "github.com/Hatsker01/Docker_implemintation/user-service/genproto"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.User, error) {
	var (
		time_at = time.Now()
		ruser   = pb.User{}
	)
	insertUserQuery := `INSERT INTO users (id, 
		first_name, 
		last_name, 
		email, 
		bio, 
		status, 
		created_at, 
		phone_number,username,password) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9,$10)
	RETURNING id, first_name, last_name, email, bio, phone_number, status, created_at,username,password`
	err := r.db.QueryRow(insertUserQuery, user.Id, user.FirstName, user.LastName, user.Email, user.Bio, user.Status, time_at, pq.Array(user.PhoneNumbers), user.Username, user.Password).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		pq.Array(&ruser.PhoneNumbers),

		&ruser.Status,
		&ruser.CreatedAt,
		&ruser.Username,
		&ruser.Password,
	)
	if err != nil {
		return &pb.User{}, err
	}
	for _, value := range user.Address {
		var adrsId string
		insertAdressQuery := `INSERT INTO adress (
			user_id, city, country, district, postal_code)
			VALUES 
			($1, $2, $3, $4, $5)
			 RETURNING user_id`
		err = r.db.QueryRow(insertAdressQuery, user.Id, value.City, value.Country, value.District, value.PostalCode).Scan(&adrsId)
		if err != nil {
			return nil, err
		}
	}
	return &ruser, nil
}

func (r *userRepo) CheckField(field, value string) (bool, error) {
	var existClient int
	if field == "username" {
		row := r.db.QueryRow(`SELECT count(1) FROM users WHERE username = $1 AND deleted_at IS NULL`, field)
		if err := row.Scan(&existClient); err != nil {
			return false, err
		}
	} else if field == "email" {
		row := r.db.QueryRow(`SELECT count(1) FROM users WHERE email = $1 and deleted_at IS NULL`, value)
		if err := row.Scan(&existClient); err != nil {
			return false, err
		}

	} else {
		return false, nil
	}
	if existClient == 0 {
		return false, nil
	}
	return true, nil

}

func (r *userRepo) UpdateUser(user *pb.User) (string, error) {
	var (
		timeAt = time.Now()
		// ruser = pb.User{}
	)

	insertUserQuery := `UPDATE users SET first_name = $1, last_name = $2, email = $3, bio = $4, status =$5,
	updated_at = $6, phone_number = $7 where id = $8 and deleted_at is null`
	_, err := r.db.Query(insertUserQuery, user.FirstName, user.LastName, user.Email, user.Bio, user.Status, timeAt, pq.Array(user.PhoneNumbers), user.Id)
	if err != nil {
		return "", err
	}

	for _, value := range user.Address {
		insertAdressQuery := `UPDATE adress SET city = $1, country = $2, district = $3, postal_code = $4 
			WHERE user_id = $5`
		_, err = r.db.Exec(insertAdressQuery, value.City, value.Country, value.District, value.PostalCode, user.Id)
		if err != nil {
			return "", err
		}
	}
	return "", err
}

func (r *userRepo) GetUserById(ID string) (*pb.User, error) {
	var ruser pb.User
	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number FROM users WHERE id = $1 and deleted_at is null`
	err := r.db.QueryRow(getByIdQuery, ID).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		&ruser.Status,
		&ruser.CreatedAt,
		pq.Array(&ruser.PhoneNumbers),
	)
	if err != nil {
		return &pb.User{}, err
	}
	getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress WHERE user_id = $1`
	rows, err := r.db.Query(getByIdAdressQuery, ID)
	if err != nil {
		return nil, err
	}
	var tempUser pb.User
	for rows.Next() {
		var adressById pb.Address
		err = rows.Scan(
			&adressById.City,
			&adressById.Country,
			&adressById.District,
			&adressById.PostalCode,
		)
		if err != nil {
			return nil, err
		}
		tempUser.Address = append(tempUser.Address, &adressById)
	}
	ruser.Address = tempUser.Address
	return &ruser, nil
}

func (r *userRepo) GetAllUser() ([]*pb.User, error) {
	var ruser1 []*pb.User

	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number,username FROM users where deleted_at is null`
	rowss, err := r.db.Query(getByIdQuery)

	if err != nil {
		return nil, err
	}

	for rowss.Next() {
		var ruser pb.User
		err = rowss.Scan(
			&ruser.Id,
			&ruser.FirstName,
			&ruser.LastName,
			&ruser.Email,
			&ruser.Bio,
			&ruser.Status,
			&ruser.CreatedAt,
			pq.Array(&ruser.PhoneNumbers),
			&ruser.Username,
		)
		if err != nil {
			return nil, err
		}

		getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress where user_id=$1`
		rows, err := r.db.Query(getByIdAdressQuery, ruser.Id)

		if err != nil {
			return nil, err
		}

		//var tempUser pb.User
		for rows.Next() {
			var adressById pb.Address
			err = rows.Scan(
				&adressById.City,
				&adressById.Country,
				&adressById.District,
				&adressById.PostalCode,
			)

			if err != nil {
				return nil, err
			}

			ruser.Address = append(ruser.Address, &adressById)
		}

		ruser1 = append(ruser1, &ruser)
	}

	return ruser1, nil
}

func (r *userRepo) GetUserFromPost(userID string) (*pb.GetUserFromPostResponse, error) {
	var user *pb.GetUserFromPostResponse

	query := `SELECT first_name, last_name FROM users WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *userRepo) DeleteUser(userId string) (pb.Empty, error) {
	queryL := `UPDATE users SET deleted_at=$1 where id = $2`
	_, err := r.db.Exec(queryL, time.Now(), userId)
	if err != nil {

		return pb.Empty{}, err
	}

	return pb.Empty{}, err
}

func (r *userRepo) GetListUsers(limit, page int64) (*pb.GetUserResponse, error) {

	var (
		users []*pb.User
		count int64
	)
	offset := (page - 1) * limit

	query := `SELECT id,first_name,last_name from users  where deleted_at is null order by first_name OFFSET $1 limit $2`
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user pb.User
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	countQuery := `SELECT count(*) FROM users where deleted_at is null`
	err = r.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Users: users,
		Count: count,
	}, nil

}
func (r *userRepo) Login(email string) (*pb.User, error) {
	var ruser pb.User
	getByIdQuery := `SELECT id, first_name, last_name, email, bio, status, created_at, phone_number,password FROM users WHERE email = $1 and deleted_at is null`
	err := r.db.QueryRow(getByIdQuery, email).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
		&ruser.Email,
		&ruser.Bio,
		&ruser.Status,
		&ruser.CreatedAt,
		pq.Array(&ruser.PhoneNumbers),
		&ruser.Password,
	)
	if err != nil {
		return nil, err
	}
	getByIdAdressQuery := `SELECT city, country, district, postal_code FROM adress WHERE user_id = $1`
	rows, err := r.db.Query(getByIdAdressQuery, ruser.Id)
	if err != nil {
		return nil, err
	}
	var tempUser pb.User
	for rows.Next() {
		var adressById pb.Address
		err = rows.Scan(
			&adressById.City,
			&adressById.Country,
			&adressById.District,
			&adressById.PostalCode,
		)
		if err != nil {
			return nil, err
		}
		tempUser.Address = append(tempUser.Address, &adressById)
	}
	ruser.Address = tempUser.Address

	return &ruser, nil
}
