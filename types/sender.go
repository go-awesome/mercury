//
//  sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package types

const (
	ApnsSenderID = "apns"
	GcmSenderID = "gcm"
)

func IsValidSenderID(senderID string) bool {
	return (senderID == ApnsSenderID || senderID == GcmSenderID)
}
