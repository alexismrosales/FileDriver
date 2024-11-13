package protocol

const (
	MsgSyn      = "SYN"
	MsgAck      = "ACK"
	MsgNak      = "NAK"
	MsgDownload = "DOWNLOAD"
	MsgUpload   = "UPLOAD"
	MsgPwd      = "PWD"
	MsgMkdir    = "MKDIR"
	MsgLs       = "LS"
	MsgCd       = "CD"
	MsgMv       = "MV"
	MsgRm       = "RM"
	MsgError    = "ERROR"
)

type Message struct {
	Type        string   `json:"Type"`        // Type of message that is sent to the other side
	WindowSize  int      `json:"WindowSize"`  // Specifying the size of the window of the selective reject
	SegmentSize int      `json:"SegmentSize"` // Size of the segment of packets to be sent
	ChunkIndex  int      `json:"ChunkIndex"`  // Number of the chunk to get them in order
	Paths       []string `json:"Paths"`       // All paths useful to execute commands
	Flags       []string `json:"Flags"`
	Output      string   `json:"Output"` // Show all output executed from the server
	Error       string   `json:"Error"`
}
