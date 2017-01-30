package main

import (
    "fmt"
)

type TCPFlag int

const (
    TCP_FLAG_FIN    = 0x01
    TCP_FLAG_SYN    = 0x02
    TCP_FLAG_RST    = 0x04
    TCP_FLAG_PUSH   = 0x08
    TCP_FLAG_ACK    = 0x10
    TCP_FLAG_URGENT = 0x20
    TCP_FLAG_CWR    = 0x40
    TCP_FLAG_ECE    = 0x80
)

var TCPFlagNameMap = map[TCPFlag]string{
    TCP_FLAG_FIN:    "F",
    TCP_FLAG_SYN:    "S",
    TCP_FLAG_RST:    "R",
    TCP_FLAG_PUSH:   "P",
    TCP_FLAG_ACK:    "A",
    TCP_FLAG_URGENT: "U",
    TCP_FLAG_CWR:    "C",
    TCP_FLAG_ECE :   "E",
}

var TCPFlagValueMap = map[string]TCPFlag{
    TCPFlagNameMap[TCP_FLAG_FIN]:    TCP_FLAG_FIN,
    TCPFlagNameMap[TCP_FLAG_SYN]:    TCP_FLAG_SYN,
    TCPFlagNameMap[TCP_FLAG_RST]:    TCP_FLAG_RST,
    TCPFlagNameMap[TCP_FLAG_PUSH]:   TCP_FLAG_PUSH,
    TCPFlagNameMap[TCP_FLAG_ACK]:    TCP_FLAG_ACK,
    TCPFlagNameMap[TCP_FLAG_URGENT]: TCP_FLAG_URGENT,
    TCPFlagNameMap[TCP_FLAG_CWR]:    TCP_FLAG_CWR,
    TCPFlagNameMap[TCP_FLAG_ECE]:    TCP_FLAG_ECE,
}

type TCPFlagOp int

const (
    TCP_FLAG_OP_OR      = 0x00
    TCP_FLAG_OP_AND     = 0x40
    TCP_FLAG_OP_END     = 0x80
    TCP_FLAG_OP_NOT     = 0x02
    TCP_FLAG_OP_MATCH   = 0x01
)

var TCPFlagOpNameMap = map[TCPFlagOp]string{
    TCP_FLAG_OP_OR:     " ",
    TCP_FLAG_OP_AND:    "&",
    TCP_FLAG_OP_END:    "E",
    TCP_FLAG_OP_NOT:    "!",
    TCP_FLAG_OP_MATCH:  "=",
}

var TCPFlagOpValueMap = map[string]TCPFlagOp{
    TCPFlagOpNameMap[TCP_FLAG_OP_OR]:       TCP_FLAG_OP_OR,
    TCPFlagOpNameMap[TCP_FLAG_OP_AND]:      TCP_FLAG_OP_AND,
    TCPFlagOpNameMap[TCP_FLAG_OP_END]:      TCP_FLAG_OP_END,
    TCPFlagOpNameMap[TCP_FLAG_OP_NOT]:      TCP_FLAG_OP_NOT,
    TCPFlagOpNameMap[TCP_FLAG_OP_MATCH]:    TCP_FLAG_OP_MATCH,
}

var examples= []string {
    "=SA&!U",
    "=SA =A",
    "=",
    "CEUAPRSF",
    "=SA",
    "SA",
    "!SA",
    "=SA&=!U",
    "!SA&!U",
    "=!SA&=!U",
    "SETUYA",
    "SA&=!U",
}


func main() {
    for _, myCmd := range(examples) {
        err, tcpFlagsValues, tcpFlowSpecOps := parseTcpFlagCmd(myCmd)
        fmt.Printf("command: %s => ", myCmd)
        if err != nil {
            fmt.Println(err)
        } else {
            for index :=0; index < len(tcpFlagsValues); index++ {
                fmt.Printf(" %d flags/ops: %b/%b", index, tcpFlagsValues[index], tcpFlowSpecOps[index])
            }
        fmt.Printf("\n")
        }
    }
}

func parseTcpFlagCmd(myCmd string) (error, []int, []int) {
    var tcpFlagBitMap = map[string]bool{
        TCPFlagNameMap[TCP_FLAG_FIN]:    false,
        TCPFlagNameMap[TCP_FLAG_SYN]:    false,
        TCPFlagNameMap[TCP_FLAG_RST]:    false,
        TCPFlagNameMap[TCP_FLAG_PUSH]:   false,
        TCPFlagNameMap[TCP_FLAG_ACK]:    false,
        TCPFlagNameMap[TCP_FLAG_URGENT]: false,
        TCPFlagNameMap[TCP_FLAG_CWR]:    false,
        TCPFlagNameMap[TCP_FLAG_ECE]:    false,
    }
    var tcpFlagOpBitMap = map[string]bool{
        TCPFlagOpNameMap[TCP_FLAG_OP_AND]:   false,
        TCPFlagOpNameMap[TCP_FLAG_OP_END]:   false,
        TCPFlagOpNameMap[TCP_FLAG_OP_NOT]:   false,
        TCPFlagOpNameMap[TCP_FLAG_OP_MATCH]: false,
    }
    var index int = 0
    var tcpFlagsValues []int
    var tcpFlowSpecOps []int
    var tcpFlagsValue int = 0
    var tcpFlowSpecOp int = 0
    for index < len(myCmd) {
        myCmdChar := myCmd[index:index+1]
        switch myCmdChar {
        case TCPFlagOpNameMap[TCP_FLAG_OP_MATCH]:
            if(tcpFlagOpBitMap[myCmdChar] == false) {
                tcpFlagOpBitMap[myCmdChar] = true
                index++
            } else {
                err := fmt.Errorf("Match flag appears multiple time")
                return err, tcpFlagsValues, tcpFlowSpecOps
            }
        case TCPFlagOpNameMap[TCP_FLAG_OP_NOT]:
            if(tcpFlagOpBitMap[myCmdChar] == false) {
                tcpFlagOpBitMap[myCmdChar] = true
                index++
            } else {
                err := fmt.Errorf("Not flag appears multiple time")
                return err, tcpFlagsValues, tcpFlowSpecOps
            }
        case TCPFlagOpNameMap[TCP_FLAG_OP_AND], TCPFlagOpNameMap[TCP_FLAG_OP_OR]:
            if(tcpFlagOpBitMap[myCmdChar] == false) {
                tcpFlagOpBitMap[myCmdChar] = true
                tcpFlagsValue, tcpFlowSpecOp = setTcpOpsBitmapWithMap(tcpFlagBitMap, tcpFlagOpBitMap)
                tcpFlagsValues = append(tcpFlagsValues, tcpFlagsValue)
                tcpFlowSpecOps = append(tcpFlowSpecOps, tcpFlowSpecOp)
                resetAllFlagsToFalse(tcpFlagBitMap, tcpFlagOpBitMap)
                index++
            } else {
                err := fmt.Errorf("AND or OR (space) operator appears multiple time")
                return err, tcpFlagsValues, tcpFlowSpecOps
            }
        case TCPFlagNameMap[TCP_FLAG_ACK], TCPFlagNameMap[TCP_FLAG_SYN], TCPFlagNameMap[TCP_FLAG_FIN],
        TCPFlagNameMap[TCP_FLAG_URGENT], TCPFlagNameMap[TCP_FLAG_ECE], TCPFlagNameMap[TCP_FLAG_RST],
        TCPFlagNameMap[TCP_FLAG_CWR], TCPFlagNameMap[TCP_FLAG_PUSH]:
            myLoopChar := myCmdChar
            loopIndex := index
            // we loop till we reach the end of TCP flags description
            // exit conditions : we reach the end of tcp flags bacause we found & or ' ' or we reach the end of the line
            for (loopIndex < len(myCmd) &&
                (myLoopChar != TCPFlagOpNameMap[TCP_FLAG_OP_AND] && myLoopChar != TCPFlagOpNameMap[TCP_FLAG_OP_OR])) {
                // we check that charater is a well known flag and is not doubled
                if (TCPFlagValueMap[myLoopChar]!= 0 && tcpFlagBitMap[myLoopChar] == false) {
                    tcpFlagBitMap[myLoopChar] = true            // we set this flag to true
                    loopIndex++                                 // we move to next character
                    if(loopIndex < len(myCmd)) {
                        myLoopChar = myCmd[loopIndex:loopIndex+1]   // we set next character only if we didn't reach the end of cmd
                    }
                } else {
                    err := fmt.Errorf("flag %s appears multiple time or is not part of TCP flags", myLoopChar)
                    return err, tcpFlagsValues, tcpFlowSpecOps
                }
            }
            // we are done with flags, we give back where we are for the first loop
            index = loopIndex
        default:
            err := fmt.Errorf("flag %s not part of tcp flags", myCmdChar)
            return err, tcpFlagsValues, tcpFlowSpecOps
        }
    }
    tcpFlagOpBitMap["E"] = true
    tcpFlagsValue, tcpFlowSpecOp = setTcpOpsBitmapWithMap(tcpFlagBitMap, tcpFlagOpBitMap)
    tcpFlagsValues = append(tcpFlagsValues, tcpFlagsValue)
    tcpFlowSpecOps = append(tcpFlowSpecOps, tcpFlowSpecOp)
    resetAllFlagsToFalse(tcpFlagBitMap, tcpFlagOpBitMap)
    return nil, tcpFlagsValues, tcpFlowSpecOps
}

func resetAllFlagsToFalse(myTcpFlagBitMap map[string]bool, myTcpFlagOpBitMap map[string]bool) {
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_FIN]]     = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_SYN]]     = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_RST]]     = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_PUSH]]    = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_ACK]]     = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_URGENT]]  = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_CWR]]     = false
    myTcpFlagBitMap[TCPFlagNameMap[TCP_FLAG_ECE]]     = false
    myTcpFlagOpBitMap[TCPFlagOpNameMap[TCP_FLAG_OP_AND]]   = false
    myTcpFlagOpBitMap[TCPFlagOpNameMap[TCP_FLAG_OP_OR]]    = false
    myTcpFlagOpBitMap[TCPFlagOpNameMap[TCP_FLAG_OP_END]]   = false
    myTcpFlagOpBitMap[TCPFlagOpNameMap[TCP_FLAG_OP_NOT]]   = false
    myTcpFlagOpBitMap[TCPFlagOpNameMap[TCP_FLAG_OP_MATCH]] = false
}

func setTcpOpsBitmapWithMap(myTcpFlagBitMap map[string]bool, myTcpFlagOpBitMap map[string]bool) (int, int) {
    var myFlags int = 0
    var myOps int = 0
    for flagString, flagSetUnset := range myTcpFlagBitMap {
        if (flagSetUnset) {
            myFlags |= int(TCPFlagValueMap[flagString])
        }
    }
    for opString, opSetUnset := range myTcpFlagOpBitMap {
        if (opSetUnset) {
            myOps |= int(TCPFlagOpValueMap[opString])
        }
    }
    return myFlags, myOps
}
