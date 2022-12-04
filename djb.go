package localcache

import (
	"crypto/rand"
	"math"
	"math/big"
	insecurerand "math/rand"
	"os"
)

func newTimes33() HashFunc {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint64
	if err != nil {
		os.Stderr.Write([]byte("WARNING: NewTimes33() failed to read from the system CSPRNG (/dev/urandom or equivalent.) Your system's security may be compromised. Continuing with an insecure seed.\n"))
		seed = uint64(insecurerand.Uint32())
	} else {
		seed = rnd.Uint64()
	}
	return djb33{seed}
}

type djb33 struct {
	seed uint64
}


func (h djb33) Sum64(k string) uint64 {
	var (
		l = uint64(len(k))
		d = 5381 + h.seed + l
		i = uint64(0)
	)
	// Why is all this 5x faster than a for loop?
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint64(k[i])
			d = (d * 33) ^ uint64(k[i+1])
			d = (d * 33) ^ uint64(k[i+2])
			d = (d * 33) ^ uint64(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint64(k[i])
	case 3:
		d = (d * 33) ^ uint64(k[i])
		d = (d * 33) ^ uint64(k[i+1])
	case 4:
		d = (d * 33) ^ uint64(k[i])
		d = (d * 33) ^ uint64(k[i+1])
		d = (d * 33) ^ uint64(k[i+2])
	}
	return d ^ (d >> 16)
}