package jwt

import (
	"fmt"
	"github.com/fangbc5/gogo/core/auth"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	privKey, err := os.ReadFile("pem/sample_key")
	if err != nil {
		t.Fatalf("Unable to read private key: %v", err)
	}
	jwt := NewAuth(auth.WithPrivateKey(string(privKey)))
	tok, err := jwt.Generate("test")
	if err != nil {
		t.Fatalf("jwt generate fail: %v", err)
	}
	fmt.Println(tok.Created)
	fmt.Println(tok.Expiry)
	fmt.Println(tok.AccessToken)
	fmt.Println(tok.RefreshToken)
}

func TestInspect(t *testing.T) {
	pubKey, err := os.ReadFile("pem/sample_key.pub")
	if err != nil {
		t.Fatalf("Unable to read public key: %v", err)
	}

	jwt := NewAuth(auth.WithPublicKey(string(pubKey)))
	inspect, err := jwt.Inspect("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiIiwic2NvcGVzIjpudWxsLCJtZXRhZGF0YSI6bnVsbCwic3ViIjoidGVzdCIsImV4cCI6MTY4MjA0NzA1OH0.i_LD7qz47PoSAHELTV__t3WtlcHBHXU8hEIG2rXXSrffq83LI8dY0RUZbP9tmlrmHxyvl_SCbwYh2tFJHV821Ae5bGJmyHOpN6Zp3acnHQZCix1RKDZW01EYAsweLAIoyW-wqEg27Bz-lKLkUiPFFvkr1o2YkL13QZ095PkXkaTFUoh-zoqDf0H3KqJsUkp-uEXQnd-RjdgNRsXP3256762GHi60Hnysmf4YUWlpZjj41CflYfkFjJB6ciLDNNop5-K5fpynGaaDjmnA3mQd0e73FHw4qkmw_JBOqaaArB5bzF2A_DtyG08ZRVPQHsjoh6mmOdAtmECJfzFI-oH0s5HxZOOsNbi35iSeG7z38Z1ppiNHegP0iER9uGsy2Lf_owm45ORz__fSSA1S-jO4Iw6DmrAofzz1BejAGy3aWCULWE3JXP_H9Yf2iv-PsKo6y3sW_7PmjjSLI-9Jg-qR9qRh9oFmi2SRNUmhuON1BPczAtzS-Ei2S2dKaoYX00ZHs8alvTB-VvEXzet7CRdGFWKtFnngREYIP638LkwNz2sLPmWYrhPuRmgd8QVRTT-c2lk6f1H0xy8IwKLeKNXVcCKxCtZEOocjXLqyKq6wk_-TqBeSRKkjDZt2B08CdQ7_6tezocc8tDVjBZ4XOU5yaJ3APt5j4TJeKs8cD09l7JA")
	if err != nil {
		t.Fatalf("Inspect Token Error: %v", err)
	}
	fmt.Println(inspect)
}

func TestRefreshToken(t *testing.T) {
	privKey, err := os.ReadFile("pem/sample_key")
	if err != nil {
		t.Fatalf("Unable to read private key: %v", err)
	}

	pubKey, err := os.ReadFile("pem/sample_key.pub")
	if err != nil {
		t.Fatalf("Unable to read public key: %v", err)
	}

	jwt := NewAuth(auth.WithPrivateKey(string(privKey)), auth.WithPublicKey(string(pubKey)))
	tok, err := jwt.Refresh(auth.WithRefreshToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiIiwic2NvcGVzIjpudWxsLCJtZXRhZGF0YSI6bnVsbCwic3ViIjoidGVzdCIsImV4cCI6MTY4MjY1MTQwN30.ULCPHHb9M97X88BDPHV-4KJ5lk8bNRkzQRGSfxD_3NQQbhsvMHNDvVgbtdJ2PvAImVaXC4r-A29jdozuF9a9_2MY_8JFUd_aYkuBCZ-KFv9oeICc70mqIksc7MPG5JgCOJBh_b-M4pU3-xaNzRJme90elcc_rvovuSkQ8nvy8he_tasWEa_HC0NxHpywJDsHWaj_Jie2Gg9QQN0NWA-IyVwchaawNo1rSyvmoSEF64fsjFNu5SoqfM2iy0VVEWDOqri0YJZ-oimhe1UeQFn78Leqb__TuLzKldRtvF2dDOeSXG5tvhfzqBROSxnpPGVu5KeZQ-J7ynWAbiK1LrU-GR7ZFhFH-S4vgY9nKvVDP6C613FNtJ_8fxY4WyDJQ3P1WNfxRe25MbZE9-YMe96EGVhma3-5CHfOBKX3aAIA6ft9XwAa--ik6YKe62kX1qkyJexPP1plHH4vYN1xyx6IrEUo5DLgBrIC758tDRxfg4E64Mugw8YuwYqKGLRRg5uav80upU0JWRIqsEAZvRdpp-nYVX-b30FFOY8YxAIawBufDD_0qcAPtSbQI4KJXrsD5HqfFDaQVDg6_rjyuaqfXvBlycbC3kciNpJYZPVNAzp87HPC14UQt8eyeSSBJRDqljc7_Lt_8bENf91-g3z2gkhhjbMHa9Hyo_sT60K49PA"))
	if err != nil {
		t.Fatalf("Inspect Token Error: %v", err)
	}
	fmt.Println(tok.Created)
	fmt.Println(tok.Expiry)
	fmt.Println(tok.AccessToken)
	fmt.Println(tok.RefreshToken)
}
