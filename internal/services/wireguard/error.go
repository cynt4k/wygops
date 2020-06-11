package wireguard

import "errors"

// ErrWrongSubnet : Error if the subnet matches not the server
var ErrWrongSubnet = errors.New("wrong subnet for this ip address")

// ErrNoIPAvailable : Error if there is no ip left in the subnet
var ErrNoIPAvailable = errors.New("no ip available")
