package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"runtime"
	"time"
)

var mem runtime.MemStats

func PrintMemory() (float64) {
	runtime.ReadMemStats(&mem)
	return float64(mem.TotalAlloc)/float64(1048576)
}

var STOP byte = 0
var EQ byte = 20
var GT byte = 17
var PUSH1 byte = 96
var PUSH2 byte = 97
var PUSH4 byte = 99
var PUSH8 byte = 103
var PUSH20 byte = 115
var PUSH32 byte = 127
var STATICCALL byte = 250
var DELETAGECALL byte = 244
var SSTORE byte = 85
var RETURN byte = 243
var JUMPDEST byte = 91
var CALLDATESIZE byte = 54
var JUMPI byte = 87
var JUMP byte = 86
var SELFDESTRUCT byte = 255

type code struct {
	Index    int
	Bytecode byte
	Func     int
}

func main() {
	all := 0
	solidity := 0
	attack := 0
	time_all := 0.0
	mem_all := 0.0
	time_max := 0.0
	mem_max := 0.0
	for i := 0; i < 10; i++ {
		dn := strconv.Itoa(i)
		dir := "../../data/contracts/" + dn
		files, _ := ioutil.ReadDir(dir)
		fmt.Println(len(files))
        all += len(files)
        attack_dir := 0
		for _, file := range files {
			data, _ := ioutil.ReadFile(dir + "/" + file.Name())
			startMem := PrintMemory()
			start := time.Now()
			if CheckSol(string(data)) {
				solidity += 1
				split := returnsplit(string(data), file.Name())
				contain, funs := split_contains(split, string(data), file.Name())
				if contain && split_same(string(data), funs, split) {
					fmt.Println(file.Name())
                    attack += 1
                    attack_dir += 1
					f, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY, 0600)
					defer f.Close()
					fmt.Fprintln(f, file.Name())
					stopMem_attack := PrintMemory()
					stop_attack := time.Now()
					mem_max += stopMem_attack-startMem
					time_max += stop_attack.Sub(start).Seconds()
				}
			}
			stopMem := PrintMemory()
			stop := time.Now()
			mem_all += stopMem-startMem
			time_all += stop.Sub(start).Seconds()
		}
		f, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY, 0600)
		defer f.Close()
        fmt.Fprintln(f, "\n")
        fmt.Print("attack: ")
        fmt.Println(attack_dir)
	}
	fmt.Println(all)
	fmt.Println(solidity)
	fmt.Println(attack)
	fmt.Println(time_max/float64(attack))
	fmt.Println(mem_max/float64(attack))
	fmt.Println(time_all/float64(all))
	fmt.Println(mem_all/float64(all))
}

func split_same(bytes string, funs [][]byte, split []code) bool {
	hex1, _ := hex.DecodeString(bytes[2:])
	count := make([]int, len(funs))
	count1 := make([]int, len(funs))
	ignore := -1
	for i, v := range hex1 {
		if v == 99 && i < len(hex1)-5 { //push4
			for j, w := range funs {
				if reflect.DeepEqual(hex1[i+1:i+5], w) && hex1[i+5] != 20 && hex1[i+6] != 20 { //eq
					count[j] = count[j] + 1
				}
			}
		}
		Ignore(v, i, &ignore)
	}
	ignore = -1
	for i, v := range split {
		if v.Bytecode == 99 && v.Func == 1 && i < len(split)-4 {
			for j, w := range funs {
				if reflect.DeepEqual(hex1[i+1:i+5], w) {
					count1[j] = count1[j] + 1
				}
			}
		}
		Ignore(v.Bytecode, i, &ignore)
	}
	for i, _ := range funs {
		if count[i]-count1[i] >= 1 {
			return true
		}
	}
	return false
}
func split_contains(s []code, bytes string, file string) (bool, [][]byte) {
	count := 0
	ignore := -1
	fun := []byte{}
	funs := [][]byte{}
    //transfer := [][]byte{{169, 5, 156, 187}, {255, 255, 255, 255}, {210, 90, 153, 44}, {111, 227, 202, 186}, {224, 166, 159, 13}, {167, 89, 195, 162}}
    //transfer := [][]byte{{255, 255, 255, 255}, {210, 90, 153, 44}, {111, 227, 202, 186}, {224, 166, 159, 13}, {167, 89, 195, 162}}
    transfer := [][]byte{{169, 5, 156, 187}, {255, 255, 255, 255}}
    //0xa9059cbb: trasfer
    //0xd25a992c ? *
    //0x6fe3caba ? *
    //0xe0a69f0d 0x
	//transfer = append(transfer, hex1[i+1:i+5])
	//transfer := [][]byte{}
	hex1, _ := hex.DecodeString(bytes[2:])
	for i, v := range s {
		if i > ignore && i < len(s)-2 && v.Func == 1 {
			if v.Bytecode == 99 && i < len(s)-5 { //push4 20も
				kaburi := false
				for _, w := range transfer {
					if reflect.DeepEqual(hex1[i+1:i+5], w) {
						kaburi = true
					}
				}
				if !kaburi {
					fun = hex1[i+1 : i+5]
				}
			}
			if v.Bytecode == 241 || v.Bytecode == 242 || v.Bytecode == 250 {
				if len(fun) > 0 {
					count = count + 1
					funs = append(funs, fun)
					transfer = append(transfer, fun)
				}
				fun = []byte{}
			}
		}
		Ignore(v.Bytecode, i, &ignore)
	}
	if count > 0 {
		return true, funs
	}
	return false, funs
}


func CheckSol(tx string) bool {
	if len(tx) <= 10 {
		return false
	} else if tx[2:4] == "60" {
		return true
	}
	return false
}

func Ignore(v byte, i int, ignore *int) {
	if v >= 96 && v <= 127 && i > *ignore {
		*ignore = i + int(v) - 95
	}
}

func returnsplit(bytes string, file string) []code {
	hex1, _ := hex.DecodeString(bytes[2:])
	codes := []code{}
	for i, v := range hex1 {
		var code1 code
		code1.Index = i
		code1.Bytecode = v
		code1.Func = 0
		codes = append(codes, code1)
	}
	afterSub := false
	sub0 := 0
	tags := []int{0}
	tagi := 0
	first := true
	ignore := -1
	gtignore := -1
	aftercalldatasize := false
	for i, v := range hex1[:len(hex1)-3] {
		if i > ignore && i > gtignore {
			if !afterSub {
				if (v == 0 && hex1[i+1] == 96) || (v == 243 && hex1[i+1] == 96) { //push1 ??
					afterSub = true
					sub0 = i + 1
				}
			} else {
				if v == 255 && hex1[i+1] == 96 { //suicide
					sub0 = i + 1
				}
				if first && v == 17 { //gt
					gtignore = sub0 + int(hex1[i+2])*256 + int(hex1[i+3])
				}
				if v == 54 { //calldatasize
					aftercalldatasize = true
				}
				if first && aftercalldatasize && v == 87 && hex1[i-3] == 97 { //jumpi
					next := sub0 + int(hex1[i-2])*256 + int(hex1[i-1])
					if next < len(hex1)-1 { //ここでかかるファイル注意
						tags = append(tags, next)
						tagi = tagi + 1
						first = false
					} else {
						break
					}
				}
				if first && v == 91 {
					break
				}
				if !first && i > tags[tagi] { //jumpi
					if v == 87 {
						next := sub0 + int(hex1[i-2])*256 + int(hex1[i-1])
						if next > tags[1] && next < len(hex1)-1 {
							tags = append(tags, next)
							tagi = tagi + 1
						}
					}
					if v == 91 { //jumpdest
						break
					}
				}
			}
		}
		Ignore(v, i, &ignore)
	}
	if len(tags) > 1 {
		for _, v := range tags[1:] {
			j := v + 1
			for hex1[j] != JUMPDEST {
				codes[j].Func = 1
				if hex1[j] >= 96 && hex1[j] <= 127 {
					for i := 0; i < int(hex1[j])-95 && j+i+1 < len(hex1); i++ {
						codes[j+i+1].Func = 1
					}
					if j+int(hex1[j])-94 >= len(hex1) {
						j = len(hex1) - 1
					} else {
						j = j + int(hex1[j]) - 94
					}
				} else {
					j = j + 1
				}
				if j == len(hex1) {
					break
				}
			}
		}
	}
	return codes
}
