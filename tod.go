package main

import (
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "reflect"
    "runtime"
    "time"
)

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
var SLOAD byte = 84
var RETURN byte = 243
var JUMPDEST byte = 91
var CALLDATESIZE byte = 54
var JUMPI byte = 87
var JUMP byte = 86
var SELFDESTRUCT byte = 255

var contract_address = "0x0000000000000000000000000000000000000000"
var func_address, _ = hex.DecodeString("0xd96a094a"[2:])//buy
var func_address2, _ = hex.DecodeString("0x8d6cc56d"[2:])//updatePrice

var mem runtime.MemStats

func PrintMemory() (float64) {
    runtime.ReadMemStats(&mem)
    return float64(mem.TotalAlloc)/float64(1048576)
}

func Ignore(v byte, i int, ignore *int) {
    if v >= PUSH1 && v <= PUSH32 && i > *ignore {
        *ignore = i + int(v) - 95
    }
}

func main(){
    data, _ := ioutil.ReadFile(contract_address)
    hex1, _ := hex.DecodeString(string(data)[2:])
    ignore := -1
    afterSub := false
    sub0 := 0
    tags := []int{0}
    tagi := 0
    next := 0
    start := time.Now()
    startMem := PrintMemory()
    for i, v := range hex1 {
        if i > ignore {
            if !afterSub {
                if (v == STOP &&  hex1[i+1] == PUSH1) || (v == RETURN && hex1[i+1] == PUSH1) {
                    afterSub = true
                    sub0 = i + 1
                } else if (v == STOP && hex1[i+2] == PUSH1) || (v == RETURN && hex1[i+2] == PUSH1) {
                    afterSub = true
                    sub0 = i + 2
                }
            } else if v == EQ && reflect.DeepEqual(hex1[i-4:i], func_address) && next == 0 {
                next = sub0 + int(hex1[i+2])*256+int(hex1[i+3])
                tags = append(tags, next)
                tagi = tagi + 1
            } else if next != 0 && i > tags[tagi] + 1 {
                if v == JUMPI {
                    next = sub0 + int(hex1[i-2])*256 + int(hex1[i-1])
                    if next > tags[1] && next < len(hex1)-1 {
                        tags = append(tags, next)
                        tagi = tagi + 1
                    }
                } else if v == JUMP && hex1[i-3] == PUSH2 {
                    next = sub0 + int(hex1[i-2])*256 + int(hex1[i-1])
                    if next > tags[1] && next < len(hex1)-1 {
                        tags = append(tags, next)
                        tagi = tagi + 1
                    }
                } else if v == JUMPDEST {
                    break
                }
            }
        }
        Ignore(v, i, &ignore)
    }
    fmt.Println(tags)
    for _, v := range tags[1:] {
        ignore = -1
        first := true
        firstpush := -1
        for j, w := range hex1[v+1:] {
            if j > ignore {
                if w == SSTORE {
                    fmt.Print("SSTORE ")
                    if hex1[j+v-1] == PUSH1 {
                        fmt.Println(hex1[j+v])
                    } else {
                        fmt.Println(firstpush)
                    }
                } else if w == SLOAD {
                    fmt.Print("SLOAD " )
                    if hex1[j+v-1] == PUSH1 {
                        fmt.Println(hex1[j+v])
                    } else {
                        fmt.Println(firstpush)
                    }
                } else if w == PUSH1 && first {
                    first = false
                    firstpush = int(hex1[j+v+2])
                } else if w == JUMPDEST {
                    break
                }
            }
            Ignore(w, j, &ignore)
        }
    }
    stopMem := PrintMemory()
    stop := time.Now()
    fmt.Println(stopMem-startMem)
    fmt.Println(stop.Sub(start).Seconds())
}
