package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
)

type AdminHtml struct {
	User      authenticate.User // To store user details
	ClientMsg string            // To store client message

	MapCategory      map[int]database.Category      // To store category details
	MapMemberType    map[int]database.MemberType    // To store member type details
	MapRequestStatus map[int]database.RequestStatus // To store request status details
}

var Selection = ""
