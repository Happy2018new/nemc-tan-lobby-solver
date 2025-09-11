package packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
)

const (
	LessonActionStart = iota
	LessonActionComplete
	LessonActionRestart
)

// LessonProgress is a packet sent by the server to the client to inform the client of updated progress on a lesson.
// This packet only functions on the Minecraft: Education Edition version of the game.
type LessonProgress struct {
	// Identifier is the identifier of the lesson that is being progressed.
	Identifier string

	/*
		PhoenixBuilder specific changes.
		Author: Happy2018new

		Action is the action the client should perform to show progress. This is one of the constants defined above.

		This is a mistake of upstream.
	*/
	Action int32
	// Action uint8

	// Score is the score the client should use when displaying the progress.
	Score int32
}

// ID ...
func (*LessonProgress) ID() uint32 {
	return IDLessonProgress
}

func (pk *LessonProgress) Marshal(io protocol.IO) {
	// PhoenixBuilder specific changes.
	// Author: Happy2018new
	//
	// This is a mistake of upstream.
	{
		io.Varint32(&pk.Action)
		// io.Uint8(&pk.Action)
	}
	io.Varint32(&pk.Score)
	io.String(&pk.Identifier)
}
