package netlifyblob

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func ticketMAC(secret, filename string, exp int64) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(filename))
	mac.Write([]byte{'\n'})
	mac.Write([]byte(strconv.FormatInt(exp, 10)))
	return mac.Sum(nil)
}

func IssueUploadTicket(
	cfg Config,
	filename string,
	ttl time.Duration,
) (exp int64, sig string, err error) {
	if ttl <= 0 || ttl > 24*time.Hour {
		return 0, "", fmt.Errorf("netlifyblob: bad ticket ttl")
	}
	exp = time.Now().Add(ttl).Unix()
	return exp, hex.EncodeToString(ticketMAC(cfg.UploadSecret, filename, exp)), nil
}

func VerifyUploadTicket(
	cfg Config,
	filename string,
	exp int64,
	sig string,
) error {
	want := ticketMAC(cfg.UploadSecret, filename, exp)
	got, err := hex.DecodeString(sig)
	if err != nil {
		return fmt.Errorf("netlifyblob: bad ticket signature")
	}
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return fmt.Errorf("netlifyblob: bad ticket signature")
	}
	now := time.Now().Unix()
	if exp < now {
		return fmt.Errorf("netlifyblob: upload ticket expired")
	}
	if exp-now > 24*time.Hour.Milliseconds()/1000+60 {
		return fmt.Errorf("netlifyblob: bad ticket expiry")
	}
	return nil
}
