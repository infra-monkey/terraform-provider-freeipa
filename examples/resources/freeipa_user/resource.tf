resource "freeipa_user" "user-1" {
  first_name        = "Roman"
  last_name         = "Roman"
  name              = "roman"
  telephone_numbers = ["+380982555429", "2-10-11"]
  email_address     = ["roman@example.com"]
}
