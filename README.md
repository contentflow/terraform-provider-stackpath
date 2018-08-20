# contentflow / terraform-provider-stackpath

This provider adds ability to [Terraform](https://www.terraform.io/) to deploy SSL certificates into StackPath CDN accounts.

## Usage

```hcl
provider "stackpath" {
  # Find those values in your StackPath account
  company_alias   = "****"
  consumer_key    = "****"
  consumer_secret = "****"
}

resource "stackpath_ssl_certificate" "mycert" {
  ssl_crt      = "${file("cert.crt")}"
  ssl_key      = "${file("cert.key")}"
  ssl_cabundle = "${file("intermediate.crt")}"
}
```
