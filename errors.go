package Quick

type QuickTransportError Int62

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
