// Copyright 2019-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// see the example at https://golang.org/pkg/image/png/#Decode

package main

import "fmt"
import "encoding/base64"
import "strings"
import "io"
import "os"
import "image/png"

const gopher = `iVBORw0KGgoAAAANSUhEUgAAAEsAAAA8CAAAAAALAhhPAAAFfUlEQVRYw62XeWwUVRzHf2+OPbo9d7tsWyiyaZti6eWGAhISoIGKECEKCAiJJkYTiUgTMYSIosYYBBIUIxoSPIINEBDi2VhwkQrVsj1ESgu9doHWdrul7ba73WNm3vOPtsseM9MdwvvrzTs+8/t95ze/33sI5BqiabU6m9En8oNjduLnAEDLUsQXFF8tQ5oxK3vmnNmDSMtrncks9Hhtt/qeWZapHb1ha3UqYSWVl2ZmpWgaXMXGohQAvmeop3bjTRtv6SgaK/Pb9/bFzUrYslbFAmHPp+3WhAYdr+7GN/YnpN46Opv55VDsJkoEpMrY/vO2BIYQ6LLvm0ThY3MzDzzeSJeeWNyTkgnIE5ePKsvKlcg/0T9QMzXalwXMlj54z4c0rh/mzEfr+FgWEz2w6uk8dkzFAgcARAgNp1ZYef8bH2AgvuStbc2/i6CiWGj98y2tw2l4FAXKkQBIf+exyRnteY83LfEwDQAYCoK+P6bxkZm/0966LxcAAILHB56kgD95PPxltuYcMtFTWw/FKkY/6Opf3GGd9ZF+Qp6mzJxzuRSractOmJrH1u8XTvWFHINNkLQLMR+XHXvfPPHw967raE1xxwtA36IMRfkAAG29/7mLuQcb2WOnsJReZGfpiHsSBX81cvMKywYZHhX5hFPtOqPGWZCXnhWGAu6lX91ElKXSalcLXu3UaOXVay57ZSe5f6Gpx7J2MXAsi7EqSp09b/MirKSyJfnfEEgeDjl8FgDAfvewP03zZ+AJ0m9aFRM8eEHBDRKjfcreDXnZdQuAxXpT2NRJ7xl3UkLBhuVGU16gZiGOgZmrSbRdqkILuL/yYoSXHHkl9KXgqNu3PB8oRg0geC5vFmLjad6mUyTKLmF3OtraWDIfACyXqmephaDABawfpi6tqqBZytfQMqOz6S09iWXhktrRaB8Xz4Yi/8gyABDm5NVe6qq/3VzPrcjELWrebVuyY2T7ar4zQyybUCtsQ5Es1FGaZVrRVQwAgHGW2ZCRZshI5bGQi7HesyE972pOSeMM0dSktlzxRdrlqb3Osa6CCS8IJoQQQgBAbTAa5l5epO34rJszibJI8rxLfGzcp1dRosutGeb2VDNgqYrwTiPNsLxXiPi3dz7LiS1WBRBDBOnqEjyy3aQb+/bLiJzz9dIkscVBBLxMfSEac7kO4Fpkngi0ruNBeSOal+u8jgOuqPz12nryMLCniEjtOOOmpt+KEIqsEdocJjYXwrh9OZqWJQyPCTo67LNS/TdxLAv6R5ZNK9npEjbYdT33gRo4o5oTqR34R+OmaSzDBWsAIPhuRcgyoteNi9gF0KzNYWVItPf2TLoXEg+7isNC7uJkgo1iQWOfRSP9NR11RtbZZ3OMG/VhL6jvx+J1m87+RCfJChAtEBQkSBX2PnSiihc/Twh3j0h7qdYQAoRVsRGmq7HU2QRbaxVGa1D6nIOqaIWRjyRZpHMQKWKpZM5feA+lzC4ZFultV8S6T0mzQGhQohi5I8iw+CsqBSxhFMuwyLgSwbghGb0AiIKkSDmGZVmJSiKihsiyOAUs70UkywooYP0bii9GdH4sfr1UNysd3fUyLLMQN+rsmo3grHl9VNJHbbwxoa47Vw5gupIqrZcjPh9R4Nye3nRDk199V+aetmvVtDRE8/+cbgAAgMIWGb3UA0MGLE9SCbWX670TDy1y98c3D27eppUjsZ6fql3jcd5rUe7+ZIlLNQny3Rd+E5Tct3WVhTM5RBCEdiEK0b6B+/ca2gYU393nFj/n1AygRQxPIUA043M42u85+z2SnssKrPl8Mx76NL3E6eXc3be7OD+H4WHbJkKI8AU8irbITQjZ+0hQcPEgId/Fn/pl9crKH02+5o2b9T/eMx7pKoskYgAAAABJRU5ErkJggg==`

// gopherPNG creates an io.Reader by decoding the base64 encoded image data string in the gopher constant.
func gopherPNG() io.Reader { return base64.NewDecoder(base64.StdEncoding, strings.NewReader(gopher)) }


func main() {
  fbdev, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0600)
  img, err := png.Decode(gopherPNG())
	if err != nil {
		fmt.Println(err)
	}
  for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
      c := []byte{byte(r >> 12), byte(g >> 12), byte(b >> 12)}
      fbdev.WriteAt(c, int64(420+x + (y+350)*640))
    }
  }
}
