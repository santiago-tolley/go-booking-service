package clients

type AuthorizeRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type AuthorizeResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	User string `json:"user"`
	Err  error  `json:"err"`
}

type CreateRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Err error `json:"err"`
}
