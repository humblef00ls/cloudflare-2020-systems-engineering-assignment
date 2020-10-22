package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"
)

func main() {

	flg := flag.String("url", "", "URL to make a GET request to")
	itr := flag.Int("profile", 1, "How many times to make the request")
	flag.Parse()

	if *flg == "" {
		log.Fatal("URL can not be empty")
	}
	if *itr < 1 {
		log.Fatal("Can not make less than 1 requests")
	}
	s := *flg

	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}

	var result string
	var leasttime int64 = 99999999
	var mosttime int64 = 0
	var leastsize int = 99999999
	var mostsize int = 0
	var timings []int64
	var total int64
	var lines []string
	var errorc = 0
	for i := 0; i < *itr; i++ {
		start := time.Now()

		conn, err := net.Dial("tcp", u.Host+":80")
		if err != nil {
			log.Fatal(err)
		}

		rt := fmt.Sprintf("GET %v HTTP/1.1\r\n", u.Path)
		rt += fmt.Sprintf("Host: %v\r\n", u.Host)
		rt += fmt.Sprintf("Connection: close\r\n")
		rt += fmt.Sprintf("\r\n")
		_, err = conn.Write([]byte(rt))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatal(err)
		}

		sizer := len(resp)
		result = string(resp)

		conn.Close()

		elapsed := time.Since(start).Milliseconds()
		scanner := bufio.NewScanner(strings.NewReader(result))
		x := 0
		for scanner.Scan() {
			if x < 1 {
				srf := strings.Split(scanner.Text(), " ")
				if srf[1] != "200" {
					lines = append(lines, srf[1])
					errorc = errorc + 1
				}
			}
			x = x + 1
		}

		if elapsed > mosttime {
			mosttime = elapsed
		}
		if elapsed < leasttime {
			leasttime = elapsed
		}
		if sizer > mostsize {
			mostsize = sizer
		}
		if sizer < leastsize {
			leastsize = sizer
		}
		timings = append(timings, elapsed)
		total = total + elapsed
	}

	fmt.Println(result)
	fmt.Printf("Made %v Requests\n", *itr)
	fmt.Println("Fastest Time:", leasttime, "ms")
	fmt.Println("Slowest Time:", mosttime, "ms")
	fmt.Println("Mean Time:", total/int64(*itr), "ms")
	sort.Slice(timings, func(i, j int) bool { return timings[i] < timings[j] })
	if *itr%2 == 0 {
		fmt.Println("Median Time:", (timings[(*itr-1)/2]+timings[*itr/2])/int64(2), "ms")
	} else {
		fmt.Println("Median Time:", timings[*itr/2], "ms")
	}
	fmt.Println("Success Rate:", (1-(errorc / *itr))*100, "%")
	fmt.Println("Error Code(s):", lines)
	fmt.Println("Smallest Size:", leastsize, "B")
	fmt.Println("Largest Size:", mostsize, "B")
}
