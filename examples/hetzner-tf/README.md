# Using trieres with Hetzner and Terraform

This example sets up all the infra in Hetzner using Terraform and feeds the output to trieres

## Steps

1. Create `terraform.tfvars` file with needed details. You can use the provided `terraform.tfvars.example` as a baseline.
2. `terraform init`
3. `terraform apply`
4. `terraform output -json | yq r - trieres_cluster.value > cluster.yml`
5. `trieres up`
6. Profit! :)