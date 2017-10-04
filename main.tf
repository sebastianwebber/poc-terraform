provider "digitalocean" {}

resource "digitalocean_tag" "testserver" {
  name = "testserver"
}

resource "digitalocean_ssh_key" "default" {
  name       = "Terraform Example"
  public_key = "${file("/Users/sebastian/gocode/src/poc/keys/id_rsa.pub")}"
}


resource "digitalocean_droplet" "testserver01" {
    image = "26505628"
    name = "testserver-terraform"
    region = "nyc3"
    size = "512mb"
    ssh_keys = ["${digitalocean_ssh_key.default.id}"]
    tags   = ["${digitalocean_tag.testserver.id}"]

    connection {
        type = "ssh"
        user = "root"
        private_key = "${file("/Users/sebastian/gocode/src/poc/keys/id_rsa")}"
    }


    provisioner "remote-exec" {
        inline = [
            "hostname",
            "echo 'it works! bye.'"
        ]
    }

}
