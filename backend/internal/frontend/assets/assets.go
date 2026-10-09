//go:build init

package assets

import _ "embed"

var assets = make(map[string][]byte)

func init() {
	assets[NotFoundImage] = notFoundImageBytes
	assets[FaviconImage] = faviconImageBytes
	assets[AyushAnandImage] = ayushAnandImageBytes
	assets[AnushkaSinghImage] = anushkaSinghImageBytes
	assets[DivyaImage] = divyaImageBytes
	assets[RitumPrabhatImage] = ritumPrabhatImageBytes
	assets[ShivamKumarImage] = shivamKumarImageBytes
	assets[TanishkaGuptaImage] = tanishkaGuptaImageBytes
	assets[AryanKumarImage] = aryanKumarImageBytes
	assets[OmkarPrasadImage] = omkarPrasadImageBytes
	assets[ShanmugasundaramImage] = shanmugasundaramImageBytes
	assets[KaashviKumarImage] = kaashviKumarImageBytes
	assets[AnirbanImage] = anirbanImageBytes
}

//go:embed not_found.png
var notFoundImageBytes []byte

//go:embed favicon.png
var faviconImageBytes []byte

//go:embed Ayush_Anand.jpg
var ayushAnandImageBytes []byte

//go:embed Anushka_Singh.jpg
var anushkaSinghImageBytes []byte

//go:embed Divya.jpg
var divyaImageBytes []byte

//go:embed Ritum_Prabhat.jpg
var ritumPrabhatImageBytes []byte

//go:embed Shivam_Kumar.jpg
var shivamKumarImageBytes []byte

//go:embed Tanishka_Gupta.jpg
var tanishkaGuptaImageBytes []byte

//go:embed Aryan_Kumar.jpg
var aryanKumarImageBytes []byte

//go:embed Omkar_Prasad.jpg
var omkarPrasadImageBytes []byte

//go:embed Shanmugasundaram.jpg
var shanmugasundaramImageBytes []byte

//go:embed Kaashvi_Kumar.jpg
var kaashviKumarImageBytes []byte

//go:embed Anirban.jpg
var anirbanImageBytes []byte
