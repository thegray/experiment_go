package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chilts/sid"
	guuid "github.com/google/uuid"
	"github.com/kjk/betterguid"
	"github.com/lithammer/shortuuid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	suuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
)

func genShortUUID() {
	id := shortuuid.New()
	fmt.Printf("github.com/lithammer/shortuuid: %s\n", id)
}

func genUUID() {
	id := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", id.String())
	fmt.Printf("github.com/google/uuid:         %d\n", id.ID())
}

func genXid() {
	id := xid.New()
	id2 := xid.New()
	fmt.Printf("github.com/rs/xid:              %s\n", id.String())
	fmt.Println("id?: ", id.Counter())

	xx := binary.BigEndian.Uint64(id.Bytes())
	fmt.Println("BE1 uint64: ", xx)
	xx2 := binary.BigEndian.Uint64(id2.Bytes())
	fmt.Println("BE2 uint64: ", xx2)

	x1 := binary.LittleEndian.Uint64(id.Bytes())
	fmt.Println("LE1 uint64: ", x1)
	x2 := binary.LittleEndian.Uint64(id2.Bytes())
	fmt.Println("LE2 uint64: ", x2)

	// yy := binary.BigEndian.Uint32(id.Bytes())
	// fmt.Println("uint32: ", yy)

	fmt.Println("byte: ", id.Bytes())
	// xy := customRep(id.Bytes())
	// fmt.Println("custom: ", xy)
}

func genKsuid() {
	id := ksuid.New()
	fmt.Printf("github.com/segmentio/ksuid:     %s\n", id.String())
}

func genBetterGUID() {
	id := betterguid.New()
	fmt.Printf("github.com/kjk/betterguid:      %s\n", id)
}

func genUlid() {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	fmt.Printf("github.com/oklog/ulid:          %s\n", id.String())
}

func genSonyflake() {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	// Note: this is base16, could shorten by encoding as base62 string
	fmt.Printf("github.com/sony/sonyflake:      %x\n", id)
}

func genSid() {
	id := sid.Id()
	fmt.Printf("github.com/chilts/sid:          %s\n", id)
}

func genUUIDv4() {
	id := suuid.NewV4()
	fmt.Printf("github.com/satori/go.uuid:      %s\n", id)
}

func main() {
	genXid()
	// genKsuid()
	// genBetterGUID()
	// genUlid()
	// genSonyflake()
	// genSid()
	// genShortUUID()
	// genUUIDv4()
	// genUUID()
}

// func customRep(b []byte) big.Int {
// 	_ = b[11]
// 	// return uint64(b[11]) | uint64(b[10])<<1 | uint64(b[9])<<2 |
// 	// 	uint64(b[8])<<3 | uint64(b[7])<<4 | uint64(b[6])<<5 |
// 	// 	uint64(b[5])<<6 | uint64(b[4])<<7 | uint64(b[3])<<8 |
// 	// 	uint64(b[2])<<9 | uint64(b[1])<<10 | uint64(b[0])<<11
// }
