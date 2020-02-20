package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func calcCRC32(iter int, data string, result []string,
	wg *sync.WaitGroup) {
	res := DataSignerCrc32(data)
	//mu.Lock()
	result[iter] = res
	//mu.Unlock()

	wg.Done()
}

func SingleHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	for dataRaw := range in {
		dataInt, ok := dataRaw.(int)
		if !ok {
			fmt.Println("cant convert result data to string")
		}
		data := strconv.Itoa(dataInt)
		//r1 := DataSignerCrc32(data)
		//tmp := DataSignerMd5(data)
		//r2 := DataSignerCrc32(tmp)

		//results := make([]string, 2)

		wg.Add(1)

		go func(data string) {
			dataHash1 := make(chan string)
			dataHash2 := make(chan string)
			go func() {
				dataHash1 <- DataSignerCrc32(data)
			}()

			go func() {
				mu.Lock()
				md5hash := DataSignerMd5(data)
				mu.Unlock()
				dataHash2 <- DataSignerCrc32(md5hash)

			}()
			out <- <-dataHash1 + "~" + <-dataHash2
			wg.Done()
		}(data)

		// fmt.Println(data, "SingleHash data", data)
		// fmt.Println(data, "SingleHash md5(data)", tmp)
		// fmt.Println(data, "SingleHash crc32(md5(data))", results[1])
		// fmt.Println(data, "SingleHash crc32(data)", results[0])
		// fmt.Println(data, "SingleHash result", res)
		//out <- res
	}
	wg.Wait()
}

/*
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			dataHash1 := make(chan string)
			dataHash2 := make(chan string)

			go func() {
				dataHash1 <- DataSignerCrc32(data)
			}()

			go func() {
				mu.Lock()
				md5hash := DataSignerMd5(data)
				mu.Unlock()
				dataHash2 <- DataSignerCrc32(md5hash)
			}()

			out <- <-dataHash1 + "~" + <-dataHash2
		}(value)
	}

	wg.Wait()
}
*/
func calcCRC32Iter(iter int, data string, result []string,
	wg *sync.WaitGroup) {
	res := DataSignerCrc32(strconv.Itoa(iter) + data)
	result[iter] = res

	wg.Done()
}

func MultiHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)
	for dataRaw := range in {

		data, ok := dataRaw.(string)
		if !ok {
			fmt.Println("cant convert result data to string")
		}

		wg.Add(1)

		go func(data string) {
			defer wg.Done()
			workerWG := new(sync.WaitGroup)

			results := make([]string, 6)
			for i := 0; i < 6; i++ {
				workerWG.Add(1)
				go func(th int) {
					results[th] = DataSignerCrc32(strconv.Itoa(th) + data)
					workerWG.Done()
				}(i)
			}
			workerWG.Wait()

			res := ""
			for i, result := range results {
				fmt.Println(data, "MultiHash: crc32(th+step1)", i, result)
				res += result
			}
			fmt.Println(data, "MultiHash result", res)
			out <- res
		}(data)

	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	concatResult := ""
	inputs := []string{}
	for dataRaw := range in {
		data, ok := dataRaw.(string)
		if !ok {
			fmt.Println("cant convert result data to string")
		}
		inputs = append(inputs, data)
	}
	sort.Strings(inputs)

	for i, input := range inputs {
		concatResult += input
		if i != len(inputs)-1 {
			concatResult += "_"
		}
	}
	out <- concatResult
}

func handleJob(curJob job, in, out chan interface{}) {
	go func() {
		//func() {
		curJob(in, out)
		close(out)
	}()
}

func ExecutePipeline(jobs ...job) {
	channels := make([]chan interface{}, len(jobs)+1)
	for i := range channels {
		channels[i] = make(chan interface{}, 100)
	}

	for i, curJob := range jobs {
		in := channels[i]
		out := channels[i+1]
		go handleJob(curJob, in, out)
	}
	//this is wait for the last out
	for range channels[len(channels)-1] {

	}
	//this  is also works
	//<-channels[len(channels)-1]
}
