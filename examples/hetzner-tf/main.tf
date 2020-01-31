variable "hcloud_token" {
    description = "Hetzner API token"
}

provider "hcloud" {
  token = "${var.hcloud_token}"
}

variable "ssh_keys" {
    default = []
}

variable "ssh_user" {
    default = "root"
}

variable "cluster_name" {
    default = "trieres"
}

variable "location" {
    default = "hel1"
}

variable "image" {
    default = "ubuntu-18.04"
}

variable "master_type" {
    default = "cx11"
}

variable "worker_count" {
    default = 1
}

variable "worker_type" {
    default = "cx11"
}

resource "hcloud_server" "master" {
    name = "${var.cluster_name}-master"
    image = "${var.image}"
    server_type = "${var.master_type}"
    ssh_keys = "${var.ssh_keys}"
    location = "${var.location}"
    labels = {
        role = "master"
    }
}

resource "hcloud_server" "worker" {
    count = "${var.worker_count}"
    name = "${var.cluster_name}-worker-${count.index}"
    image = "${var.image}"
    server_type = "${var.worker_type}"
    ssh_keys = "${var.ssh_keys}"
    location = "${var.location}"
    labels = {
        role = "worker"
    }
}

output "trieres_cluster" {
    value = {
        hosts = [
            for host in concat([hcloud_server.master], hcloud_server.worker)  : {
                address      = host.ipv4_address
                user    = "root"
                role    = host.labels.role
            }
        ]
    }

}