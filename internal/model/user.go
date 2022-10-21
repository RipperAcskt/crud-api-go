package model

type User struct {
	Id       int
	Name     string
	LastName string
	Age      int
}

func (u *User) Validate() string {
	var err string
	if u.Name == "" {
		err += "Fill name."
	}
	if u.LastName == "" {
		err += "Fill surname."
	}
	if u.Age <= 0 {
		err += "Age need to be upper zero."
	}
	return err
}
