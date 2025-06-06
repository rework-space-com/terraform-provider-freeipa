# import id must be of format <name>;<zone_name>;<type>;<set_identifier>
# or without the optional set_identifier : <name>;<zone_name>;<type>

import {
    to = freeipa_dns_record.txtrecord
    id = "_kerberos;testimport.ipatest.lan;TXT;krbtxt"
}

resource "freeipa_dns_record" "txtrecord" {
  name           = "_kerberos"
  zone_name      = "testimport.ipatest.lan"
  type           = "TXT"
  records        = ["IPATEST.LAN"]
  set_identifier = "krbtxt"
}

import {
    to = freeipa_dns_record.arecord
    id = "test;testimport.ipatest.lan;A"
}

resource "freeipa_dns_record" "arecord" {
  name      = "test"
  zone_name = "testimport.ipatest.lan"
  type      = "A"
  records   = ["172.27.2.2"]
}

import {
    to = freeipa_dns_record.ptrrecord
    id = "2;2.27.172.in-addr.arpa;PTR"
}

resource "freeipa_dns_record" "ptrrecord" {
  name      = "2"
  zone_name = "2.27.172.in-addr.arpa"
  type      = "PTR"
  records   = ["test.testimport.ipatest.lan."]
}