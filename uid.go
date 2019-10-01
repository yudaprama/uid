package uid

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	mrand "math/rand"
)

const Version = "1.0.0"

const (
	digits   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base     = 62
	preLen   = 12
	seqLen   = 10
	maxSeq   = int64(839299365868340224)
	minInc   = int64(33)
	maxInc   = int64(333)
	totalLen = preLen + seqLen
)

type UID struct {
	pre []byte
	seq int64
	inc int64
}

type lockedUID struct {
	sync.Mutex
	*UID
}

var globalUID *lockedUID

func init() {
	r, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		mrand.Seed(time.Now().UnixNano())
	} else {
		mrand.Seed(r.Int64())
	}
	globalUID = &lockedUID{UID: New()}
	globalUID.RandomizePrefix()
}

func New() *UID {
	n := &UID{
		seq: mrand.Int63n(maxSeq),
		inc: minInc + mrand.Int63n(maxInc-minInc),
		pre: make([]byte, preLen),
	}
	n.RandomizePrefix()
	return n
}

func Next() string {
	globalUID.Lock()
	id := globalUID.Next()
	globalUID.Unlock()
	return id
}

func (n *UID) Next() string {
	n.seq += n.inc
	if n.seq >= maxSeq {
		n.RandomizePrefix()
		n.resetSequential()
	}
	seq := n.seq

	var b [totalLen]byte
	bs := b[:preLen]
	copy(bs, n.pre)

	for i, l := len(b), seq; i > preLen; l /= base {
		i -= 1
		b[i] = digits[l%base]
	}
	return string(b[:])
}

func (n *UID) resetSequential() {
	n.seq = mrand.Int63n(maxSeq)
	n.inc = minInc + mrand.Int63n(maxInc-minInc)
}

func (n *UID) RandomizePrefix() {
	var cb [preLen]byte
	cbs := cb[:]
	if nb, err := rand.Read(cbs); nb != preLen || err != nil {
		panic(fmt.Sprintf("failed generate crypto random number: %v\n", err))
	}

	for i := 0; i < preLen; i++ {
		n.pre[i] = digits[int(cbs[i])%base]
	}
}
