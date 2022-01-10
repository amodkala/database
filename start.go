package raft

//
// Start has a few main responsibilities:
// 		- listen for requests from peers
// 		- initialize the peer as a follower
//
func (cm *CM) Start() {

	cm.becomeFollower(0)
}
