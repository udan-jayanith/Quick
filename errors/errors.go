package QuicErr

import (
	"github.com/udan-jayanith/Quick/varint"
)

type QuickTransportError varint.Int62

// Do not change the enum order.
const (
	NO_ERROR                  QuickTransportError = 0x00 + iota //No error
	INTERNAL_ERROR                                              //Implementation error
	CONNECTION_REFUSED                                          //Server refuses a connection
	FLOW_CONTROL_ERROR                                          //Flow control error
	STREAM_LIMIT_ERROR                                          //Too many streams opened
	STREAM_STATE_ERROR                                          //Frame received in invalid stream state
	FINAL_SIZE_ERROR                                            //Change to final size
	FRAME_ENCODING_ERROR                                        //Frame encoding error
	TRANSPORT_PARAMETER_ERROR                                   //Error in transport parameters
	CONNECTION_ID_LIMIT_ERROR                                   //Too many connection IDs received
	PROTOCOL_VIOLATION                                          //Generic protocol violation
	INVALID_TOKEN                                               //Invalid Token received
	APPLICATION_ERROR                                           //Application error
	CRYPTO_BUFFER_EXCEEDED                                      //CRYPTO data buffer overflowed
	KEY_UPDATE_ERROR                                            //Invalid packet protection update
	AEAD_LIMIT_REACHED                                          //Excessive use of packet protection keys
	NO_VIABLE_PATH                                              //No viable network path exists
	//CRYPTO_ERROR                                              //0x0100-​0x01ff //TLS alert code
)

var (
	errMap map[QuickTransportError]string = map[QuickTransportError]string{
		NO_ERROR:                  "No error",
		INTERNAL_ERROR:            "The endpoint encountered an internal error and cannot continue with the connection",
		CONNECTION_REFUSED:        "The server refused to accept a new connection.",
		FLOW_CONTROL_ERROR:        "An endpoint received more data than it permitted in its advertised data limits",
		STREAM_LIMIT_ERROR:        "An endpoint received a frame for a stream identifier that exceeded its advertised stream limit",
		STREAM_STATE_ERROR:        "An endpoint received a frame for a stream that was not in a state that permitted that frame",
		FINAL_SIZE_ERROR:          "(1) An endpoint received a STREAM frame containing data that exceeded the previously established final size, \n(2) an endpoint received a STREAM frame or a RESET_STREAM frame containing a final size that was lower than the size of stream data that was already received, or \n(3) an endpoint received a STREAM frame or a RESET_STREAM frame containing a different final size to the one already established.",
		FRAME_ENCODING_ERROR:      "An endpoint received a frame that was badly formatted",
		TRANSPORT_PARAMETER_ERROR: "An endpoint received transport parameters that were badly formatted",
		CONNECTION_ID_LIMIT_ERROR: "The number of connection IDs provided by the peer exceeds the advertised active_connection_id_limit.",
		PROTOCOL_VIOLATION:        "An endpoint detected an error with protocol compliance that was not covered by more specific error codes.",
		INVALID_TOKEN:             "A server received a client Initial that contained an invalid Token field.",
		APPLICATION_ERROR:         "The application or application protocol caused the connection to be closed.",
		CRYPTO_BUFFER_EXCEEDED:    "An endpoint has received more data in CRYPTO frames than it can buffer.",
		KEY_UPDATE_ERROR:          "An endpoint detected errors in performing key updates",
		AEAD_LIMIT_REACHED:        "An endpoint has reached the confidentiality or integrity limit for the AEAD algorithm used by the given connection.",
		NO_VIABLE_PATH:            "An endpoint has determined that the network path is incapable of supporting QUIC.",
	}
)

func (code QuickTransportError) Error() string {
	if 0x0100 <= code && 0x01ff >= code {
		return "The cryptographic handshake failed"
	}
	return errMap[code]
}
