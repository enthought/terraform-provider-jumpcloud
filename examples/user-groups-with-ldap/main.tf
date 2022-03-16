resource "jumpcloud_user_group" "test_group" {
  name = "test_group"
  attributes = {
    posix_groups = "32:testerino"
  }
}

resource "jumpcloud_user" "test_user1" {
  username   = "testuser1"
  firstname  = "test"
  lastname   = "user1"
  email      = "testuser1@testorg.org"
  enable_mfa = true
}

resource "jumpcloud_user" "test_user2" {
  username   = "testuser2"
  firstname  = "test"
  lastname   = "user2"
  email      = "testuser2@testorg.org"
  enable_mfa = false
}


resource "jumpcloud_user_group_membership" "test_membership" {
  userid  = jumpcloud_user.test_user1.id
  groupid = jumpcloud_user_group.test_group.id
}

data "jumpcloud_ldap_server" "default" {
  name="JumpCloud LDAP" 
}

resource "jumpcloud_user_group_ldap_membership" "test_group-ldap_membership" {
  usergroup_id = jumpcloud_user_group.test_group.id
  ldap_id = data.jumpcloud_ldap_server.default.ldap_id
}
