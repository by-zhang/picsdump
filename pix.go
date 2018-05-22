package main

import(
	"net/http"
	"io/ioutil"
	"log"
	"time"
	"math/rand"
	"strconv"
	"sync"
	"errors"
)

const(
    JPEG = 0
    PNG = 1
	
	DIR = "./pix"
)

// pix type information
var pixSet *PixInfo

var wg sync.WaitGroup

type PixInfo struct{
	num int
    pixType map[int]*PixType
	contentType map[string]bool
}
func NewPixInfo(pixType map[int]*PixType, num int) *PixInfo{
	ct := make(map[string]bool, num)
	for _,v := range pixType{
	     ct[v.ct] = true
	}
	return &PixInfo{pixType:pixType,contentType:ct}
}
// get file ext 
func (p *PixInfo)getExt(id int) string{
    return "." + p.pixType[id].name
}

type PixType struct{
    id   int
	name string
	ct   string
	prefix []byte
}
func NewPixType(id int, name string, prefix []byte) *PixType{
	return &PixType{id, name, "image/"+name , prefix}
}

//count the number of pix
type Counter struct {
	//number of pix created
    c int
	muc sync.Mutex
	//number of response handled
	s int
	mus sync.Mutex	
}
func NewCounter() *Counter{
	return &Counter{c:0, muc:sync.Mutex{}, s:0, mus:sync.Mutex{}}
}


type PixHandler struct{
    buf chan []byte
	workerNum int
	c *Counter
}
func NewPixHandler(workerNum int) (ph *PixHandler){
	ph = &PixHandler{buf:make(chan []byte, 1024), workerNum:workerNum, c:NewCounter(),}
	for i:=0;i<workerNum;i++ {
		go ph.startWorker()
	}
	return 
}
//start a worker for saving pix
func (ph *PixHandler)startWorker(){
	for{
	    b := <-ph.buf
		err, id := ph.filterBytes(b)
		if err != nil {
			wg.Done()
		    continue
		}
		err = ioutil.WriteFile(DIR + "/" + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(100)) + pixSet.getExt(id), b, 0666)
		wg.Done()
		if err != nil {
			log.Println(err)
			continue
		}
		ph.c.muc.Lock()
		ph.c.c++
		ph.c.muc.Unlock()
	}
}
//read bytes from response body 
func (ph *PixHandler)handle (resp *http.Response) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Print("failed to read response body.")
	}
	wg.Add(1)
	ph.buf <- body
}
//filter wrong type data
func (ph *PixHandler)filterBytes(body []byte) (error, int){
    res := -1
	if len(body) == 0 {
		return errors.New("empty response body"), res
	}
	res = JPEG
	for k,v := range pixSet.pixType[JPEG].prefix {
		if body[k]^v != 0 {
		   res = -1
		   break; 
		}		
	}
	if res == -1 {
	    res = PNG
		for k,v := range pixSet.pixType[PNG].prefix {
			if body[k]^v != 0 {
		       res = -1
		       break; 
		    }		
	    }
	}
	if res == -1 {
		return errors.New("wrong data type"), res
	}
	ph.c.mus.Lock()
	ph.c.s++
	ph.c.mus.Unlock()
	return nil,res
}
//print statistics
func (ph *PixHandler)statics() {
	logger.Print("sniffed:" + strconv.Itoa(ph.c.s) +  " created:" +  strconv.Itoa(ph.c.c)+ " rate:" +strconv.FormatFloat(float64(ph.c.c/ph.c.s*100),'f',2,64)+"%")
}