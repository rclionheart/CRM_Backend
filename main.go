package main

type Customer struct {
	ID        int
	Name      string
	Role      string
	Email     string
	Phone     string
	Contacted bool
}

var customers = []Customer{
	{ID: 1, Name: "Bob", Role: "Member", Phone: "0678304019", Contacted: true},
	{ID: 2, Name: "Scotty", Role: "Premium Member", Phone: "0683240193", Contacted: true},
	{ID: 3, Name: "Michael", Role: "New Member", Phone: "0745320845", Contacted: false},
	{ID: 4, Name: "Andria", Role: "Member", Phone: "0836891203", Contacted: false},
}

func CreateCustomer() {

}

func GetCustomer() {

}

func UpdateCustomer() {

}

func DeleteCustomer() {

}

func main() {

}
