package daitem

type Options struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	MasterCode string `json:"masterCode" mapstruct:"masterCode"`
}
