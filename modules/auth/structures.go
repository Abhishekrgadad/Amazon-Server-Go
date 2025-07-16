package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	FullName        string             `json:"full_name" bson:"full_name" validate:"required,min=3,max=20"`
	Email           string             `json:"email" bson:"email" validate:"required,email"`
	PhoneNumber     string             `json:"phone_number" bson:"phone_number" validate:"required,e164,len=13"`
	Password        string             `json:"password" bson:"password" validate:"required,min=6"`
	ShippingAddress string             `json:"shipping_address" bson:"shipping_address"`
	Role            string             `json:"role" bson:"role" validate:"required,oneof=customer seller admin"`
}

type Seller struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	BusinessName      string             `json:"business_name" bson:"business_name" validate:"required,min=3,max=20"`
	BusinessType      string             `json:"business_type" bson:"business_type" validate:"required,oneof=individual company"`
	ContactPersonName string             `json:"contact_person_name" bson:"contact_person_name" validate:"required,min=3,max=20"`
	Email             string             `json:"email" bson:"email" validate:"required,email"`
	Password          string             `json:"password" bson:"password" validate:"required,min=6"`
	PhoneNumber       string             `json:"phone_number" bson:"phone_number" validate:"required,e164,len=13"`
	BusinessAddress   string             `json:"business_address" bson:"business_address" validate:"required"`
	TaxID             string             `json:"tax_id" bson:"tax_id" validate:"required"`
	BankAccount       string             `json:"bank_account_details" bson:"bank_account_details" validate:"required"`
	ProductCategories []string           `json:"product_categories" bson:"product_categories" validate:"required,dive,min=3"`
}

type Admin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AdminName   string             `json:"admin_name" bson:"admin_name" validate:"required,min=3,max=50"`
	AdminRole   string             `json:"admin_role" bson:"admin_role" validate:"required,oneof=admin superadmin"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number" validate:"required,e164"`
	Password    string             `json:"password" bson:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=customer seller admin"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}
type PasswordReset struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	Role        string `json:"role"`
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}
type AuthEntity interface {
	GetID() primitive.ObjectID
	GetPassword() string
}

func (u *User) GetID() primitive.ObjectID { return u.ID }
func (u *User) GetPassword() string       { return u.Password }

func (s *Seller) GetID() primitive.ObjectID { return s.ID }
func (s *Seller) GetPassword() string       { return s.Password }

func (a *Admin) GetID() primitive.ObjectID { return a.ID }
func (a *Admin) GetPassword() string       { return a.Password }
