package main

import(
	"time"
	"log"
	"sync"
	"os"
	"flag"
	
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
)

var counter *Counter
var fname = flag.String("f", "", "pcap. file to be read")
func init(){
	//pics type info
	pixType := map[int]*PixType{}
	pixType[JPEG] = NewPixType(JPEG, "jpeg", []byte{0xFF, 0xD8, 0xFF})
	pixType[PNG] = NewPixType(PNG, "png", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0X1A, 0X0A})
	pixSet = NewPixInfo(pixType, len(pixType))
	
	wg = sync.WaitGroup{}
	
	//dir that stores pics
	_, err := os.Stat(DIR)
	if err != nil {
		if err = os.Mkdir(DIR, os.ModePerm); err != nil {
			log.Println("make dir failed.")
		}
	}
	
	//logger
	logger = NewLogger()

	//launch handler
    pixHandler = NewPixHandler(3)
	logger.Print("handler established.")
}


func main(){
	var handle *pcap.Handle
	var tcp interface{}
	flag.Parse()
	if *fname == "" {
		logger.Fatal("please input file as an argument with -f flag")
	}
	handle, err := pcap.OpenOffline(*fname)
	if err != nil {
		logger.Fatal(err.Error())
	} 
	if err := handle.SetBPFFilter("tcp and src port 80"); err != nil {
		logger.Fatal(err.Error())
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	assembler := NewAssembler()
    ticker := time.Tick(time.Minute)
	logger.Print("start sniffing and handling.")
	ForEnd:
	for{
		select{
			case packet := <-packets:
			    if packet == nil {
			         break ForEnd
			    }
			    if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				    //Unusable
				    continue
			    }
			    tcp = packet.TransportLayer()
			    assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp.(*layers.TCP), packet.Metadata().Timestamp)
			case <-ticker:
		}
	}
	wg.Wait()
	logger.Print("task finished.")
	pixHandler.statics()
	return
}