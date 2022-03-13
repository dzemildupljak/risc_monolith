package docs

// swagger:route GET /user/users User ListUsers
// You receive an list of users
// responses:
//   200: userListResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:route GET /user/{user_id} User GetUserById
// You receive an list of users
// responses:
//   200: userResponse
//   500: genericResponse

/////////////////////////////////////////////////
// swagger:parameters UpdateUserById
type updateUserRequest struct {
	// in:body
	Body struct {
		Address  string `json:"addres"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}
}

// swagger:route PUT /user/{user_id} User UpdateUserById
// You receive an list of users
// responses:
//   200: userResponse
//   500: genericResponse

/////////////////////////////////////////////////
// swagger:route GET /user/current User CurrentUser
// You receive an user with id from JWT
// responses:
//   200: userResponse
//   500: genericResponse

/////////////////////////////////////////////////
