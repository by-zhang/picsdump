package main

import (
	"bufio"
	"net/http"
	"errors"
	"io"
	
	"github.com/google/gopacket"
    "github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

// pix handler read bytes from reader
var pixHandler *PixHandler

// HTTPStreamFactory implements tcpassembly.StreamFactory interface
type HTTPStreamFactory struct {}
// New returns a new stream for a given TCP key
func (h *HTTPStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hs := &HTTPStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hs.run()

	return &hs.r
}


type HTTPStream struct {
	net, transport gopacket.Flow
	//ReaderStream implements both tcpassembly.Stream and io.Reader
	r              tcpreader.ReaderStream
}
func (h *HTTPStream)run() {
	buf := bufio.NewReader(&h.r)
	for {
		resp, err := http.ReadResponse(buf, nil)
		if err == io.EOF {
			return
		}
		if err == nil && h.filterStream(resp) == nil{
			// Dont forget resp.Body.Close()
			pixHandler.handle(resp)
		}
	}
}
//return nil if this response has an image-type content 
func (h *HTTPStream)filterStream(resp *http.Response) error{
	if resp.StatusCode != 200 {
		return errors.New("Wrong response status.")
	}
	ct, flag := resp.Header["Content-Type"]
	if !flag || !checkImageType(ct[0]) {
		return errors.New("Wrong content type.")
	}
	return nil
}


//Assembler handles reassembling TCP streams.
func NewAssembler() (assembler *tcpassembly.Assembler){
    streamFactory := &HTTPStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler = tcpassembly.NewAssembler(streamPool)
    return
}


func checkImageType(ct string) bool {
	 _,flag := pixSet.contentType[ct]
	return flag
}


