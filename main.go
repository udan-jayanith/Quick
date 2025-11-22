package Quick

/*
x (A):
Indicates that x is A bits long

x (i):
Indicates that x holds an integer value using the variable-length encoding described in Section 16

x (A..B):
Indicates that x can be any length from A to B; A can be omitted to indicate a minimum of zero bits, and B can be omitted to indicate no set upper limit; values in this format always end on a byte boundary

x (L) = C:
Indicates that x has a fixed value of C; the length of x is described by L, which can use any of the length forms above

x (L) = C..D:
Indicates that x has a value in the range from C to D, inclusive, with the length described by L, as above

[x (L)]:
Indicates that x is optional and has a length of L

x (L) ...:
Indicates that x is repeated zero or more times and that each instance has a length of L
*/

//This document uses network byte order (that is, big endian) values.

var (
	notImplemented = func() {
		panic("Not implemented")
	}
)

// If a StreamFrame sends frames for the same stream with different data to offset that is already filled it should return PROTOCOL_VIOLATION.
// Sender must ensure that stream is no longer then int62 max value.

// # Features that will be added.
// * Provide a way to control stream prioritization.
