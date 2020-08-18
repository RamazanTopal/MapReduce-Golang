package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

func dosyakaydet(url string, sourcekac string) {
	// Make request
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create output file
	outFile, err := os.Create(sourcekac)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Copy data from HTTP response to file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//dosya kaydetme kısmı
	dosyakaydet("http://www.elzempen.com/cambalkon.html", "source1.txt")
	dosyakaydet("https://www.tarihiolaylar.com/tarihi-olaylar/yazi-161", "source2.txt")
	dosyakaydet("http://moradergisi.com/haber/yazinin-icadi-ve-tarihcesi-69406.html", "source3.txt")
	dosyakaydet("https://onlineoyuntarihi.wordpress.com/", "source4.txt")
	dosyakaydet("https://www.muhasebedersleri.com/ekonomi/para.html", "source5.txt")
	//1.bölüm
	uzunluk := 18
	text1_1 := make(chan string, uzunluk)
	text2_2 := make(chan string, uzunluk)
	text3_3 := make(chan string, uzunluk)
	text4_4 := make(chan string, uzunluk)
	text5_5 := make(chan string, uzunluk)
	map1_1 := make(chan map[string]int, uzunluk)
	map2_2 := make(chan map[string]int, uzunluk)
	map3_3 := make(chan map[string]int, uzunluk)
	map4_4 := make(chan map[string]int, uzunluk)
	map5_5 := make(chan map[string]int, uzunluk)
	reduce1_1 := make(chan int, uzunluk)
	reduce2_2 := make(chan int, uzunluk)
	reduce3_3 := make(chan int, uzunluk)
	reduce4_4 := make(chan int, uzunluk)
	toplam1_1 := make(chan float32, uzunluk)
	toplam2_2 := make(chan float32, uzunluk)
	toplam3_3 := make(chan float32, uzunluk)
	toplam4_4 := make(chan float32, uzunluk)
	go inputReader1([5]chan<- string{text1_1, text2_2, text3_3, text4_4, text5_5})
	go mapper1(text1_1, map1_1)
	go mapper1(text2_2, map2_2)
	go mapper1(text3_3, map3_3)
	go mapper1(text4_4, map4_4)
	go mapper1(text5_5, map5_5)
	go shuffler1([]<-chan map[string]int{map1_1, map2_2, map3_3, map4_4, map5_5}, [4]chan<- int{reduce1_1, reduce2_2, reduce3_3, reduce4_4})
	go reducer1(reduce1_1, toplam1_1)
	go reducer1(reduce2_2, toplam2_2)
	go reducer1(reduce3_3, toplam3_3)
	go reducer1(reduce4_4, toplam4_4)
	outputWriter1([]<-chan float32{toplam1_1, toplam2_2, toplam3_3, toplam4_4})

	//2.bölüm
	size := 18
	text1 := make(chan string, size)
	text2 := make(chan string, size)
	text3 := make(chan string, size)
	text4 := make(chan string, size)
	text5 := make(chan string, size)
	map1 := make(chan map[string]int, size)
	map2 := make(chan map[string]int, size)
	map3 := make(chan map[string]int, size)
	map4 := make(chan map[string]int, size)
	map5 := make(chan map[string]int, size)
	reduce1 := make(chan int, size)
	reduce2 := make(chan int, size)
	reduce3 := make(chan int, size)
	reduce4 := make(chan int, size)
	toplam1 := make(chan float32, size)
	toplam2 := make(chan float32, size)
	toplam3 := make(chan float32, size)
	toplam4 := make(chan float32, size)
	go inputReader([5]chan<- string{text1, text2, text3, text4, text5})
	go mapper(text1, map1)
	go mapper(text2, map2)
	go mapper(text3, map3)
	go mapper(text4, map4)
	go mapper(text5, map5)
	go shuffler([]<-chan map[string]int{map1, map2, map3, map4, map5}, [4]chan<- int{reduce1, reduce2, reduce3, reduce4})
	go reducer(reduce1, toplam1)
	go reducer(reduce2, toplam2)
	go reducer(reduce3, toplam3)
	go reducer(reduce4, toplam4)
	outputWriter([]<-chan float32{toplam1, toplam2, toplam3, toplam4})
}
func inputReader(out [5]chan<- string) {
	data1, err := ioutil.ReadFile("source1.txt")
	if err != nil {
		fmt.Println(err)
	}
	data2, err := ioutil.ReadFile("source2.txt")
	if err != nil {
		fmt.Println(err)
	}
	data3, err := ioutil.ReadFile("source3.txt")
	if err != nil {
		fmt.Println(err)
	}
	data4, err := ioutil.ReadFile("source4.txt")
	if err != nil {
		fmt.Println(err)
	}
	data5, err := ioutil.ReadFile("source5.txt")
	if err != nil {
		fmt.Println(err)
	}

	s1 := strings.Fields(string(data1))
	s2 := strings.Fields(string(data2))
	s3 := strings.Fields(string(data3))
	s4 := strings.Fields(string(data4))
	s5 := strings.Fields(string(data5))
	input := [][]string{
		s1,
		s2,
		s3,
		s4,
		s5,
	}
	for i := range out {
		go func(ch chan<- string, word []string) {
			for _, w := range word {
				ch <- w
			}
			close(ch)
		}(out[i], input[i])
	}
}

func mapper(in <-chan string, out chan<- map[string]int) {
	count := map[string]int{}
	for word := range in {
		count[word] = count[word] + 1
	}
	out <- count
	close(out)
}
func shuffler(in []<-chan map[string]int, out [4]chan<- int) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func(c <-chan map[string]int) {
			for m := range c {
				nc, ok := m["</p>"]
				if ok {
					out[0] <- nc
				}
				vc, ok := m["</h1>"]
				if ok {
					out[1] <- vc
				}
				mc, ok := m["</li>"]
				if ok {
					out[2] <- mc
				}
				kc, ok := m["</ul>"]
				if ok {
					out[3] <- kc
				}
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out[0])
		close(out[1])
		close(out[2])
		close(out[3])
	}()
}
func reducer(in <-chan int, out chan<- float32) {
	sum := 0
	for n := range in {
		sum += n

	}
	out <- float32(sum)
	close(out)
}
func outputWriter1(in []<-chan float32) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	name := []string{"</p>", "</h1>", "</li>", "</ul>"}
	for i := 0; i < len(in); i++ {
		go func(n int, c <-chan float32) {
			for toplam := range c {
				fmt.Printf(" %s etiketinin toplamı : %f\n", name[n], toplam)
			}
			wg.Done()
		}(i, in[i])
	}
	wg.Wait()
}

//2.Bölüm
func inputReader1(out [5]chan<- string) {
	data1, err := ioutil.ReadFile("source1.txt")
	if err != nil {
		fmt.Println(err)
	}
	data2, err := ioutil.ReadFile("source2.txt")
	if err != nil {
		fmt.Println(err)
	}
	data3, err := ioutil.ReadFile("source3.txt")
	if err != nil {
		fmt.Println(err)
	}
	data4, err := ioutil.ReadFile("source4.txt")
	if err != nil {
		fmt.Println(err)
	}
	data5, err := ioutil.ReadFile("source5.txt")
	if err != nil {
		fmt.Println(err)
	}

	s1 := strings.Fields(string(data1))
	s2 := strings.Fields(string(data2))
	s3 := strings.Fields(string(data3))
	s4 := strings.Fields(string(data4))
	s5 := strings.Fields(string(data5))
	input := [][]string{
		s1,
		s2,
		s3,
		s4,
		s5,
	}
	for i := range out {
		go func(ch chan<- string, word []string) {
			for _, w := range word {
				ch <- w
			}
			close(ch)
		}(out[i], input[i])
	}
}
func mapper1(in <-chan string, out chan<- map[string]int) {
	count := map[string]int{}
	for word := range in {
		count[word] = count[word] + 1
	}
	out <- count
	close(out)
}
func shuffler1(in []<-chan map[string]int, out [4]chan<- int) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func(c <-chan map[string]int) {
			for m := range c {
				nc, ok := m["</p>"]
				if ok {
					out[0] <- nc
				}
				vc, ok := m["</h1>"]
				if ok {
					out[1] <- vc
				}
				mc, ok := m["</li>"]
				if ok {
					out[2] <- mc
				}
				kc, ok := m["</ul>"]
				if ok {
					out[3] <- kc
				}
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out[0])
		close(out[1])
		close(out[2])
		close(out[3])
	}()
}

func reducer1(in <-chan int, out chan<- float32) {
	sum := 0
	for n := range in {
		sum += n

	}
	out <- float32(sum)
	close(out)
}
func outputWriter(in []<-chan float32) {

	a := make(chan float32)
	go sumcountp(a, "http://www.elzempen.com/cambalkon.html")
	cikti1 := <-a
	//2.p
	b := make(chan float32)
	go sumcountp(b, "https://www.tarihiolaylar.com/tarihi-olaylar/yazi-161")
	cikti2 := <-b
	//3.p
	c := make(chan float32)
	go sumcountp(c, "http://moradergisi.com/haber/yazinin-icadi-ve-tarihcesi-69406.html")
	cikti3 := <-c
	//4.p
	d := make(chan float32)
	go sumcountp(d, "https://onlineoyuntarihi.wordpress.com/")
	cikti4 := <-d
	//5.p
	e := make(chan float32)
	go sumcountp(e, "https://www.muhasebedersleri.com/ekonomi/para.html")
	cikti5 := <-e
	ciktitoplamp := cikti1 + cikti2 + cikti3 + cikti4 + cikti5
	//1.h
	f := make(chan float32)
	go sumcounth(f, "http://www.elzempen.com/cambalkon.html")
	cikti6 := <-f
	//2.h
	g := make(chan float32)
	go sumcounth(g, "https://www.tarihiolaylar.com/tarihi-olaylar/yazi-161")
	cikti7 := <-g
	//3.h
	h := make(chan float32)
	go sumcounth(h, "http://moradergisi.com/haber/yazinin-icadi-ve-tarihcesi-69406.html")
	cikti8 := <-h
	//4.h
	j := make(chan float32)
	go sumcounth(j, "https://onlineoyuntarihi.wordpress.com/")
	cikti9 := <-j
	//5.h
	k := make(chan float32)
	go sumcounth(k, "https://www.muhasebedersleri.com/ekonomi/para.html")
	cikti10 := <-k
	ciktitoplamh := cikti6 + cikti7 + cikti8 + cikti9 + cikti10

	//--------------------------------------------------------------------------
	//1.li

	l := make(chan float32)
	go sumcountli(l, "http://www.elzempen.com/cambalkon.html")
	cikti11 := <-l
	//2.li
	m := make(chan float32)
	go sumcountli(m, "https://www.tarihiolaylar.com/tarihi-olaylar/yazi-161")
	cikti12 := <-m
	//3.li
	n := make(chan float32)
	go sumcountli(n, "http://moradergisi.com/haber/yazinin-icadi-ve-tarihcesi-69406.html")
	cikti13 := <-n
	//4.li
	o := make(chan float32)
	go sumcountli(o, "https://onlineoyuntarihi.wordpress.com/")
	cikti14 := <-o
	//5.li
	p := make(chan float32)
	go sumcountli(p, "https://www.muhasebedersleri.com/ekonomi/para.html")
	cikti15 := <-p
	ciktitoplamli := cikti11 + cikti12 + cikti13 + cikti14 + cikti15

	//---------------------------------------------------------------------------------

	//1.strong
	r := make(chan float32)
	go sumcountul(r, "http://www.elzempen.com/cambalkon.html")
	cikti16 := <-r
	//2.strong
	s := make(chan float32)
	go sumcountul(s, "https://www.tarihiolaylar.com/tarihi-olaylar/yazi-161")
	cikti17 := <-s
	//3.strong
	t := make(chan float32)
	go sumcountul(t, "http://moradergisi.com/haber/yazinin-icadi-ve-tarihcesi-69406.html")
	cikti18 := <-t
	//4.strong
	u := make(chan float32)
	go sumcountul(u, "https://onlineoyuntarihi.wordpress.com/")
	cikti19 := <-u
	//5.strong
	v := make(chan float32)
	go sumcountul(v, "https://www.muhasebedersleri.com/ekonomi/para.html")
	cikti20 := <-v
	ciktitoplamul := cikti16 + cikti17 + cikti18 + cikti19 + cikti20

	var wg sync.WaitGroup
	wg.Add(len(in))
	name := []string{"</p>", "</h1>", "</li>", "</ul>"}
	for i := 0; i < len(in); i++ {
		go func(n int, c <-chan float32) {
			for toplam := range c {
				if name[n] == "</p>" {
					fmt.Printf(" %s etiketinin icindeki ortalama harf uzunluğu : %f\n", name[n], ciktitoplamp/toplam)
				}
				if name[n] == "</h1>" {
					fmt.Printf(" %s etiketinin  icindeki ortalama harf uzunluğu : %f\n", name[n], ciktitoplamh/toplam)
				}
				if name[n] == "</li>" {
					fmt.Printf(" %s etiketinin icindeki ortalama harf uzunluğu : %f\n", name[n], ciktitoplamli/toplam)
				}
				if name[n] == "</ul>" {
					fmt.Printf("%s etiketinin icindeki ortalama harf uzunluğu : %f\n", name[n], ciktitoplamul/toplam)
				}
			}
			wg.Done()
		}(i, in[i])
	}
	wg.Wait()
}

func sumcountp(c chan float32, url string) {
	sum := 0.0
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading HTTP body. ", err)
	}

	// Create a regular expression to find comments
	re := regexp.MustCompile("<p>(.|\n)*?</p>")
	comments := re.FindAllString(string(body), -1)

	if comments == nil {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	} else {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	}

}
func sumcounth(c chan float32, url string) {
	sum := 0.0
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading HTTP body. ", err)
	}

	// Create a regular expression to find comments
	re := regexp.MustCompile("<h1>(.|\n)*?</h1>")
	comments := re.FindAllString(string(body), -1)

	if comments == nil {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	} else {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	}

}
func sumcountli(c chan float32, url string) {
	sum := 0.0
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading HTTP body. ", err)
	}

	// Create a regular expression to find comments
	re := regexp.MustCompile("<li>(.|\n)*?</li>")
	comments := re.FindAllString(string(body), -1)

	if comments == nil {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	} else {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	}

}
func sumcountul(c chan float32, url string) {
	sum := 0.0
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading HTTP body. ", err)
	}

	// Create a regular expression to find comments
	re := regexp.MustCompile("<ul>(.|\n)*?</ul>")
	comments := re.FindAllString(string(body), -1)

	if comments == nil {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	} else {
		for _, comment := range comments {
			sum += float64(len(comment))
		}
		c <- float32(sum)
	}

}
