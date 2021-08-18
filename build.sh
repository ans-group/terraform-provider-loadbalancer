go build -o terraform-provider-loadbalancer
export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
mkdir -p ~/.terraform.d/plugins/ukfast.io/test/loadbalancer/0.1/$OS_ARCH
mv terraform-provider-loadbalancer ~/.terraform.d/plugins/ukfast.io/test/loadbalancer/0.1/$OS_ARCH