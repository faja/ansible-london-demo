consul {
  address = "consul:8500"
}

template {
  source = "/demo/consul-template/templates/ansible/roles.ctmpl"
  destination = "/demo/consul-template/templates/ansible/roles"
  command = "sh -c 'ANSIBLE_LIBRARY=/opt/ansible/ansible/library PYTHONPATH=/opt/ansible/ansible/lib /demo/differ/bin/differ -hostname node1 -type roles -old /demo/consul-template/templates/ansible/roles.old -new /demo/consul-template/templates/ansible/roles || true'"
}

template {
  source = "/demo/consul-template/templates/ansible/vars.ctmpl"
  destination = "/demo/consul-template/templates/ansible/vars"
  command = "sh -c 'ANSIBLE_LIBRARY=/opt/ansible/ansible/library PYTHONPATH=/opt/ansible/ansible/lib /demo/differ/bin/differ -hostname node1 -type vars -old /demo/consul-template/templates/ansible/vars.old -new /demo/consul-template/templates/ansible/vars || true'"
}
