# terraform-provider-expensify

## Usage

```terraform
variable "expensify_partner_user_id" {}
variable "expensify_partner_user_secret" {}

provider "expensify" {
  partner_user_id     = var.expensify_partner_user_id
  partner_user_secret = var.expensify_partner_user_secret
}

data "expensify_policy" "policy" {
  name = "My Policy Name"
}

resource "expensify_report" "test" {
  email     = "paultyng@example.com"
  title     = "April Report"
  policy_id = data.expensify_policy.policy.id

  expense {
    date         = "2020-04-20"
    merchant     = "Internet"
    amount_cents = 8972 // in cents
  }
}
```