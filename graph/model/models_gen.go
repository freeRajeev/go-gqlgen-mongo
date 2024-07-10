// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateCustomerProfileInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Email       string `json:"email"`
}

type CustomerProfile struct {
	ID          string `json:"_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Email       string `json:"email"`
}

type DeleteCustomerResponse struct {
	DeletedCusID string `json:"deletedCusId"`
}

type Mutation struct {
}

type Query struct {
}

type UpdateCustomerProfileInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Eamil       *string `json:"eamil,omitempty"`
}
